package wredis

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	"github.com/stretchr/testify/suite"
)

type TestClientSuite struct {
	suite.Suite

	ctx context.Context

	cntSingle  *atomic.Int64
	cntLru     *atomic.Int64
	addedToLru *atomic.Int64
	clients    map[string]UniversalClient
}

func TestClient(t *testing.T) {
	suite.Run(t, new(TestClientSuite))
}

func (t *TestClientSuite) SetupTest() {
	ctx := context.Background()
	t.ctx = ctx
	t.cntSingle = new(atomic.Int64)
	t.cntLru = new(atomic.Int64)
	t.addedToLru = new(atomic.Int64)

	var (
		cfg    = NewBuildConfig()
		lruCfg = NewLRUConfig().
			WithTTL(time.Minute).
			WithCacheErrors(true, time.Minute).
			WithOnAdd(func(key string, _ any) {
				t.addedToLru.Add(1)
			}).
			WithOnGet(func(key string, _ any) {
				fmt.Println("CNT_LRU", t.cntLru.Add(1))
			})
		singleCfg = SingleFlightConfigs{
			WithCounter(t.cntSingle),
		}
	)

	t.clients = make(map[string]UniversalClient)
	t.clients[basicClient] = NewFactory(cfg).Build(ctx, Options{
		Addr: "localhost:6379",
	})
	t.clients[lruClient] = NewFactory(
		cfg.WithLRU(lruCfg),
		LRU,
	).Build(ctx, Options{
		Addr: "localhost:6379",
	})
	t.clients[lruSingleFlightClient] = NewFactory(
		cfg.WithLRU(lruCfg).WithSingleFlight(singleCfg...),
		LRU, SingleFlight,
	).Build(ctx, Options{
		Addr: "localhost:6379",
	})
	t.clients[singleFlightClient] = NewFactory(
		cfg.WithSingleFlight(singleCfg...),
		SingleFlight,
	).Build(ctx, Options{
		Addr: "localhost:6379",
	})
}

const (
	basicClient           = "basicClient"
	lruSingleFlightClient = "lruSingleFlightClient"
	singleFlightClient    = "singleFlightClient"
	lruClient             = "lruClient"
)

