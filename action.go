package gisk

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type ActionInterface interface {
	Parse(gisk *Gisk) error
}

type ActionType struct {
	ActionType string `json:"action_type" yaml:"action_type"`
}

type Assignment struct {
	ActionType ActionType
	Variate    string `json:"variate" yaml:"variate"`
	Value      string `json:"value" yaml:"value"`
}

func (a *Assignment) Parse(gisk *Gisk) error {

	key := strings.ReplaceAll(a.Variate, " ", "") //删除空格字符
	keys := strings.Split(key, "_")
	if len(keys) < 3 {
		return errors.New("variate trait error")
	}
	//判断变量类型
	if keys[0] != "variate" {
		return errors.New("variate trait error")
	}

	k := strings.Join(keys[1:len(keys)-1], "_")
	version := keys[len(keys)-1]
	//获取变量dsl
	dsl, err := gisk.DslGetter.GetDsl(VARIATE, k, version)
	if err != nil {
		return err
	}
	var variate Variate
	err = gisk.Unmarshal([]byte(dsl), &variate)
	if err != nil {
		return err
	}

	v, err := GetValueByTrait(gisk, a.Value)
	if err != nil {
		return err
	}
	if variate.ValueType != v.ValueType {
		return fmt.Errorf("不能为变量`%s`类型`%s`赋值类型为`%s`的值", variate.Key, variate.ValueType, v.ValueType)
	}
	gisk.Context.SetVariate(k, v)
	return nil
}

var actionMap sync.Map

func RegisterAction(name string, actionStruct ActionInterface) {
	actionMap.Store(name, actionStruct)
}

func getAction(name string) (ActionInterface, bool) {
	action, ok := actionMap.Load(name)
	if !ok {
		return nil, false
	}
	return action.(ActionInterface), true
}

func init() {
	RegisterAction("assignment", &Assignment{})
}
