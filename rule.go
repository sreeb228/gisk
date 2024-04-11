package gisk

type Rule struct {
	Left     string   `json:"left" yaml:"left"`
	Operator Operator `json:"operator" yaml:"operator"` // 比较符号
	Right    string   `json:"right" yaml:"right"`
}

func (rule *Rule) Parse(gisk *Gisk) (bool bool, err error) {
	left, err := GetValueByTrait(gisk, rule.Left)
	if err != nil {
		return
	}
	right, err := GetValueByTrait(gisk, rule.Right)
	if err != nil {
		return
	}
	return Operation(left, rule.Operator, right)
}
