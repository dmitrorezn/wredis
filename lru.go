package wredis

import (
	"context"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/redis/go-redis/v9"
)

type LRUClient struct {
	cfg LRUConfig
	UniversalClient
	logger Logger
	cache  *expirable.LRU[string, any]
	quit   chan struct{}
	work   chan string
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
	Size               int
	TTL                time.Duration
	InvalidationEvents Events
	SubChannel         SubscriptionChannelConfig
	OnAdd              callback[string, any]
	OnEvict            callback[string, any]
	logger             Logger
}

type SubscriptionChannelConfig struct {
	HealthCheckInterval time.Duration
	SendTimeout         time.Duration
	Size                int
}

const (
	defaultLruTTL  = time.Minute
	defaultLruSize = 1_000
)

func NewLRUConfig() LRUConfig {
	return LRUConfig{
		Size:   defaultLruSize,
		TTL:    defaultLruTTL,
		logger: NewStdLogger(),
	}
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

func (lc LRUConfig) WitOnAdd(c callback[string, any]) LRUConfig {
	lc.OnAdd = c

	return lc
}

func (lc LRUConfig) WitOnEvict(c callback[string, any]) LRUConfig {
	lc.OnEvict = c

	return lc
}
func (lc LRUConfig) WitLogger(logger Logger) LRUConfig {
	lc.logger = logger

	return lc
}

func (lc LRUConfig) WithCallbackSubChannelCfg(
	size int,
	sendTimeout time.Duration,
	healthCheckInterval time.Duration,
) LRUConfig {
	lc.SubChannel.Size = size
	lc.SubChannel.SendTimeout = sendTimeout
	lc.SubChannel.HealthCheckInterval = healthCheckInterval

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
	go c.subscribeForCallbacks(ctx)
}

func (c *LRUClient) subscribeForCallbacks(ctx context.Context) {
	if len(c.cfg.InvalidationEvents) == 0 {
		return
	}

	var (
		chCfg        = c.cfg.SubChannel
		subPattern   = keysPatternChannel.String(c.cfg.InvalidationEvents.String())
		subscription = c.UniversalClient.PSubscribe(ctx, subPattern)
	)
	subChan := subscription.Channel(
		redis.WithChannelHealthCheckInterval(chCfg.HealthCheckInterval),
		redis.WithChannelSendTimeout(chCfg.SendTimeout),
		redis.WithChannelSize(chCfg.Size),
	)

	go c.worker(ctx, sub{subPattern, subscription}, subChan)
}

type sub struct {
	pattern string
	*redis.PubSub
}

func (c *LRUClient) worker(ctx context.Context, sub sub, subChan <-chan *redis.Message) {
	for {
		select {
		case msg := <-subChan:
			go c.del(msg.Payload)
		case <-ctx.Done():
			if err := sub.PUnsubscribe(ctx, sub.pattern); err != nil {
				c.logger.Error("subscribeForCallbacks -> PUnsubscribe",
					"channel", sub.pattern, "err", err,
				)
			}

			return
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

	cmd = c.UniversalClient.Get(ctx, key)
	c.add(key, cmd)

	return cmd
}

func (c *LRUClient) GetDel(ctx context.Context, key string) *redis.StringCmd {
	ckey := key
	cmd, ok := c.getStringCmd(ckey)
	if ok {
		c.del(ckey)
		dcmd := c.UniversalClient.Del(ctx, key)
		cmd.SetErr(dcmd.Err())

		return cmd
	}

	return c.UniversalClient.GetDel(ctx, key)
}

func (c *LRUClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	ckey := key + field
	cmd, ok := c.getStringCmd(ckey)
	if ok {
		return cmd
	}

	cmd = c.UniversalClient.HGet(ctx, key, field)
	c.add(key, cmd)

	return cmd
}

func (c *LRUClient) HMGet(ctx context.Context, key string, field ...string) *redis.SliceCmd {
	ckey := fieldsToKey(field...)
	cmd, ok := c.getSliceCmd(ckey)
	if ok {
		return cmd
	}
	cmd = c.UniversalClient.HMGet(ctx, key, field...)
	c.add(key, cmd)

	return cmd
}

type IErrCmd interface {
	Err() error
}

func (c *LRUClient) add(key string, anyCmd any) {
	switch cmd := anyCmd.(type) {
	case IErrCmd:
		if err := cmd.Err(); err != nil {
			return
		}
	}

	if c.cache.Add(key, anyCmd) {
		c.cfg.OnAdd.Call(key, anyCmd)
	}
}

func (c *LRUClient) getStringCmd(key string) (cmd *redis.StringCmd, ok bool) {
	value, ok := c.cache.Get(key)
	if ok {
		cmd, ok = value.(*redis.StringCmd)
	}
	if !ok || cmd == nil {
		return cmd, false
	}

	return cmd, true
}

func (c *LRUClient) getSliceCmd(key string) (cmd *redis.SliceCmd, ok bool) {
	value, ok := c.cache.Get(key)
	if ok {
		cmd, ok = value.(*redis.SliceCmd)
	}
	if !ok || cmd == nil {
		return cmd, false
	}

	return cmd, true
}
