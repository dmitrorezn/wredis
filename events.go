package wredis

import (
	"fmt"
	"strings"
)

type (
	Event  string
	Events []Event
)

func (e Events) String() string {
	events := make([]string, len(e))
	for i, ev := range e {
		events[i] = string(ev)
	}

	return strings.Join(events, ",")
}

const (
	DelEvent    Event = "deleted"
	ExpireEvent       = "expired"
	SetEvent          = "set"
)

type Pattern string

const keysPatternChannel Pattern = "__key*__:%s"

func (p Pattern) String(s string) string {
	return fmt.Sprintf(string(p), s)
}
