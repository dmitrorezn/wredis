package wredis

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/redis/go-redis/v9"

	"github.com/dmitrorezn/golang-lru/v2/expirable"
	"github.com/dmitrorezn/wredis/pkg/logger"
)

type CacheClient struct {
	cfg LocalCacheConfig
	UniversalClient
	logger logger.Logger
	cache  Cache
	quit   chan struct{}
}

type callback[K comparable, V any] expirable.EvictCallback[K, V]

func (c callback[K, V]) Call(key K, value V) {
	if c == nil {
		return
	}

	c(key, value)
}

type LocalCacheConfig struct {
	GroupName string
	Size      int           `env:"WREDIS_LOCAL_CACHE_SIZE"`
	Samples   int           `env:"WREDIS_LOCAL_CACHE_SAMPLES"`
	TTL       time.Duration `env:"WREDIS_LOCAL_CACHE_TTL"`
	TTLOffset time.Duration `env:"WREDIS_LOCAL_CACHE_TTL"`

	CacheErrors        bool          `env:"WREDIS_LOCAL_CACHE_ERRS"`
	ErrTTL             time.Duration `env:"WREDIS_LOCAL_CACHE_ERR_TTL"`
	InvalidationEvents Events        `env:"WREDIS_LOCAL_CACHE_INVALIDATION_EVENTS" envSeparator:","`
	SubscriptionChannelConfig
	OnAdd   callback[string, any]
	OnEvict callback[string, any]
	OnGet   callback[string, any]
	logger  logger.Logger
}

type SubscriptionChannelConfig struct {
	HealthCheckInterval time.Duration `env:"WREDIS_LOCAL_CACHE_SUB_HEALS_INTERVAL"`
	SendTimeout         time.Duration `env:"WREDIS_LOCAL_CACHE_SUB_SEND_TIMEOUT"`
	Size                int           `env:"WREDIS_LOCAL_CACHE_SUB_BUF_SIZE"`
	Workers             int           `env:"WREDIS_LOCAL_CACHE_SUB_WORKERS"`
}

const (
	defaultTTL          = time.Minute
	defaultErrTTL       = 20 * time.Second
	defaultSize         = 1_000
	defaultSamples      = 1_000
	subBufSize          = 1_000
	subWorkers          = 10
	healthCheckInterval = time.Minute
	sendTimeout         = 100 * time.Millisecond
)

func NewLocalCacheConfig() LocalCacheConfig {
	return LocalCacheConfig{
		Size:    defaultSize,
		Samples: defaultSamples,
		TTL:     defaultTTL,
		ErrTTL:  defaultErrTTL,
		logger:  logger.NewStdLogger(),
		SubscriptionChannelConfig: SubscriptionChannelConfig{
			HealthCheckInterval: healthCheckInterval,
			Size:                subBufSize,
			Workers:             subWorkers,
			SendTimeout:         sendTimeout,
		},
	}
}

func (lc *LocalCacheConfig) LoadFromEnv() error {
	return env.Parse(lc)
}

func (lc LocalCacheConfig) WithCacheErrors(needCacheErrors bool, errTTL time.Duration) LocalCacheConfig {
	lc.CacheErrors = needCacheErrors
	lc.ErrTTL = errTTL

	return lc
}

func (lc LocalCacheConfig) WithSize(size int) LocalCacheConfig {
	lc.Size = size

	return lc
}

func (lc LocalCacheConfig) WithSamples(samples int) LocalCacheConfig {
	lc.Samples = samples

	return lc
}

func (lc LocalCacheConfig) WithTTL(ttl time.Duration) LocalCacheConfig {
	lc.TTL = ttl

	return lc
}

func (lc LocalCacheConfig) WithGroupName(name string) LocalCacheConfig {
	lc.GroupName = name

	return lc
}

func (lc LocalCacheConfig) WithOnAdd(onAdd callback[string, any]) LocalCacheConfig {
	lc.OnAdd = onAdd

	return lc
}

func (lc LocalCacheConfig) WitOnEvict(onEvict callback[string, any]) LocalCacheConfig {
	lc.OnEvict = onEvict

	return lc
}

