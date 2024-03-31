package wredis

import (
	"context"
	"errors"
	"sort"
	"time"
)

type Factory struct {
	cfg        FactoryConfig
	decorators decorators
}

func NewFactory(cfg FactoryConfig, ds ...decoratorID) *Factory {
	sort.Sort(decorators(ds))

	return &Factory{
		cfg:        cfg,
		decorators: ds,
	}
}

type FactoryConfig struct {
	CacheConfig  LocalCacheConfig
	SingleFlight SingleFlightConfigs
}

func NewBuildConfig() FactoryConfig {
	return FactoryConfig{
		CacheConfig:  NewLocalCacheConfig(),
		SingleFlight: make(SingleFlightConfigs, 0),
	}
}

func (bc FactoryConfig) WithCacheConfig(cacheConfig LocalCacheConfig) FactoryConfig {
	bc.CacheConfig = cacheConfig

	return bc
}

func (bc FactoryConfig) WithSingleFlight(opts ...Config[SingleFlightClient]) FactoryConfig {
	bc.SingleFlight = opts

	return bc
}

func (f Factory) decorate(
	ctx context.Context,
	client UniversalClient,
	id decoratorID,
	cfg FactoryConfig,
) UniversalClient {
	switch id {
	case LRUDecorator:
		{
			cache := NewLRUCache(
				cfg.CacheConfig.Size, cfg.CacheConfig.OnEvict, cfg.CacheConfig.TTL,
			)

			return NewWithCache(ctx, client, NewTTLCache(cache, time.Second), cfg.CacheConfig)
		}
	case LFUDecorator:
		{
			cache := NewLFUCache(
				cfg.CacheConfig.Size, cfg.CacheConfig.Samples, cfg.CacheConfig.OnEvict, cfg.CacheConfig.TTL,
			)

			return NewWithCache(ctx, client, NewTTLCache(cache, time.Second), cfg.CacheConfig)
		}
	case SingleFlightDecorator:
		{
			return NewSingleFlight(client, cfg.SingleFlight...)
		}
	}

	return client
}

var ErrCacheDecoratorDuplicate = errors.New("cache decorator duplicate")

func (f Factory) Build(
	ctx context.Context,
	options IOptions,
	configurations ...Configuration,
) (client UniversalClient, err error) {
	if !f.decorators.UniqCache() {
		return nil, ErrCacheDecoratorDuplicate
	}

	switch opt := (options).(type) {
	case Options:
		client = New(opt, configurations...)
	case ClusterOptions:
		client = NewCluster(opt, configurations...)
	case RingOptions:
		client = NewRing(opt, configurations...)
	}
	for _, d := range f.decorators {
		client = f.decorate(ctx, client, d, f.cfg)
	}

	return client, nil
}
