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

type command[T any] func(client UniversalClient) T

func newGroup[T any](
	g *singleflight.Group,
	client UniversalClient,
) *group[T] {
	return &group[T]{
		group:  g,
		client: client,
	}
}

type group[T any] struct {
	group  *singleflight.Group
	client UniversalClient
}

func (g *group[T]) Do(key string, cmd command[T]) (T, error, bool) {
	val, err, shared := g.group.Do(key, func() (interface{}, error) {
		return cmd(g.client), nil
	})
	return val.(T), err, shared
}

func (g *SingleFlightClient) prepareStringCmd() *group[*redis.StringCmd] {
	return newGroup[*redis.StringCmd](g.group, g.UniversalClient)
}

func (g *SingleFlightClient) prepareSliceCmd() *group[*redis.SliceCmd] {
	return newGroup[*redis.SliceCmd](g.group, g.UniversalClient)
}

func (s *SingleFlightClient) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd, _, shared := s.prepareStringCmd().Do(key, func(client UniversalClient) *redis.StringCmd {
		return client.Get(ctx, key)
	})
	if shared {
		s.sharedInMinuteCount.Add(1)
	}

	return cmd
}

func (s *SingleFlightClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	cmd, _, shared := s.prepareStringCmd().Do(key+field, func(client UniversalClient) *redis.StringCmd {
		return client.HGet(ctx, key, field)
	})
	if shared {
		s.sharedInMinuteCount.Add(1)
	}

	return cmd
}

func (s *SingleFlightClient) GetDel(ctx context.Context, key string) *redis.StringCmd {
	cmd, _, shared := s.prepareStringCmd().Do(key, func(client UniversalClient) *redis.StringCmd {
		return client.GetDel(ctx, key)
	})
	if shared {
		s.sharedInMinuteCount.Add(1)
	}

	return cmd
}

func (s *SingleFlightClient) HMGet(ctx context.Context, key string, field ...string) *redis.SliceCmd {
	cmd, _, shared := s.prepareSliceCmd().Do(key+fieldsToKey(field...), func(client UniversalClient) *redis.SliceCmd {
		return client.HMGet(ctx, key, field...)
	})
	if shared {
		s.sharedInMinuteCount.Add(1)
	}

	return cmd
}
