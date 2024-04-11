package tests

import (
	"gitee.com/sreeb/gisk"
	"github.com/gogf/gf/v2/util/gconv"
	"testing"
)

func TestFunc(t *testing.T) {
	//注册求和函数
	gisk.RegisterFunc("sum", func(parameters ...gisk.Value) (gisk.Value, error) {
		var v float64
		for _, parameter := range parameters {
			v += gconv.Float64(parameter.Value)
		}
		return gisk.Value{ValueType: gisk.NUMBER, Value: v}, nil
	})

	//解析函数表达式
	g := gisk.New()
	g.SetDslGetter(&fileDslGetter{})

	res, err := gisk.GetValueByTrait(g, "func_rand(input_1_number,func_sum(input_10_number,input_20_string))")
	if err != nil {
		t.Errorf("err: %v", err)
	}
	t.Logf("res: %v", res)
}
