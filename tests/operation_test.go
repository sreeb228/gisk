package tests

import (
	"gitee.com/sreeb/gisk"
	"regexp"
	"testing"
)

func TestOperation(t *testing.T) {
	// 注册自定义操作符
	gisk.RegisterOperation("reg", func(left gisk.Value, operator gisk.Operator, right gisk.Value) (result bool, err error) {
		if left.ValueType != gisk.STRING || right.ValueType != gisk.STRING {
			t.Errorf("left or right is not string")
			return
		}
		// 正则表达式匹配
		reg := regexp.MustCompile(right.Value.(string))
		if reg.MatchString(left.Value.(string)) {
			result = true
		}
		return
	})
	_ = gisk.New()
	res, err := gisk.Operation(gisk.Value{ValueType: gisk.STRING, Value: "A111"}, "reg", gisk.Value{ValueType: gisk.STRING, Value: "\\d+"})
	if err != nil {
		t.Errorf("err: %v", err)
	}
	t.Logf("res: %v", res)

}
