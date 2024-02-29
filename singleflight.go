package wredis

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type SingleFlightClient struct {
	group *singleflight.Group
	UniversalClient
	logger Logger

	sharedInMinuteCount atomic.Int64
}

func WithLogger(logger Logger) func(client *SingleFlightClient) {
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
		logger:          NewStdLogger(),
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

func (s *SingleFlightClient) do(ctx context.Context, key string, fn func() (interface{}, error)) (interface{}, error) {
	result, err, shared := s.group.Do(key, fn)
	if shared {
		s.sharedInMinuteCount.Add(1)
		// ...
	}

	return result, err
}

func (s *SingleFlightClient) Get(ctx context.Context, key string) (cmd *redis.StringCmd) {
	_, err := s.do(ctx, key, func() (interface{}, error) {
		cmd = s.UniversalClient.Get(ctx, key)

		return cmd, cmd.Err()
	})
	if cmd == nil {
		cmd = redis.NewStringCmd(ctx)
	}
	if err != nil {
		cmd.SetErr(err)
	}

	return cmd
}

func (s *SingleFlightClient) HGet(ctx context.Context, key, field string) (cmd *redis.StringCmd) {
	_, err := s.do(ctx, key+field, func() (interface{}, error) {
		cmd = s.UniversalClient.HGet(ctx, key, field)

		return cmd, cmd.Err()
	})
	if cmd == nil {
		cmd = redis.NewStringCmd(ctx)
	}
	if err != nil {
		cmd.SetErr(err)
	}

	return cmd
}

func (s *SingleFlightClient) GetDel(ctx context.Context, key string) (cmd *redis.StringCmd) {
	_, err := s.do(ctx, key, func() (_ interface{}, err error) {
		cmd = s.UniversalClient.GetDel(ctx, key)

		return cmd, cmd.Err()
	})
	if cmd == nil {
		cmd = redis.NewStringCmd(ctx)
	}
	if err != nil {
		cmd.SetErr(err)
	}

	return cmd
}

func (s *SingleFlightClient) HMGet(ctx context.Context, key string, field ...string) (cmd *redis.SliceCmd) {
	_, err := s.do(ctx, key+fieldsToKey(field...), func() (interface{}, error) {
		cmd = s.UniversalClient.HMGet(ctx, key, field...)

		return cmd, cmd.Err()
	})
	if cmd == nil {
		cmd = redis.NewSliceCmd(ctx)
	}
	if err != nil {
		cmd.SetErr(err)
	}

	return cmd
}
