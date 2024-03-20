package wredis

import (
	"context"
	"github.com/dmitrorezn/wredis/pkg/logger"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"sync/atomic"
)

type SingleFlightClient struct {
	group *singleflight.Group
	UniversalClient
	logger logger.Logger

	sharedInMinuteCount *atomic.Int64
}

type SingleFlightConfigs []Config[SingleFlightClient]

func WithLogger(logger logger.Logger) func(client *SingleFlightClient) {
	return func(client *SingleFlightClient) {
		client.logger = logger
	}
}

func WithCounter(sharedInMinuteCount *atomic.Int64) func(client *SingleFlightClient) {
	return func(client *SingleFlightClient) {
		client.sharedInMinuteCount = sharedInMinuteCount
	}
}

func NewSingleFlight(
	client UniversalClient,
	conf ...Config[SingleFlightClient],
) *SingleFlightClient {
	sfc := &SingleFlightClient{
		group:               new(singleflight.Group),
		UniversalClient:     client,
		logger:              logger.NewStdLogger(),
		sharedInMinuteCount: new(atomic.Int64),
	}
	for _, c := range conf {
		c(sfc)
	}

	return sfc
}

func (g *SingleFlightClient) getSingleStringCmd() *group[*redis.StringCmd] {
	return newGroup[*redis.StringCmd](g.group, g.UniversalClient)
}

func (g *SingleFlightClient) getSingleSliceCmd() *group[*redis.SliceCmd] {
	return newGroup[*redis.SliceCmd](g.group, g.UniversalClient)
}

func (g *SingleFlightClient) getSingleMapStringStringCmd() *group[*redis.MapStringStringCmd] {
	return newGroup[*redis.MapStringStringCmd](g.group, g.UniversalClient)
}

func (s *SingleFlightClient) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd, _, shared := s.getSingleStringCmd().Do(key, func(client UniversalClient) *redis.StringCmd {
		return client.Get(ctx, key)
	})
	if shared {
		s.sharedInMinuteCount.Add(1)
	}

	return cmd
}

func (s *SingleFlightClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	cmd, _, shared := s.getSingleStringCmd().Do(key+field, func(client UniversalClient) *redis.StringCmd {
		return client.HGet(ctx, key, field)
	})
	if shared {
		s.sharedInMinuteCount.Add(1)
	}

	return cmd
}

func (s *SingleFlightClient) GetDel(ctx context.Context, key string) *redis.StringCmd {
	cmd, _, shared := s.getSingleStringCmd().Do(key, func(client UniversalClient) *redis.StringCmd {
		return client.GetDel(ctx, key)
	})
	if shared {
		s.sharedInMinuteCount.Add(1)
	}

	return cmd
}

func (s *SingleFlightClient) HMGet(ctx context.Context, key string, field ...string) *redis.SliceCmd {
	single := s.getSingleSliceCmd()
	cmd, _, _ := single.Do(key+fieldsToKey(field...), func(client UniversalClient) *redis.SliceCmd {
		return client.HMGet(ctx, key, field...)
	})

	return cmd
}

func (s *SingleFlightClient) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	single := s.getSingleMapStringStringCmd()
	cmd, _, _ := single.Do(key, func(client UniversalClient) *redis.MapStringStringCmd {
		return client.HGetAll(ctx, key)
	})

	return cmd
}
