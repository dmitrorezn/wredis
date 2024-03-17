package wredis

import (
	"fmt"
	"strings"
)

type (
	Event  string
	Events []Event
)

func (e Events) Empty() bool {
	return len(e) == 0
}

func (e Events) String() string {
	events := make([]string, len(e))
	for i, ev := range e {
		events[i] = string(ev)
	}

	return strings.Join(events, ",")
}

const (
	DelEvent    Event = "deleted"
	ExpireEvent Event = "expired"
	SetEvent    Event = "set"
)

type Pattern string

const keysPatternChannel Pattern = "__key*__:%s"

func (p Pattern) String(s ...any) string {
	return fmt.Sprintf(string(p), s...)
}
