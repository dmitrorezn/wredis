package wredis

import (
	"context"
	"sort"
)

type Factory struct {
	cfg        BuildConfig
	decorators decorators
}

func NewFactory(cfg BuildConfig, ds ...decoratorID) *Factory {
	sort.Sort(decorators(ds))

	return &Factory{
		cfg:        cfg,
		decorators: ds,
	}
}

type BuildConfig struct {
	LRU LRUConfig
}

func NewBuildConfig() BuildConfig {
	return BuildConfig{
		LRU: NewLRUConfig(),
	}
}

func (bc BuildConfig) WithLRU(lru LRUConfig) BuildConfig {
	bc.LRU = lru

	return bc
}

func (f Factory) decorate(
	ctx context.Context,
	client UniversalClient,
	id decoratorID,
	cfg BuildConfig,
) UniversalClient {
	switch id {
	case LRU:
		return NewWithLRU(ctx, client, cfg.LRU)
	case SingleFlight:
		return NewSingleFlight(client)
	}

	return client
}

func (f Factory) Build(
	ctx context.Context,
	options IOptions,
	configurations ...Configuration,
) (client UniversalClient) {
	switch opt := (options).(type) {
	case Options:
		client = New(opt, configurations...)
	case ClusterOptions:
		client = NewCluster(opt, configurations...)
	}
	for _, d := range f.decorators {
		client = f.decorate(ctx, client, d, f.cfg)
	}

	return client
}
