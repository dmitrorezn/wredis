package wredis

type Configs[T any] []Config[T]

type Config[T any] func(*T)

func (c Configs[T]) Apply(v *T) {
	for _, cc := range c {
		cc(v)
	}
}
