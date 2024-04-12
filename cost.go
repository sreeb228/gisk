package gisk

// ValueType 值类型
type ValueType string

const (
	NUMBER ValueType = "number"
	STRING ValueType = "string"
	ARRAY  ValueType = "array"
	MAP    ValueType = "map"
	BOOL   ValueType = "bool"
)

// ElementType 元素类型
type ElementType string

const (
	VARIATE ElementType = "variate" //变量
	INPUT   ElementType = "input"   //输入值
	FUNC    ElementType = "func"    //函数
	RULE    ElementType = "rule"    //规则
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

// BreakMode 规则集规则中断模式
type BreakMode string

const (
	HitBreak  = "hit_break"  //规则命中中断后续规则
	MissBreak = "miss_break" //规则未命中中断后续规则
)
