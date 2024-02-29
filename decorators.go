package wredis

import "sort"

// decoratorID value is the order in decoration chain
// less decorator value means that its usage is closer to original wrapped realization
//
//go:generate stringer -type=decoratorID
type decoratorID int

type decorators []decoratorID

const (
	LRU decoratorID = iota + 1
	SingleFlight
)

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
