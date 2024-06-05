package wredis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dmitrorezn/wredis/pkg/logger"
	"go.uber.org/zap"
	"strings"
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
		name:       name,
		meter:      meter,
		counters:   make(map[string]metric.Int64Counter),
		histograms: make(map[string]metric.Float64Histogram),
		gauge:      make(map[string]atomic.Int64),
	}
}

func (m *Meter) Observe(ctx context.Context, cmd string) func(err error) {
	var (
		start   = time.Now()
		attrOpt = metric.WithAttributes(
			attribute.String("name", m.name),
			attribute.String("cmd", cmd),
		)
		metricName = strings.ToLower(
			fmt.Sprintf("redis_%s", cmd),
		)
	)
	m.cmu.RLock()
	cnt, ok := m.counters[metricName]
	m.cmu.RUnlock()
	var err error
	if !ok {
		if cnt, err = m.meter.Int64Counter(metricName + "_counter"); err == nil {
			m.cmu.Lock()
			m.counters[metricName] = cnt
			m.cmu.Unlock()
		}
	}

	gauge, ok := m.gauge[metricName]
	if !ok {
		_, _ = m.meter.Int64ObservableGauge(metricName+"_gauge", metric.WithInt64Callback(func(ctx context.Context, observer metric.Int64Observer) error {
			observer.Observe(gauge.Load(), attrOpt)
			return nil
		}))
	}
	gauge.Add(1)

	m.hmu.RLock()
	histogram, ok := m.histograms[metricName]
	m.hmu.RUnlock()
	if !ok {
		if histogram, err = m.meter.Float64Histogram(metricName + "_histogram"); err == nil {
			m.cmu.Lock()
			m.histograms[metricName] = histogram
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
	record := s.meter.Observe(ctx, "Get")
	defer func() {
		record(cmd.Err())
	}()

	return s.UniversalClient.Get(ctx, key)
}

func (s *MetricsClient) HGet(ctx context.Context, key, field string) (cmd *redis.StringCmd) {
	record := s.meter.Observe(ctx, "HGet")
	defer func() {
		record(cmd.Err())
	}()

	return s.UniversalClient.HGet(ctx, key, field)
}

func (s *MetricsClient) GetDel(ctx context.Context, key string) (cmd *redis.StringCmd) {
	record := s.meter.Observe(ctx, "GetDel")
	defer func() {
		record(cmd.Err())
	}()

	return s.UniversalClient.GetDel(ctx, key)
}

func (s *MetricsClient) HMGet(ctx context.Context, key string, field ...string) (cmd *redis.SliceCmd) {
	record := s.meter.Observe(ctx, "HMGet")
	defer func() {
		record(cmd.Err())
	}()

	return s.UniversalClient.HMGet(ctx, key, field...)
}

func (s *MetricsClient) HGetAll(ctx context.Context, key string) (cmd *redis.MapStringStringCmd) {
	record := s.meter.Observe(ctx, "HGetAll")
	defer func() {
		record(cmd.Err())
	}()

	return s.UniversalClient.HGetAll(ctx, key)
}

/*
	TO DO
	....
*/

type BufferedCollector struct {
	Collector
	wg sync.WaitGroup

	quit chan struct{}

	cfg MetricConfig

	measurements chan Measurement
	flushTicker  *time.Ticker
	buffer       []Measurement
	bufSize      int
	bufMaxSize   int
}

type MetricConfig struct {
	flushInterval     time.Duration
	bufSize           int
	timeoutSend       time.Duration
	heartbeatInterval time.Duration
}

func WithFlushInterval(dur time.Duration) func(config *MetricConfig) {
	return func(config *MetricConfig) {
		config.flushInterval = dur
	}
}

func WithHeartbeatInterval(dur time.Duration) func(config *MetricConfig) {
	return func(config *MetricConfig) {
		config.heartbeatInterval = dur
	}
}

func WithBufSize(bufSize int) func(config *MetricConfig) {
	return func(config *MetricConfig) {
		config.bufSize = bufSize
	}
}

func WithTimeoutSend(timeoutSend time.Duration) func(config *MetricConfig) {
	return func(config *MetricConfig) {
		config.timeoutSend = timeoutSend
	}
}

func NewCfg(opts ...func(config *MetricConfig)) MetricConfig {
	cfg := MetricConfig{
		flushInterval:     2 * time.Second,
		bufSize:           1000, // RPS: 500/s
		timeoutSend:       50 * time.Millisecond,
		heartbeatInterval: 10 * time.Second,
	}
	for _, op := range opts {
		op(&cfg)
	}

	return cfg
}

type Measurement struct {
	Name      string
	InputSize int64
	RespSize  int64
	Dur       time.Duration
	Missed    bool
	Hit       bool
	HasErr    bool
	Attr      []attribute.KeyValue
}

func createMeasurement(start time.Time, name string, cmd string, inputSize, respSize int64, err error) Measurement {
	return Measurement{
		Dur:       time.Since(start),
		RespSize:  respSize,
		InputSize: inputSize,
		Attr: []attribute.KeyValue{
			attribute.String("cache", name),
			attribute.String("cmd", cmd),
		},
		Missed: err == redis.Nil,
		HasErr: err != nil && err != redis.Nil,
		Hit:    err == nil,
	}
}

func (m Measurement) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

type Collector interface {
	RecordMeasurement(ctx context.Context, m Measurement)
}
type MetricCollector struct {
	counters map[string]metric.Int64Counter
}

const (
	call = "calls"
	miss = "miss"
)

type name string

func (n name) With(add string)name  {
	return n+"_"+name(add)
}
func (n name) String()string  {
	return string(n)
}

func NewMetricCollector(provider metric.MeterProvider) (*MetricCollector, error) {
	meter := provider.Meter("github.com/dmitrorezn/wredis")
var prefix name = "redis_client"
	var err error
	m := &MetricCollector{
		counters: make(map[string]metric.Int64Counter),
	}
	if m.counters[call], err = meter.Int64Counter(prefix.With(call).String()); err != nil {
		return nil, err
	}
	if m.counters[miss], err = meter.Int64Counter(prefix.With(miss).String()); err != nil {
		return nil, err
	}

	return m, nil
}
func (mc *MetricCollector) RecordMeasurement(ctx context.Context, m Measurement) {
	mc.counters[call].Add(ctx, 1)
	if m.Missed{
		mc.counters[miss].Add(ctx, 1)
	}
}

func newBufColl(c Collector, cfg MetricConfig) *BufferedCollector {
	return &BufferedCollector{
		Collector:    c,
		cfg:          cfg,
		quit:         make(chan struct{}),
		buffer:       make([]Measurement, 0, cfg.bufSize),
		measurements: make(chan Measurement, cfg.bufSize),
		flushTicker:  time.NewTicker(cfg.flushInterval),
	}
}

func (bc *BufferedCollector) Start(ctx context.Context) {
	bc.wg.Add(1)
	go func() {
		defer bc.wg.Done()
		bc.collect(ctx)
	}()
}

func (bc *BufferedCollector) Close() {
	close(bc.quit)

	bc.flushTicker.Reset(bc.cfg.flushInterval)

	bc.wg.Wait()

	close(bc.measurements)
}

func (bc *BufferedCollector) collect(ctx context.Context) {
	var (
		heartbeat = time.NewTicker(1 * time.Second)
	)

	for {
		select {
		case <-ctx.Done():
			bc.flushBuffer(ctx)
			return
		case <-bc.quit:
			logger.Info("BufferedCollector collect closed")
			bc.flushBuffer(ctx)
			return
		case <-heartbeat.C:
			logger.Info("BufferedCollector collect alive",
				zap.Any("buf_chan_len", len(bc.measurements)),
				zap.Any("buf_size", bc.bufMaxSize),
			)
		case <-bc.flushTicker.C:
			bc.flushBuffer(ctx)
		case m, opened := <-bc.measurements:
			if !opened {
				logger.Info("BufferedCollector collect closed")

				return
			}

			bc.addToBuffer(m)
		}
	}
}

func (bc *BufferedCollector) addToBuffer(m Measurement) {
	bc.buffer = append(bc.buffer, m)
	bc.bufSize += 1
	if bc.bufSize == bc.bufMaxSize {
		bc.flushTicker.Reset(bc.cfg.flushInterval)
	}
}

func (bc *BufferedCollector) flushBuffer(ctx context.Context) {
	if bc.bufSize == 0 || len(bc.buffer) == 0 {
		return
	}

	buf := make([]Measurement, len(bc.buffer))
	n := copy(buf, bc.buffer)
	bc.bufSize = 0
	bc.buffer = bc.buffer[:0:n] // set capacity as n
	bc.wg.Add(1)
	go func() {
		logger.Info("flushBuffer ", zap.Int("len", len(buf)))

		defer bc.wg.Done()
		for _, mm := range buf {
			bc.Collector.RecordMeasurement(ctx, mm)
		}
	}()
}

func (bc *BufferedCollector) Record(ctx context.Context, name string, cmd string, inputSize int64) func(err error, respSize int64) {
	start := time.Now()

	return func(err error, respSize int64) {
		var measurement = createMeasurement(start, name, cmd, inputSize, respSize, err)

		select {
		case <-ctx.Done():
			ctx = context.Background()
		case <-time.After(bc.cfg.timeoutSend):
			logger.Warn("Record -> timeout push", zap.String("cmd", cmd), zap.String("name", name))
		case <-bc.quit:
			logger.Warn("Record -> push quit chan closed", zap.String("cmd", cmd), zap.String("name", name))
		case bc.measurements <- measurement:
			return
		}

		bc.Collector.RecordMeasurement(ctx, measurement)
	}
}
