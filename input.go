package gisk

type Input struct {
	ValueType ValueType
	Value     interface{}
}

func (input *Input) Parse(gisk *Gisk) (Value Value, err error) {
	value, err := ConvertValueType(input.Value, input.ValueType)
	Value.Value = value
	Value.ValueType = input.ValueType
	return
}
