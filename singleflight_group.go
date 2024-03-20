package wredis

import "golang.org/x/sync/singleflight"

type command[T any] func(client UniversalClient) T

func newGroup[T any](
	g *singleflight.Group,
	client UniversalClient,
) *group[T] {
	return &group[T]{
		group:  g,
		client: client,
	}
}

type group[T any] struct {
	group  *singleflight.Group
	client UniversalClient
}

func (g *group[T]) Do(key string, cmd command[T]) (T, error, bool) {
	val, err, shared := g.group.Do(key, func() (interface{}, error) {
		return cmd(g.client), nil
	})

	return val.(T), err, shared
}