func (lc LocalCacheConfig) WithOnGet(onGet callback[string, any]) LocalCacheConfig {
	lc.OnGet = onGet

	return lc
}

func (lc LocalCacheConfig) WitLogger(logger logger.Logger) LocalCacheConfig {
	lc.logger = logger

	return lc
}

func (lc LocalCacheConfig) WithCallbackSubChannelCfg(
	size int,
	sendTimeout time.Duration,
	healthCheckInterval time.Duration,
) LocalCacheConfig {
	lc.SubscriptionChannelConfig.Size = size
	lc.SubscriptionChannelConfig.SendTimeout = sendTimeout
	lc.SubscriptionChannelConfig.HealthCheckInterval = healthCheckInterval

	return lc
}

func (lc LocalCacheConfig) WithInvalidationCallbackEvent(events ...Event) LocalCacheConfig {
	lc.InvalidationEvents = events

	return lc
}

func NewWithCache(ctx context.Context, client UniversalClient, cache Cache, cfg LocalCacheConfig) *CacheClient {
	lc := &CacheClient{
		cfg:             cfg,
		logger:          cfg.logger,
		UniversalClient: client,
		cache:           cache,
		quit:            make(chan struct{}),
	}

	lc.start(ctx)

	return lc
}

func (c *CacheClient) start(ctx context.Context) {
	c.subscribeForCallbacks(ctx)
}

func (c *CacheClient) subscribeForCallbacks(ctx context.Context) {
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

func (c *CacheClient) processMessage(ctx context.Context, subChan <-chan *redis.Message) {
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

func (c *CacheClient) del(key string) {
	c.cache.Del(key)
}

func (c *CacheClient) Get(ctx context.Context, key string) (cmd *redis.StringCmd) {
	cmd, ok := c.getStringCmd(key)
	if ok {
		return cmd
	}
	cmd = c.UniversalClient.Get(ctx, key)
	c.set(key, cmd)

	return cmd
}

func (c *CacheClient) GetDel(ctx context.Context, key string) *redis.StringCmd {
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

func (c *CacheClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	cacheKey := fieldsToKey(key, field)
	cmd, ok := c.getStringCmd(cacheKey)
	if ok {
		return cmd
	}
	cmd = c.UniversalClient.HGet(ctx, key, field)
	c.set(cacheKey, cmd)

	return cmd
}

func (c *CacheClient) HMGet(ctx context.Context, key string, field ...string) *redis.SliceCmd {
	cacheKey := key + fieldsToKey(field...)
	cmd, ok := c.getSliceCmd(cacheKey)
	if ok {
		return cmd
	}
	cmd = c.UniversalClient.HMGet(ctx, key, field...)
	c.set(key, cmd)

	return cmd
}

func (c *CacheClient) Close() error {
	close(c.quit)

	return c.UniversalClient.Close()
}

type IErrCmd interface {
	Err() error
}

func (c *CacheClient) set(key string, cmd IErrCmd) {
	var ttl = c.cfg.TTL
	if err := cmd.Err(); err != nil {
		if !c.cfg.CacheErrors {
			return
		}
		if c.cfg.ErrTTL > 0 {
			ttl = c.cfg.ErrTTL
		}
	}
	if c.cache.Set(key, cmd, ttl) {
		c.cfg.OnEvict.Call(key, cmd)
	}
	c.cfg.OnAdd.Call(key, cmd)
}

func (c *CacheClient) getStringCmd(key string) (*redis.StringCmd, bool) {
	value, ok := c.cache.Get(key)
	if !ok {
		return nil, false
	}
	cmd, asserted := value.(*redis.StringCmd)
	if !asserted {
		return nil, false
	}
	c.cfg.OnGet.Call(key, value)

	return cmd, true
}

func (c *CacheClient) getSliceCmd(key string) (*redis.SliceCmd, bool) {
	value, ok := c.cache.Get(key)
	if !ok {
		return nil, false
	}
	cmd, asserted := value.(*redis.SliceCmd)
	if !asserted || cmd == nil {
		return nil, false
	}
	c.cfg.OnGet.Call(key, value)

	return cmd, true
}

func fieldsToKey(field ...string) string {
	sort.Strings(field)

	return strings.Join(field, ":")
}
