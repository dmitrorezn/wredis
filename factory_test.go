package wredis

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
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
		cfg      = NewBuildConfig()
		cacheCfg = NewLocalCacheConfig().
				WithTTL(time.Minute).
				WithCacheErrors(true, 30*time.Second).
				WithOnAdd(func(key string, _ any) {
				t.addedToLru.Add(1)
			}).
			WithSize(10_000).
			WithOnGet(func(key string, _ any) {
				t.cntLru.Add(1)
			})
		singleCfg = SingleFlightConfigs{
			WithCounter(t.cntSingle),
		}
	)

	t.clients = make(map[string]UniversalClient)
	var err error
	t.clients[basicClient], err = NewFactory(cfg).Build(ctx, Options{
		Addr: "localhost:6379",
	})
	assert.NoError(t.T(), err)
	t.clients[lruClient], err = NewFactory(
		cfg.WithCacheConfig(cacheCfg),
		LRUDecorator,
	).Build(ctx, Options{
		Addr: "localhost:6379",
	})
	assert.NoError(t.T(), err)

	t.clients[lfuClient], err = NewFactory(
		cfg.WithCacheConfig(cacheCfg.WithSamples(100_000)),
		LFUDecorator,
	).Build(ctx, Options{
		Addr: "localhost:6379",
	})
	assert.NoError(t.T(), err)

	t.clients[lruSingleFlightClient], err = NewFactory(
		cfg.WithCacheConfig(cacheCfg).WithSingleFlight(singleCfg...),
		LRUDecorator, SingleFlightDecorator,
	).Build(ctx, Options{
		Addr: "localhost:6379",
	})
	assert.NoError(t.T(), err)

	t.clients[lfuSingleFlightClient], err = NewFactory(
		cfg.WithCacheConfig(cacheCfg.WithSamples(100_000)).WithSingleFlight(singleCfg...),
		LFUDecorator, SingleFlightDecorator,
	).Build(ctx, Options{
		Addr: "localhost:6379",
	})
	assert.NoError(t.T(), err)

	t.clients[singleFlightClient], err = NewFactory(
		cfg.WithSingleFlight(singleCfg...),
		SingleFlightDecorator,
	).Build(ctx, Options{
		Addr: "localhost:6379",
	})
	assert.NoError(t.T(), err)

}

const (
	basicClient           = "basicClient"
	lruSingleFlightClient = "lruSingleFlightClient"
	lfuSingleFlightClient = "lfuSingleFlightClient"
	singleFlightClient    = "singleFlightClient"
	lruClient             = "lruClient"
	lfuClient             = "lfuClient"
)

