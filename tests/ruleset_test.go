package tests

import (
	"encoding/json"
	"fmt"
	"gitee.com/sreeb/gisk"
	"os"
	"testing"
)

func TestRuleset_Parse(t *testing.T) {

	g := gisk.New()
	g.SetDslGetter(&fileDslGetter{})

	bytes, _ := os.ReadFile("./dsl/ruleset.json")
	var ruleset gisk.Ruleset
	json.Unmarshal(bytes, &ruleset)
	err := ruleset.Parse(g)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(g.GetVariates())

}
