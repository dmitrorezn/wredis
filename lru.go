package wredis

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/dmitrorezn/golang-lru/v2/expirable"
	"github.com/redis/go-redis/v9"

	"github.com/dmitrorezn/wredis/pkg/logger"
)

type LRUClient struct {
	cfg LRUConfig
	UniversalClient
	logger logger.Logger
	cache  *expirable.LRU[string, any]
	quit   chan struct{}
}

type callback[K comparable, V any] expirable.EvictCallback[K, V]

func (c callback[K, V]) Call(key K, value V) {
	if c == nil {
		return
	}

	c(key, value)
}

type LRUConfig struct {
	GroupName          string
	Size               int           `env:"WREDIS_LRU_SIZE"`
	TTL                time.Duration `env:"WREDIS_LRU_TTL"`
	CacheErrors        bool          `env:"WREDIS_LRU_CACHE_ERRS"`
	ErrTTL             time.Duration `env:"WREDIS_LRU_ERR_TTL"`
	InvalidationEvents Events        `env:"WREDIS_LRU_INVALIDATION_EVENTS" envSeparator:","`
	SubscriptionChannelConfig
	OnAdd   callback[string, any]
	OnEvict callback[string, any]
	OnGet   callback[string, any]
	logger  logger.Logger
}

type SubscriptionChannelConfig struct {
	HealthCheckInterval time.Duration `env:"WREDIS_LRU_SUB_HEALS_INTERVAL"`
	SendTimeout         time.Duration `env:"WREDIS_LRU_SUB_SEND_TIMEOUT"`
	Size                int           `env:"WREDIS_LRU_SUB_BUF_SIZE"`
	Workers             int           `env:"WREDIS_LRU_SUB_WORKERS"`
}

const (
	defaultLruTTL       = time.Minute
	defaultLruErrTTL    = 10 * time.Second
	defaultLruSize      = 1_000
	subBufSize          = 1_000
	subWorkers          = 10
	healthCheckInterval = time.Minute
	sendTimeout         = 100 * time.Millisecond
)

func NewLRUConfig() LRUConfig {
	return LRUConfig{
		Size:   defaultLruSize,
		TTL:    defaultLruTTL,
		ErrTTL: defaultLruErrTTL,
		logger: logger.NewStdLogger(),
		SubscriptionChannelConfig: SubscriptionChannelConfig{
			HealthCheckInterval: healthCheckInterval,
			Size:                subBufSize,
			Workers:             subWorkers,
			SendTimeout:         sendTimeout,
		},
	}
}

func (lc *LRUConfig) LoadFromEnv() error {
	return env.Parse(lc)
}

func (lc LRUConfig) WithCacheErrors(needCacheErrors bool, errTTL time.Duration) LRUConfig {
	lc.CacheErrors = needCacheErrors
	lc.ErrTTL = errTTL

	return lc
}

func (lc LRUConfig) WithSize(size int) LRUConfig {
	lc.Size = size

	return lc
}

func (lc LRUConfig) WithTTL(ttl time.Duration) LRUConfig {
	lc.TTL = ttl

	return lc
}

func (lc LRUConfig) WithGroupName(name string) LRUConfig {
	lc.GroupName = name

	return lc
}

func (lc LRUConfig) WithOnAdd(onAdd callback[string, any]) LRUConfig {
	lc.OnAdd = onAdd

	return lc
}

func (lc LRUConfig) WitOnEvict(onEvict callback[string, any]) LRUConfig {
	lc.OnEvict = onEvict

	return lc
}

func (lc LRUConfig) WithOnGet(onGet callback[string, any]) LRUConfig {
	lc.OnGet = onGet

	return lc
}

func (lc LRUConfig) WitLogger(logger logger.Logger) LRUConfig {
	lc.logger = logger

	return lc
}

func (lc LRUConfig) WithCallbackSubChannelCfg(
	size int,
	sendTimeout time.Duration,
	healthCheckInterval time.Duration,
) LRUConfig {
	lc.SubscriptionChannelConfig.Size = size
	lc.SubscriptionChannelConfig.SendTimeout = sendTimeout
	lc.SubscriptionChannelConfig.HealthCheckInterval = healthCheckInterval

	return lc
}

func (lc LRUConfig) WithInvalidationCallbackEvent(events ...Event) LRUConfig {
	lc.InvalidationEvents = events

	return lc
}

func NewWithLRU(ctx context.Context, client UniversalClient, cfg LRUConfig) *LRUClient {
	lc := &LRUClient{
		cfg:             cfg,
		logger:          cfg.logger,
		UniversalClient: client,
		cache:           expirable.NewLRU[string, any](cfg.Size, cfg.OnEvict.Call, cfg.TTL),
		quit:            make(chan struct{}),
	}

	lc.start(ctx)

	return lc
}

