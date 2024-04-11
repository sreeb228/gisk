package gisk

import (
	"errors"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
	"sync"
)

// Operation 运算
func Operation(left Value, operator Operator, right Value) (result bool, err error) {

	//优先使用用户注册比较符
	fun, ok := getOperation(operator)
	if ok {
		return fun(left, operator, right)
	}

	switch operator {
	case EQ:
		result = gconv.String(left.Value) == gconv.String(right.Value)
		return
	case NEQ:
		result = gconv.String(left.Value) != gconv.String(right.Value)
		return
	case GT:
		result = gconv.Float64(left.Value) > gconv.Float64(right.Value)
		return
	case LT:
		result = gconv.Float64(left.Value) < gconv.Float64(right.Value)
		return
	case GTE:
		result = gconv.Float64(left.Value) >= gconv.Float64(right.Value)
		return
	case LTE:
		result = gconv.Float64(left.Value) <= gconv.Float64(right.Value)
		return
	case IN:
		switch right.ValueType {
		case MAP:
			for _, v := range right.Value.(map[string]interface{}) {
				if gconv.String(left.Value) == gconv.String(v) {
					result = true
					return
				}
			}
		case ARRAY:
			for _, v := range right.Value.([]string) {
				if gconv.String(left.Value) == gconv.String(v) {
					result = true
					return
				}
			}
		default:
			err = errors.New("IN运算符的右侧必须是数组或map")
		}
		return
	case NOTIN:
		switch right.ValueType {
		case MAP:
			for _, v := range right.Value.(map[interface{}]interface{}) {
				if gconv.String(left.Value) == gconv.String(v) {
					result = false
					return
				}
			}
		case ARRAY:
			for _, v := range right.Value.([]string) {
				if gconv.String(left.Value) == gconv.String(v) {
					result = false
					return
				}
			}
		default:
			err = errors.New("NOTIN运算符的右侧必须是数组或map")
		}
		return
	case LIKE:
		l := gconv.String(left.Value)
		r := gconv.String(right.Value)
		result = strings.Contains(r, l)
		return
	case NOTLIKE:
		l := gconv.String(left.Value)
		r := gconv.String(right.Value)
		result = !strings.Contains(r, l)
		return
	default:
		err = errors.New("operator not support")
	}
	return
}

type OperationFunc func(left Value, operator Operator, right Value) (result bool, err error)

var operationMap sync.Map

func RegisterOperation(name Operator, op OperationFunc) {
	operationMap.Store(name, op)
}

func getOperation(name Operator) (OperationFunc, bool) {
	op, ok := operationMap.Load(name)
	if !ok {
		return nil, false
	}
	return op.(OperationFunc), true
}
