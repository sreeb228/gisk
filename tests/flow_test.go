package tests

import (
	"fmt"
	"gitee.com/sreeb/gisk"
	"testing"
)

func TestFlow_Parse(t *testing.T) {
	g := gisk.New()
	g.SetDslGetter(&fileDslGetter{})

	err := g.Parse(gisk.FLOW, "flow", "1")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(g.GetVariates())
}
