package wredis

import redis "github.com/redis/go-redis/v9"

type IOptions interface {
	opt()
}

type Options redis.Options

func (Options) opt() {}

func (o Options) WithAddr(network string) {
	o.Addr = network
}

type ClusterOptions redis.ClusterOptions

func (ClusterOptions) opt() {}

func (o ClusterOptions) WithAddrs(addrs []string) {
	o.Addrs = addrs
}
