package wredis

import (
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	ds := decorators{SingleFlight, LRU}
	sort.Sort(ds)

	if ds[0] != LRU || ds[1] != SingleFlight {
		t.Error("LRU should be first in decorators list", ds)
	}
}
