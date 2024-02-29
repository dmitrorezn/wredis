package wredis

import (
	"sort"
	"strings"
)

func fieldsToKey(field ...string) string {
	sort.Strings(field)

	return strings.Join(field, "")
}
