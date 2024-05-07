package wredis

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	ds := decorators{SingleFlightDecorator, LRUDecorator}
	sort.Sort(ds)

	if ds[0] != LRUDecorator || ds[1] != SingleFlightDecorator {
		t.Error("LRUDecorator should be first in decorators list", ds)
	}
}

func TestUniqCache(t *testing.T) {
	ds := decorators{SingleFlightDecorator, LRUDecorator}

	assert.True(t, ds.UniqCache())
	ds = decorators{SingleFlightDecorator, LRUDecorator, LFUDecorator}

	assert.False(t, ds.UniqCache())
}
