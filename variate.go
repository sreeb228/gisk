package gisk

type Variate struct {
	Key       string      `json:"key" yaml:"key"`               //唯一标识
	Name      string      `json:"name" yaml:"name"`             //变量名称
	Desc      string      `json:"desc" yaml:"desc"`             //描述
	Version   string      `json:"version" yaml:"version"`       //版本
	ValueType ValueType   `json:"value_type" yaml:"value_type"` //值类型
	Default   interface{} `json:"default" yaml:"default"`       //默认值
	IsInput   bool        `json:"is_input" yaml:"is_input"`     //是否要从输入值中匹配
}

func (variate *Variate) Parse(gisk *Gisk) (Value Value, err error) {
	if Value, ok := gisk.Context.GetVariate(variate.Key); ok {
		return Value, nil
	}

	v := variate.Default
	if variate.IsInput {
		if input, ok := gisk.Input[variate.Key]; ok {
			v = input
		}
	}
	value, err := ConvertValueType(v, variate.ValueType)
	Value.ValueType = variate.ValueType
	Value.Value = value

	gisk.Context.SetVariate(variate.Key, Value)
	return
}
