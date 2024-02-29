package wredis

import "github.com/redis/go-redis/v9"

type (
	UniversalClient redis.UniversalClient
	Client          struct {
		UniversalClient
	}
)

type iClient interface {
	AddHook(hook redis.Hook)
	UniversalClient
}

type Configuration func(client iClient)

func WithHooks(hooks ...redis.Hook) Configuration {
	return func(client iClient) {
		for _, hook := range hooks {
			client.AddHook(hook)
		}
	}
}

func New(options Options, configurations ...Configuration) *Client {
	var (
		opts   = redis.Options(options)
		client = redis.NewClient(&opts)
	)
	for _, c := range configurations {
		c(client)
	}

	return &Client{
		UniversalClient: client,
	}
}

func NewCluster(options ClusterOptions, configurations ...Configuration) *Client {
	var (
		opts   = redis.ClusterOptions(options)
		client = redis.NewClusterClient(&opts)
	)
	for _, c := range configurations {
		c(client)
	}

	return &Client{
		UniversalClient: client,
	}
}
