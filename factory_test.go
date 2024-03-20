package wredis

import (
	"context"
	"fmt"
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
		cfg    = NewBuildConfig()
		lruCfg = NewLRUConfig().
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

type fifo struct {
	items chan string
}

func (f *fifo) Push(item string) {
	select {
	case f.items <- item:
	default:
	}
}

func (f *fifo) Pop() string {
	select {
	case it := <-f.items:
		return it
	default:
		return "default"
	}
}

func (f *fifo) RepeatPush(item string, n int) {
	for i := 0; i < n; i++ {
		f.items <- item
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

func BenchmarkBasicClient(b *testing.B) {
	t := new(TestClientSuite)
	t.SetupTest()

	ctx, cancel := context.WithTimeout(t.ctx, 10*time.Minute)
	defer cancel()

	t.initBench(ctx)

	tableTest := []struct {
		op     string
		client string
	}{
		{
			client: basicClient,
		},
		//{
		//	client: lruClient,
		//},
		//{
		//	client: singleFlightClient,
		//},
		//{
		//	client: lruSingleFlightClient,
		//},
	}

	one := sync.Once{}
	keys := make([]string, maxkeyixdss)

	for _, test := range tableTest {
		test := test
		client := t.clients[test.client]

		b.Run(test.op+" "+test.client, func(b *testing.B) {
			defer func() {
				fmt.Println("FINISHED", test.op+" "+test.client)
			}()
			b.N = maxkeyixdss / 2
			one.Do(func() {
				for i := 0; i < maxkeyixdss; i++ {
					keys[i] = fmt.Sprint(rand.Intn(maxkeyixd))
				}
			})
			for i := b.N; i > 0; i-- {
				b.StopTimer()
				key := ""
				if i < len(keys) {
					key = test.op + keys[i]
				} else {
					key = fmt.Sprint(test.op, rand.Intn(maxkeyixd))
				}
				b.StartTimer()

				_, _ = client.Get(t.ctx, key).Result()
			}
		})
	}
}

func BenchmarkClientLruClient(b *testing.B) {
	t := new(TestClientSuite)
	t.SetupTest()

	ctx, cancel := context.WithTimeout(t.ctx, 10*time.Minute)
	defer cancel()

	t.initBench(ctx)

	tableTest := []struct {
		op     string
		client string
	}{
		{
			client: lruClient,
		},
	}

	one := sync.Once{}
	keys := make([]string, maxkeyixdss)

	for _, test := range tableTest {
		test := test
		client := t.clients[test.client]

		b.Run(test.op+" "+test.client, func(b *testing.B) {
			defer func() {
				fmt.Println("FINISHED", test.op+" "+test.client)
			}()
			b.N = maxkeyixdss / 2
			one.Do(func() {
				for i := 0; i < maxkeyixdss; i++ {
					keys[i] = fmt.Sprint(rand.Intn(maxkeyixd))
				}
			})
			for i := b.N; i > 0; i-- {
				b.StopTimer()
				key := ""
				if i < len(keys) {
					key = test.op + keys[i]
				} else {
					key = fmt.Sprint(test.op, rand.Intn(maxkeyixd))
				}
				b.StartTimer()

				_, _ = client.Get(t.ctx, key).Result()
			}
		})
	}
}

func BenchmarkClientHmap(b *testing.B) {
	t := new(TestClientSuite)
	t.SetupTest()

	const (
		hmap      = "hmap"
		maxkeyixd = 1000
	)

	ctx, cancel := context.WithTimeout(t.ctx, time.Hour)
	defer cancel()

	for i := 0; i < maxkeyixd; i++ {
		if i%2 == 0 {
			t.clients[basicClient].Set(ctx, fmt.Sprint(i), fmt.Sprint(i), time.Hour)
			hkey := fmt.Sprint(hmap, i)
			t.clients[basicClient].HSet(ctx, hkey, hkey, hkey)
		}
	}

	tableTest := []struct {
		op     string
		client string
	}{
		{
			op:     hmap,
			client: basicClient,
		},
		{
			op:     hmap,
			client: lruClient,
		},
		{
			op:     hmap,
			client: singleFlightClient,
		},
		{
			op:     hmap,
			client: lruSingleFlightClient,
		},
	}

	one := sync.Once{}
	keys := make([]string, b.N)

	for _, test := range tableTest {
		test := test
		client := t.clients[test.client]

		b.Run(test.op+" "+test.client, func(b *testing.B) {
			one.Do(func() {
				for i := 0; i < b.N; i++ {
					keys[i] = fmt.Sprint(rand.Intn(maxkeyixd))
				}
			})

			for i := b.N; i > 0; i-- {
				b.StopTimer()
				key := ""
				if i < len(keys) {
					key = test.op + keys[i]
				} else {
					key = fmt.Sprint(test.op, rand.Intn(maxkeyixd))
				}
				if i%1000 == 0 {
					fmt.Println("key", key)
				}
				b.StartTimer()

				_, _ = client.HGet(t.ctx, key, key).Result()
			}
		})
	}
}
