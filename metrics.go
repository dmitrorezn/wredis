package wredis

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/redis/go-redis/v9"
)

type MetricsClient struct {
	UniversalClient
	meter *Meter
}

func NewMetricsClient(
	name string,
	client UniversalClient,
	metricProvider metric.MeterProvider,
) *MetricsClient {
	return &MetricsClient{
		UniversalClient: client,
		meter:           NewMeter(name, metricProvider.Meter("wredis")),
	}
}

type Meter struct {
	name string

	cmu      sync.RWMutex
	counters map[string]metric.Int64Counter

	hmu        sync.RWMutex
	histograms map[string]metric.Float64Histogram

	gauge map[string]atomic.Int64

	meter metric.Meter
}

func NewMeter(name string, meter metric.Meter) *Meter {
	return &Meter{
		name:  name,
		meter: meter,
	}
}

func (m *Meter) Observe(ctx context.Context, cmd string) func(err error) {
	var (
		start   = time.Now()
		attrOpt = metric.WithAttributes(
			attribute.String("name", m.name),
			attribute.String("cmd", cmd),
		)
	)
	m.cmu.RLock()
	cnt, ok := m.counters[cmd]
	m.cmu.RUnlock()
	var err error
	if !ok {
		if cnt, err = m.meter.Int64Counter(cmd); err == nil {
			m.cmu.Lock()
			m.counters[cmd] = cnt
			m.cmu.Unlock()
		}
	}

	gauge, ok := m.gauge[cmd]
	if !ok {
		_, _ = m.meter.Int64ObservableGauge(cmd, metric.WithInt64Callback(func(ctx context.Context, observer metric.Int64Observer) error {
			observer.Observe(gauge.Load(), attrOpt)
			return nil
		}))
	}
	gauge.Add(1)

	m.hmu.RLock()
	histogram, ok := m.histograms[cmd]
	m.hmu.RUnlock()
	if !ok {
		if histogram, err = m.meter.Float64Histogram(cmd); err == nil {
			m.cmu.Lock()
			m.histograms[cmd] = histogram
			m.cmu.Unlock()
		}
	}

	return func(err error) {
		cnt.Add(ctx, 1, attrOpt, metric.WithAttributes(attribute.Bool("error", err != nil)))

		histogram.Record(
			ctx,
			time.Since(start).Seconds(),
			attrOpt,
			metric.WithAttributes(
				attribute.Bool("error", err != nil),
			),
		)

		gauge.Add(-1)
	}
}

func (s *MetricsClient) Get(ctx context.Context, key string) (cmd *redis.StringCmd) {
	defer s.meter.Observe(ctx, "Get")(cmd.Err())

	return s.UniversalClient.Get(ctx, key)
}

func (s *MetricsClient) HGet(ctx context.Context, key, field string) (cmd *redis.StringCmd) {
	defer s.meter.Observe(ctx, "HGet")(cmd.Err())

	return s.UniversalClient.HGet(ctx, key, field)
}

func (s *MetricsClient) GetDel(ctx context.Context, key string) (cmd *redis.StringCmd) {
	defer s.meter.Observe(ctx, "GetDel")(cmd.Err())

	return s.UniversalClient.GetDel(ctx, key)
}

func (s *MetricsClient) HMGet(ctx context.Context, key string, field ...string) (cmd *redis.SliceCmd) {
	defer s.meter.Observe(ctx, "HMGet")(cmd.Err())

	return s.UniversalClient.HMGet(ctx, key, field...)
}

func (s *MetricsClient) HGetAll(ctx context.Context, key string) (cmd *redis.MapStringStringCmd) {
	defer s.meter.Observe(ctx, "HGetAll")(cmd.Err())

	return s.UniversalClient.HGetAll(ctx, key)
}

/*
	TO DO
	....
*/
