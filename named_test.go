package wredis

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	var n NamedClient
	ctype := reflect.TypeOf(n)

	b := strings.Builder{}

out:
	for i := 0; i < ctype.NumMethod(); i++ {
		m := ctype.Method(i)

		fn := strings.Replace(m.Func.String(), "<func", "", 1)
		fn = strings.Replace(fn, "Value>", "", 1)
		fn = strings.Replace(fn, "wredis.NamedClient,", "", 1)
		fn = strings.Replace(fn, "( ", "(ctx ", 1)

		l := strings.Count(fn, ",") - 2
		if l == 0 {
			continue
		}
		needKey := true
		for j := 0; j < strings.Count(fn, ","); j++ {
			if i == 0 {
				fn = strings.Replace(fn, ", string", ",key string", 1)
				if !strings.Contains(fn, "key string") {
					continue out
				}
				needKey = true
				continue
			}

			fn = strings.Replace(fn, ", ", ",arg"+strconv.Itoa(j+1)+" ", 1)
		}
		b.WriteString(`func (nc *NamedClient)`)
		b.WriteString(m.Name)
		b.WriteString(fn)
		b.WriteString("{\n")
		if needKey {
			b.WriteString("	arg1 = nc.Name + arg1\n\n")
		} else {
			b.WriteString("\n")
		}
		b.WriteString(`	return nc.UniversalClient.`)
		b.WriteString(m.Name)
		b.WriteString("(")
		b.WriteString("ctx, ")

		for j, l := 0, strings.Count(fn, ","); j < strings.Count(fn, ","); j++ {
			b.WriteString("arg")
			b.WriteString(strconv.Itoa(j + 1))
			if j == l-1 {
				continue
			}
			b.WriteString(", ")
		}
		b.WriteString(")\n}")

		b.WriteString("\n")
		fmt.Println(b.String())

	}
	os.WriteFile("code.txt", []byte(b.String()), 0666)

}
