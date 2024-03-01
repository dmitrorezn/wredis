package wredis

import (
	"context"
	"github.com/dmitrorezn/wredis/pkg/logger"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type SingleFlightClient struct {
	group *singleflight.Group
	UniversalClient
	logger logger.Logger

	sharedInMinuteCount atomic.Int64
}

func WithLogger(logger logger.Logger) func(client *SingleFlightClient) {
	return func(client *SingleFlightClient) {
		client.logger = logger
	}
}

const (
	logMetricInterval = time.Minute
)

func NewSingleFlight(
	client UniversalClient,
	conf ...Config[SingleFlightClient],
) *SingleFlightClient {
	sfc := &SingleFlightClient{
		group:           new(singleflight.Group),
		UniversalClient: client,
		logger:          logger.NewStdLogger(),
	}
	for _, c := range conf {
		c(sfc)
	}
	go sfc.handleMetric(logMetricInterval)

	return sfc

}
func (s *SingleFlightClient) handleMetric(interval time.Duration) {
	for range time.Tick(interval) {
		s.logAndPurgeMetric()
	}
}
func (s *SingleFlightClient) logAndPurgeMetric() {
	s.logger.Info("METRIC", "shared_keys_count", s.sharedInMinuteCount.Swap(0))
}

type cmd[T any] func(client UniversalClient) T

func newGroup[T any](
	g *singleflight.Group,
	client UniversalClient,
	onShare func(key string),
) *group[T] {
	return &group[T]{
		group:   g,
		client:  client,
		onShare: onShare,
	}
}

type group[T any] struct {
	group   *singleflight.Group
	client  UniversalClient
	onShare func(key string)
}

func (g *group[T]) Do(key string, cmd cmd[T]) (T, bool) {
	val, _, shared := g.group.Do(key, func() (interface{}, error) {
		return cmd(g.client), nil
	})
	if shared && g.onShare != nil {
		g.onShare(key)
	}

	return val.(T), shared
}

func (g *SingleFlightClient) prepareStringCmd() *group[*redis.StringCmd] {
	return newGroup[*redis.StringCmd](g.group, g.UniversalClient, func(key string) {
		g.sharedInMinuteCount.Add(1)
	})
}

func (g *SingleFlightClient) prepareSliceCmd() *group[*redis.SliceCmd] {
	return newGroup[*redis.SliceCmd](g.group, g.UniversalClient, func(key string) {
		g.sharedInMinuteCount.Add(1)
	})
}

func (s *SingleFlightClient) do(_ context.Context, key string, fn func() (interface{}, error)) (interface{}, error) {
	result, err, shared := s.group.Do(key, fn)
	if shared {
		s.sharedInMinuteCount.Add(1)
	}

	return result, err
}

func (s *SingleFlightClient) Get(ctx context.Context, key string) (cmd *redis.StringCmd) {
	cmd, _ = s.prepareStringCmd().Do(key, func(client UniversalClient) *redis.StringCmd {
		return client.Get(ctx, key)
	})

	return cmd
}

func (s *SingleFlightClient) HGet(ctx context.Context, key, field string) (cmd *redis.StringCmd) {
	cmd, _ = s.prepareStringCmd().Do(key+field, func(client UniversalClient) *redis.StringCmd {
		return client.HGet(ctx, key, field)
	})

	return cmd
}

func (s *SingleFlightClient) GetDel(ctx context.Context, key string) (cmd *redis.StringCmd) {
	cmd, _ = s.prepareStringCmd().Do(key, func(client UniversalClient) *redis.StringCmd {
		return client.GetDel(ctx, key)
	})

	return cmd
}

func (s *SingleFlightClient) HMGet(ctx context.Context, key string, field ...string) (cmd *redis.SliceCmd) {
	cmd, _ = s.prepareSliceCmd().Do(key+fieldsToKey(field...), func(client UniversalClient) *redis.SliceCmd {
		return client.HMGet(ctx, key, field...)
	})

	return cmd
}
