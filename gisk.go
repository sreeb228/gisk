package gisk

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"sync"
)

const (
	JSON = "json"
	YAML = "yaml"
)

type Gisk struct {
	DslFormat string                 //dsl解析格式：json 或 yaml
	Input     map[string]interface{} //输入值
	Context   *Context               //上下文
	DslGetter *DslGetter             //dsl获取器
}

func New() *Gisk {
	return &Gisk{
		DslFormat: JSON,
		Input:     make(map[string]interface{}),
		Context: &Context{
			variateData:  make(map[string]Value),
			variateMutex: sync.RWMutex{},
		},
		DslGetter: &DslGetter{
			Dsl: sync.Map{},
		},
	}
}

// SetDslFormat 设置dsl格式
func (gisk *Gisk) SetDslFormat(format string) *Gisk {
	if format != JSON && format != YAML {
		panic("dsl解析格式仅支持json或yaml")
	}
	gisk.DslFormat = format
	return gisk
}

// SetDslGetter 设置dsl获取器
func (gisk *Gisk) SetDslGetter(getter DslGetterInterface) *Gisk {
	gisk.DslGetter.Getter = getter
	return gisk
}

func (gisk *Gisk) Parse(elementType ElementType, key string, version string) error {
	switch elementType {
	case RULE:
		rule, err := GetRule(gisk, key, version)
		if err != nil {
			return err
		}
		_, err = rule.Parse(gisk)
		return err
	case RULESET:
		ruleset, err := GetRuleset(gisk, key, version)
		if err != nil {
			return err
		}
		err = ruleset.Parse(gisk)
		return err
	}

	return nil
}

// GetVariates 获取所有赋值变量
func (gisk *Gisk) GetVariates() map[string]interface{} {
	res := make(map[string]interface{})
	variates := gisk.Context.GetVariates()
	for k, v := range variates {
		res[k] = v.Value
	}
	return res
}

func (gisk *Gisk) Unmarshal(data []byte, v any) error {
	var err error
	if gisk.DslFormat == JSON {
		err = json.Unmarshal(data, &v)
	} else if gisk.DslFormat == YAML {
		err = yaml.Unmarshal(data, &v)
	}
	return err
}
