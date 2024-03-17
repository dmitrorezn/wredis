package wredis

import (
	"context"
	"sort"
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
	LRU          LRUConfig
	SingleFlight SingleFlightConfigs
}

func NewBuildConfig() FactoryConfig {
	return FactoryConfig{
		LRU:          NewLRUConfig(),
		SingleFlight: make(SingleFlightConfigs, 0),
	}
}

func (bc FactoryConfig) WithLRU(lru LRUConfig) FactoryConfig {
	bc.LRU = lru

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
	case LRU:
		return NewWithLRU(ctx, client, cfg.LRU)
	case SingleFlight:
		return NewSingleFlight(client, cfg.SingleFlight...)
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
