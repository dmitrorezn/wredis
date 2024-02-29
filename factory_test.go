package wredis

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TestClientSuite struct {
	suite.Suite

	ctx context.Context

	factory *Factory
	client  UniversalClient
}

func TestClient(t *testing.T) {
	suite.Run(t, new(TestClientSuite))
}

func (t *TestClientSuite) SetupTest() {
	ctx := context.Background()
	t.ctx = ctx

	cfg := NewBuildConfig().
		WithLRU(NewLRUConfig())

	t.factory = NewFactory(
		cfg,
		LRU, SingleFlight,
	)
	client, err := miniredis.Run()
	t.NoError(err)

	t.client = t.factory.Build(ctx, Options{
		Addr: client.Addr(),
	})
}

func (t *TestClientSuite) TestGet() {
	const (
		key       = "ket"
		existskey = "existskey"
	)

	err := t.client.Set(t.ctx, existskey, existskey, time.Second).Err()
	t.NoError(err)

	tabletest := []struct {
		assert func(val string, err error)
		input  string
	}{
		{
			assert: func(val string, err error) {
				t.ErrorIs(err, redis.Nil)
				t.Equal("", val)
			},
			input: key,
		},
		{
			assert: func(val string, err error) {
				t.ErrorIs(err, redis.Nil)
				t.Equal("", val)
			},
			input: key,
		},
		{
			assert: func(val string, err error) {
				t.NoError(err)
				t.Equal(existskey, val)
			},
			input: existskey,
		},
		{
			assert: func(val string, err error) {
				t.NoError(err)
				t.Equal(existskey, val)
			},
			input: existskey,
		},
	}
	for _, test := range tabletest {
		test.assert(
			t.client.Get(t.ctx, test.input).Result(),
		)
	}
}
