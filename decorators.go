package wredis

import "sort"

// decoratorID value is the order in decoration chain
// less decorator value means that its usage is closer to original wrapped realization
//
//go:generate stringer -type=decoratorID
type decoratorID int

type decorators []decoratorID

const (
	LRUDecorator decoratorID = iota + 1
	LFUDecorator
	SingleFlightDecorator
	EncodingDecorator
)

var caches = map[decoratorID]int{
	LRUDecorator: 1,
	LFUDecorator: 1,
}

func (ds decorators) UniqCache() bool {
	cachesCount := 0
	for _, id := range ds {
		cachesCount += caches[id]
	}
	return cachesCount <= 1
}

var _ sort.Interface = new(decorators)

func (d decorators) Len() int {
	return len(d)
}

func (d decorators) Less(i, j int) bool {
	return d[i] < d[j]
}

func (d decorators) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