func (c *LRUClient) start(ctx context.Context) {
	c.subscribeForCallbacks(ctx)
}

func (c *LRUClient) subscribeForCallbacks(ctx context.Context) {
	if c.cfg.InvalidationEvents.Empty() {
		return
	}

	var (
		subCfg       = c.cfg.SubscriptionChannelConfig
		subPattern   = keysPatternChannel.String(c.cfg.InvalidationEvents.String())
		subscription = c.UniversalClient.PSubscribe(ctx, subPattern)
	)
	subChan := subscription.Channel(
		redis.WithChannelHealthCheckInterval(subCfg.HealthCheckInterval),
		redis.WithChannelSendTimeout(subCfg.SendTimeout),
		redis.WithChannelSize(subCfg.Size),
	)

	var pubSub = sub{subPattern, subscription, c.logger}
	for i := 0; i <= subCfg.Workers; i++ {
		go c.processMessage(ctx, subChan)
	}

	go pubSub.waitQuit(ctx, c.quit)
}

type sub struct {
	pattern string
	*redis.PubSub
	logger logger.Logger
}

func (s sub) waitQuit(ctx context.Context, quit chan struct{}) {
	select {
	case <-quit:
	case <-ctx.Done():
	}
	if err := s.PUnsubscribe(ctx, s.pattern); err != nil {
		s.logger.Error("subscribeForCallbacks -> PUnsubscribe",
			"channel", s.pattern, "err", err,
		)
	}
}

func (c *LRUClient) processMessage(ctx context.Context, subChan <-chan *redis.Message) {
loop:
	for {
		select {
		case msg := <-subChan:
			c.del(msg.Payload)
		case <-c.quit:
			break loop
		case <-ctx.Done():
			break loop
		}
	}
}

func (c *LRUClient) del(key string) {
	c.cache.Remove(key)
}

func (c *LRUClient) Get(ctx context.Context, key string) (cmd *redis.StringCmd) {
	ckey := key
	cmd, ok := c.getStringCmd(ckey)
	if ok {
		return cmd
	}
	cmd = c.UniversalClient.Get(ctx, ckey)
	c.add(ckey, cmd)

	return cmd
}

func (c *LRUClient) GetDel(ctx context.Context, key string) *redis.StringCmd {
	cmd, ok := c.getStringCmd(key)
	if !ok || ok && cmd.Err() != nil {
		return c.UniversalClient.GetDel(ctx, key)
	}
	c.del(key)
	if err := c.UniversalClient.Del(ctx, key).Err(); err != nil {
		cmd.SetErr(err)
	}

	return cmd
}

func (c *LRUClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	cacheKey := key + field
	cmd, ok := c.getStringCmd(cacheKey)
	if ok {
		return cmd
	}

	cmd = c.UniversalClient.HGet(ctx, cacheKey, field)
	c.add(cacheKey, cmd)

	return cmd
}

func (c *LRUClient) HMGet(ctx context.Context, key string, field ...string) *redis.SliceCmd {
	cacheKey := fieldsToKey(field...)
	cmd, ok := c.getSliceCmd(cacheKey)
	if ok {
		return cmd
	}
	cmd = c.UniversalClient.HMGet(ctx, cacheKey, field...)
	c.add(key, cmd)

	return cmd
}

func (c *LRUClient) Close() error {
	close(c.quit)

	c.cache.Purge()

	return c.UniversalClient.Close()
}

type IErrCmd interface {
	Err() error
}

func (c *LRUClient) add(key string, cmd IErrCmd) {
	var ttl = c.cfg.TTL
	if err := cmd.Err(); err != nil {
		if !c.cfg.CacheErrors {
			return
		}
		if c.cfg.ErrTTL > time.Nanosecond {
			ttl = c.cfg.ErrTTL
		}
	}
	if c.cache.Set(key, cmd, ttl) {
		c.cfg.OnEvict.Call(key, cmd)
	}
	c.cfg.OnAdd.Call(key, cmd)
}

func (c *LRUClient) getStringCmd(key string) (*redis.StringCmd, bool) {
	value, ok := c.cache.Get(key)
	if !ok {
		return nil, false
	}
	cmd, asserted := value.(*redis.StringCmd)
	if !asserted || cmd == nil {
		return nil, false
	}
	c.cfg.OnGet.Call(key, value)

	return cmd, true
}

func (c *LRUClient) getSliceCmd(key string) (*redis.SliceCmd, bool) {
	value, ok := c.cache.Get(key)
	if !ok {
		return nil, false
	}
	cmd, asserted := value.(*redis.SliceCmd)
	if !asserted || cmd == nil {
		panic("getSliceCmd")
		return nil, false
	}
	c.cfg.OnGet.Call(key, value)

	return cmd, true
}

func fieldsToKey(field ...string) string {
	sort.Strings(field)

	return strings.Join(field, "")
}
