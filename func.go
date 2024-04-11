package gisk

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"math/rand"
	"sync"
	"time"
)

type Func func(parameters ...Value) (Value, error)

var funcMap sync.Map

func RegisterFunc(name string, function Func) {
	funcMap.Store(name, function)
}

func getFunc(name string) (Func, bool) {
	function, ok := funcMap.Load(name)
	if !ok {
		return nil, false
	}
	return function.(Func), true
}

func init() {
	//注册随机函数
	RegisterFunc("rand", func(parameters ...Value) (Value, error) {
		if len(parameters) < 2 {
			return Value{}, errors.New("rand function need two parameters")
		}
		n1 := gconv.Int(parameters[0].Value)
		n2 := gconv.Int(parameters[1].Value)
		if n1 >= n2 || n1 < 0 {
			return Value{}, fmt.Errorf("rand函数参数错误，参数2应该大于参数1[%d,%d]", n1, n2)
		}
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		res := rng.Intn(n2-n1+1) + n1
		return Value{ValueType: NUMBER, Value: res}, nil
	})
	//注册求和函数
	RegisterFunc("sum", func(parameters ...Value) (Value, error) {
		var v float64
		for _, parameter := range parameters {
			v += gconv.Float64(parameter.Value)
		}
		return Value{ValueType: NUMBER, Value: v}, nil
	})
}