func (t *TestClientSuite) TestClients() {
	const (
		keyNotExists  = "keyNotExist"
		keyExists     = "keyExists"
		hMapKeys      = "hMapKeys"
		hMapKeysKey   = "h_key"
		hMapKeysValue = "h_value"
		repeat        = 10_000
	)
	// preconditions
	{
		t.NoError(t.clients[basicClient].Set(t.ctx, keyExists, keyExists, time.Minute).Err())
		t.NoError(t.clients[basicClient].HSet(t.ctx, hMapKeys, hMapKeysKey, hMapKeysValue).Err())
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
				t.True(t.addedToLru.Swap(0) > 0)
				t.True(t.cntLru.Swap(0) > repeat*0.3)
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
				t.True(t.cntSingle.Swap(0) > repeat*0.8)
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
				t.True(t.addedToLru.Swap(0) >= 1)
				t.True(t.cntLru.Swap(0) >= repeat*0.3)
				t.True(t.cntSingle.Swap(0) >= repeat*0.4)
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
				t.True(t.addedToLru.Swap(0) >= 1)
				t.True(t.cntLru.Swap(0) >= repeat*0.3)
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
				t.True(t.cntSingle.Swap(0) >= repeat*0.4)
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
				t.True(t.addedToLru.Swap(0) >= 1)
				t.True(t.cntLru.Swap(0) >= repeat*0.3)
				t.True(t.cntSingle.Swap(0) >= repeat*0.4)
			},
		},
		{
			name:   "hmap not exist",
			client: basicClient,
			action: func(client UniversalClient) {
				_, err := client.HGet(t.ctx, hMapKeys, hMapKeysKey+keyNotExists).Result()
				t.ErrorIs(err, redis.Nil)
			},
		},
		{
			name:   "hmap not exist",
			client: lruClient,
			action: func(client UniversalClient) {
				_, err := client.HGet(t.ctx, hMapKeys, hMapKeysKey+keyNotExists).Result()
				t.ErrorIs(err, redis.Nil)
			},
			assert: func() {
				t.True(
					t.addedToLru.Swap(0) >= 1,
				)
				t.True(
					t.cntLru.Swap(0) >= repeat*0.4,
				)
			},
		},
		{
			name:   "hmap not exist",
			client: singleFlightClient,
			action: func(client UniversalClient) {
				_, err := client.HGet(t.ctx, hMapKeys, hMapKeysKey+keyNotExists).Result()
				t.ErrorIs(err, redis.Nil)
			},
			assert: func() {
				t.True(
					t.cntSingle.Swap(0) >= repeat*0.4,
				)
			},
		},
		{
			name:   "hmap not exist",
			client: lruSingleFlightClient,
			action: func(client UniversalClient) {
				_, err := client.HGet(t.ctx, hMapKeys, hMapKeysKey+keyNotExists).Result()
				t.ErrorIs(err, redis.Nil)
			},
			assert: func() {
				t.True(
					t.addedToLru.Swap(0) >= 1,
				)
				t.True(
					t.cntLru.Swap(0) >= repeat*0.3,
				)
				t.True(
					t.cntSingle.Swap(0) >= repeat*0.4,
				)
			},
		},
		{
			name:   "exits",
			client: basicClient,
			action: func(client UniversalClient) {
				result, err := client.HGet(t.ctx, hMapKeys, hMapKeysKey).Result()
				t.NoError(err)
				t.Equal(result, hMapKeysValue)
			},
		},
		{
			name:   "hmap exist",
			client: lruClient,
			action: func(client UniversalClient) {
				result, err := client.HGet(t.ctx, hMapKeys, hMapKeysKey).Result()
				t.NoError(err)
				t.Equal(result, hMapKeysValue)
			},
			assert: func() {
				t.True(
					t.addedToLru.Swap(0) >= 1,
				)
				t.True(
					t.cntLru.Swap(0) >= repeat*0.3,
				)
			},
		},
		{
			name:   "hmap exist",
			client: singleFlightClient,
			action: func(client UniversalClient) {
				result, err := client.HGet(t.ctx, hMapKeys, hMapKeysKey).Result()
				t.NoError(err)
				t.Equal(result, hMapKeysValue)
			},
			assert: func() {
				t.True(
					t.cntSingle.Swap(0) >= repeat*0.4,
				)
			},
		},
		{
			name:   "hmap exist",
			client: lruSingleFlightClient,
			action: func(client UniversalClient) {
				result, err := client.HGet(t.ctx, hMapKeys, hMapKeysKey).Result()
				t.NoError(err)
				t.Equal(result, hMapKeysValue)
			},
			assert: func() {
				t.True(
					t.addedToLru.Swap(0) >= 1,
				)
				t.True(
					t.cntLru.Swap(0) >= repeat*0.3,
				)
				t.True(
					t.cntSingle.Swap(0) >= repeat*0.4,
				)
			},
		},
	}
	for _, test := range tableTest {
		test := test
		client := t.clients[test.client]

		t.Run(
			fmt.Sprintf("%s %s count %d", test.name, test.client, repeat),
			func() {
				wg := errgroup.Group{}
				for i := 0; i < repeat; i++ {
					if i%2 == 0 {
						test.action(client)
						continue
					}
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

const (
	hmap        = "hmap"
	maxkeyixd   = 1000
	maxkeyixdss = 10000
)

func (t *TestClientSuite) initBench(ctx context.Context) {
	for i := 0; i < maxkeyixd; i++ {
		if i%2 == 0 {
			t.clients[basicClient].Set(ctx, fmt.Sprint(i), fmt.Sprint(i), 10*time.Minute)
			hkey := fmt.Sprint(hmap, i)
			t.clients[basicClient].HSet(ctx, hkey, hkey, hkey)
		}
	}
}

/*
goos: darwin
goarch: arm64
pkg: github.com/dmitrorezn/wredis
BenchmarkBasicClient
BenchmarkBasicClient/Get
BenchmarkBasicClient/Get-10         	   81182	     15701 ns/op
BenchmarkBasicClient/HGet
BenchmarkBasicClient/HGet-10        	   65368	     16395 ns/op
*/
func BenchmarkBasicClient(b *testing.B) {
	t := new(TestClientSuite)
	b.ReportAllocs()

	b.Run("Get", func(b *testing.B) {
		benchGet(t, basicClient, b)
	})
	b.Run("HGet", func(b *testing.B) {
		benchHGet(t, basicClient, b)
	})
}

/*
goos: darwin
goarch: arm64
pkg: github.com/dmitrorezn/wredis
BenchmarkClientLruClient
BenchmarkClientLruClient/Get
BenchmarkClientLruClient/Get-10         	  137955	      7268 ns/op
BenchmarkClientLruClient/HGet
BenchmarkClientLruClient/HGet-10        	  134301	      7928 ns/op
*/
func BenchmarkClientLruClient(b *testing.B) {
	t := new(TestClientSuite)
	b.ReportAllocs()

	b.Run("Get", func(b *testing.B) {
		benchGet(t, lruClient, b)
	})
	b.Run("HGet", func(b *testing.B) {
		benchHGet(t, lruClient, b)
	})
}

/*
goos: darwin
goarch: arm64
pkg: github.com/dmitrorezn/wredis
BenchmarkClientLfuClient
BenchmarkClientLfuClient/Get
BenchmarkClientLfuClient/Get-10         	  389989	      2567 ns/op
BenchmarkClientLfuClient/HGet
BenchmarkClientLfuClient/HGet-10        	  815414	      1234 ns/op
*/
func BenchmarkClientLfuClient(b *testing.B) {
	t := new(TestClientSuite)
	b.ReportAllocs()

	b.Run("Get", func(b *testing.B) {
		benchGet(t, lfuClient, b)
	})
	b.Run("HGet", func(b *testing.B) {
		benchHGet(t, lfuClient, b)
	})
}

/*
goos: darwin
goarch: arm64
pkg: github.com/dmitrorezn/wredis
BenchmarkSingleFlightClient
BenchmarkSingleFlightClient/Get
BenchmarkSingleFlightClient/Get-10         	   56251	     21283 ns/op
BenchmarkSingleFlightClient/HGet
BenchmarkSingleFlightClient/HGet-10        	   41181	     26600 ns/op
*/
func BenchmarkSingleFlightClient(b *testing.B) {
	t := new(TestClientSuite)
	b.ReportAllocs()

	b.Run("Get", func(b *testing.B) {
		benchGet(t, singleFlightClient, b)
	})
	b.Run("HGet", func(b *testing.B) {
		benchHGet(t, singleFlightClient, b)
	})
}

/*
goos: darwin
goarch: arm64
pkg: github.com/dmitrorezn/wredis
BenchmarkLruSingleFlightClient
BenchmarkLruSingleFlightClient/Get
BenchmarkLruSingleFlightClient/Get-10         	  224241	      5576 ns/op
BenchmarkLruSingleFlightClient/HGet
BenchmarkLruSingleFlightClient/HGet-10        	  115236	      9987 ns/op
*/
func BenchmarkLruSingleFlightClient(b *testing.B) {
	t := new(TestClientSuite)
	b.ReportAllocs()

	b.Run("Get", func(b *testing.B) {
		benchGet(t, lruSingleFlightClient, b)
	})
	b.Run("HGet", func(b *testing.B) {
		benchHGet(t, lruSingleFlightClient, b)
	})
}

/*
goos: darwin
goarch: arm64
pkg: github.com/dmitrorezn/wredis
BenchmarkLfuSingleFlightClient
BenchmarkLfuSingleFlightClient/Get
BenchmarkLfuSingleFlightClient/Get-10         	  177110	      9872 ns/op
BenchmarkLfuSingleFlightClient/HGet
BenchmarkLfuSingleFlightClient/HGet-10        	  181047	      5816 ns/op
*/
func BenchmarkLfuSingleFlightClient(b *testing.B) {
	t := new(TestClientSuite)
	b.ReportAllocs()

	b.Run("Get", func(b *testing.B) {
		benchGet(t, lfuSingleFlightClient, b)
	})
	b.Run("HGet", func(b *testing.B) {
		benchHGet(t, lfuSingleFlightClient, b)
	})
}

func benchGet(t *TestClientSuite, clientName string, b *testing.B) {
	b.Helper()
	t.SetupTest()

	ctx, cancel := context.WithTimeout(t.ctx, 10*time.Minute)
	defer cancel()

	t.initBench(ctx)

	one := sync.Once{}
	keys := make([]string, maxkeyixdss)
	one.Do(func() {
		for i := 0; i < maxkeyixdss; i++ {
			keys[i] = fmt.Sprint(rand.Intn(maxkeyixd))
		}
	})

	client := t.clients[clientName]
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			b.StopTimer()
			key := ""
			if i < len(keys) {
				key = keys[i]
			} else {
				key = fmt.Sprint(rand.Intn(maxkeyixd))
			}
			b.StartTimer()

			_, _ = client.Get(t.ctx, key).Result()
			i++
		}
	})
}

func benchHGet(t *TestClientSuite, clientName string, b *testing.B) {
	b.Helper()
	t.SetupTest()

	ctx, cancel := context.WithTimeout(t.ctx, 10*time.Minute)
	defer cancel()

	t.initBench(ctx)

	one := sync.Once{}
	keys := make([]string, maxkeyixdss)
	one.Do(func() {
		for i := 0; i < maxkeyixdss; i++ {
			keys[i] = fmt.Sprint(rand.Intn(maxkeyixd))
		}
	})

	client := t.clients[clientName]
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			b.StopTimer()
			key := ""
			if i < len(keys) {
				key = keys[i]
			} else {
				key = fmt.Sprint(rand.Intn(maxkeyixd))
			}
			b.StartTimer()

			_, _ = client.HGet(t.ctx, key, key).Result()
			i++
		}
	})
}
