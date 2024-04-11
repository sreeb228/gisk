package gisk

type ValueType string

const (
	NUMBER ValueType = "number"
	STRING ValueType = "string"
	ARRAY  ValueType = "array"
	MAP    ValueType = "map"
	BOOL   ValueType = "bool"
)

// 元素类型
type ElementType string

const (
	VARIATE ElementType = "variate" //变量
	INPUT   ElementType = "input"   //输入值
	FUNC    ElementType = "func"    //函数
	RULES   ElementType = "rules"   //规则
)

// Operator 运算符
type Operator string

const (
	EQ      Operator = "eq"
	NEQ     Operator = "neq"
	GT      Operator = "gt"
	LT      Operator = "lt"
	GTE     Operator = "gte"
	LTE     Operator = "lte"
	IN      Operator = "in"
	NOTIN   Operator = "notIn"
	LIKE    Operator = "like"
	NOTLIKE Operator = "notLike"
)
