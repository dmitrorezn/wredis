package wredis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/sdk/resource"

	 "github.com/prometheus/client_golang/prometheus"

	"go.opentelemetry.io/otel/sdk/metric"
	"testing"
)

type TestMetricsClientSuite struct {
	suite.Suite

	ctx context.Context

	meter  *Meter
	basic  UniversalClient
	client *MetricsClient
}

func TestMetricsClient(t *testing.T) {
	suite.Run(t, new(TestMetricsClientSuite))
}

type coll struct {

}

func (c coll) Describe(descs chan<- *prometheus.Desc) {
	descs <-&prometheus.Desc{

	}
}

func (c coll) Collect(metrics chan<- prometheus.Metric) {
	metrics <-prometheus.Metric{

	}
}

func (t *TestMetricsClientSuite) SetupTest() {
	t.ctx = context.Background()

	pclient :=prometheus.Register(new(coll))
	pclient.

	res, err := resource.New(t.ctx, resource.WithSchemaURL("http://localhost:9090"))
	t.NoError(err)
	meter := metric.NewMeterProvider(metric.WithResource(res))

	t.basic, err = NewFactory(NewBuildConfig()).Build(t.ctx, Options{
		Addr: "localhost:6379",
	})
	t.NoError(err)

	t.client = NewMetricsClient("text-chache-module", t.basic, meter)
}

func (t *TestMetricsClientSuite) TestGet() {
	err := t.client.Get(t.ctx, "test").Err()
	t.ErrorIs(err, redis.Nil)

	err = t.client.Set(t.ctx, "test", "test", time.Minute).Err()
	t.NoError(err)

	err = t.client.Get(t.ctx, "test").Err()
	t.NoError(err)
}