func (t *TestClientSuite) TestClients() {
	const (
		keyNotExists = "keyNotExist"
		keyExists    = "keyExists"
		hMapKeys     = "hMapKeys"
		repeat       = 1_000
	)
	// preconditions
	{
		t.NoError(t.clients[basicClient].Set(t.ctx, keyExists, keyExists, time.Minute).Err())
		t.NoError(t.clients[basicClient].HSet(t.ctx, hMapKeys, hMapKeys, hMapKeys).Err())
	}
	tableTest := []struct {
		name   string
		client string
		action func(client UniversalClient)
		assert func()
	}{
		{
			name:   "not exist",
			client: basicClient,
			action: func(client UniversalClient) {
				_, err := client.Get(t.ctx, keyNotExists).Result()
				t.ErrorIs(err, redis.Nil)
			},
		},
		{
			name:   "not exist",
			client: lruClient,
			action: func(client UniversalClient) {
				_, err := client.Get(t.ctx, keyNotExists).Result()
				t.ErrorIs(err, redis.Nil)
			},
			assert: func() {
				t.True(int64(1) <= t.addedToLru.Swap(0))
				t.True(
					repeat-1 <= t.cntLru.Load(),
				)
			},
		},
		{
			name:   "not exist",
			client: singleFlightClient,
			action: func(client UniversalClient) {
				_, err := client.Get(t.ctx, keyNotExists).Result()
				t.ErrorIs(err, redis.Nil)
			},
			assert: func() {
				t.True(int64(repeat*0.8) <= t.cntSingle.Swap(0))
			},
		},
		{
			name:   "not exist",
			client: lruSingleFlightClient,
			action: func(client UniversalClient) {
				_, err := client.Get(t.ctx, keyNotExists).Result()
				t.ErrorIs(err, redis.Nil)
			},
			assert: func() {
				t.True(
					int64(1) <= t.addedToLru.Swap(0),
				)
				t.True(
					repeat-1 <= t.cntLru.Load(),
				)
				t.True(
					int64(repeat*0.8) <= t.cntSingle.Swap(0),
				)
			},
		},
		{
			name:   "exits",
			client: basicClient,
			action: func(client UniversalClient) {
				result, err := client.Get(t.ctx, keyExists).Result()
				t.NoError(err)
				t.Equal(keyExists, result)
			},
		},
		{
			name:   "exits",
			client: lruClient,
			action: func(client UniversalClient) {
				result, err := client.Get(t.ctx, keyExists).Result()
				t.NoError(err)
				t.Equal(keyExists, result)
			},
			assert: func() {
				t.True(
					int64(1) <= t.addedToLru.Swap(0),
				)
				t.True(
					repeat-1 <= t.cntLru.Load(),
				)
			},
		},
		{
			name:   "exits",
			client: singleFlightClient,
			action: func(client UniversalClient) {
				result, err := client.Get(t.ctx, keyExists).Result()
				t.NoError(err)
				t.Equal(keyExists, result)
			},
			assert: func() {
				t.True(
					int64(repeat*0.8) <= t.cntSingle.Swap(0),
				)
			},
		},
		{
			name:   "exits",
			client: lruSingleFlightClient,
			action: func(client UniversalClient) {
				result, err := client.Get(t.ctx, keyExists).Result()
				t.NoError(err)
				t.Equal(keyExists, result)
			},
			assert: func() {
				t.True(
					int64(1) <= t.addedToLru.Swap(0),
				)
				t.True(
					repeat-t.cntSingle.Load()-1 <= t.cntLru.Swap(0),
				)
				t.True(
					int64(repeat*0.8) <= t.cntSingle.Swap(0),
				)
			},
		},
		{
			name:   "hmap not exist",
			client: basicClient,
			action: func(client UniversalClient) {
				_, err := client.HGet(t.ctx, hMapKeys, hMapKeys+keyNotExists).Result()
				t.ErrorIs(err, redis.Nil)
			},
		},
		{
			name:   "hmap not exist",
			client: lruClient,
			action: func(client UniversalClient) {
				_, err := client.HGet(t.ctx, hMapKeys, hMapKeys+keyNotExists).Result()
				t.ErrorIs(err, redis.Nil)
			},
			assert: func() {
				t.True(
					int64(1) <= t.addedToLru.Swap(0),
				)
				t.True(
					repeat-t.cntSingle.Load()-1 <= t.cntLru.Swap(0),
				)
			},
		},
		{
			name:   "hmap not exist",
			client: singleFlightClient,
			action: func(client UniversalClient) {
				_, err := client.HGet(t.ctx, hMapKeys, hMapKeys+keyNotExists).Result()
				t.ErrorIs(err, redis.Nil)
			},
			assert: func() {
				t.True(
					int64(repeat*0.8) <= t.cntSingle.Swap(0),
				)
			},
		},
		{
			name:   "hmap not exist",
			client: lruSingleFlightClient,
			action: func(client UniversalClient) {
				_, err := client.HGet(t.ctx, hMapKeys, hMapKeys+keyNotExists).Result()
				t.ErrorIs(err, redis.Nil)
			},
			assert: func() {
				t.True(
					int64(1) <= t.addedToLru.Swap(0),
				)
				t.True(
					int64(repeat*0.8) <= t.cntSingle.Swap(0),
				)
			},
		},
		{
			name:   "exits",
			client: basicClient,
			action: func(client UniversalClient) {
				result, err := client.HGet(t.ctx, hMapKeys, hMapKeys).Result()
				t.NoError(err)
				t.Equal(result, hMapKeys)
			},
		},
		{
			name:   "hmap exist",
			client: lruClient,
			action: func(client UniversalClient) {
				result, err := client.HGet(t.ctx, hMapKeys, hMapKeys).Result()
				t.NoError(err)
				t.Equal(result, hMapKeys)
			},
			assert: func() {
				t.True(
					int64(1) <= t.addedToLru.Swap(0),
				)
				t.True(
					int64(repeat*0.8) <= t.cntSingle.Swap(0),
				)
			},
		},
		{
			name:   "hmap exist",
			client: singleFlightClient,
			action: func(client UniversalClient) {
				result, err := client.HGet(t.ctx, hMapKeys, hMapKeys).Result()
				t.NoError(err)
				t.Equal(result, hMapKeys)
			},
			assert: func() {
				t.True(
					int64(repeat*0.8) <= t.cntSingle.Swap(0),
				)
			},
		},
		{
			name:   "hmap exist",
			client: lruSingleFlightClient,
			action: func(client UniversalClient) {
				result, err := client.HGet(t.ctx, hMapKeys, hMapKeys).Result()
				t.NoError(err)
				t.Equal(result, hMapKeys)
			},
			assert: func() {
				t.True(
					int64(repeat*0.8) <= t.cntSingle.Swap(0),
				)
			},
		},
	}
	for _, test := range tableTest {
		test := test
		client := t.clients[test.client]

		t.Run(
			fmt.Sprintf("%s %s count %d %T", test.name, test.client, repeat, client),
			func() {
				wg := errgroup.Group{}
				for i := 0; i < repeat; i++ {
					wg.Go(func() error {
						test.action(client)
						return nil
					})
				}
				t.NoError(wg.Wait())
				if test.assert != nil {
					test.assert()
				}
			},
		)

	}
}
