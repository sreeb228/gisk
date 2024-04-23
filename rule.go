package gisk

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
	"strings"
	"sync"
)

type Compare struct {
	Left     string   `json:"left" yaml:"left"`
	Operator Operator `json:"operator" yaml:"operator"` // 比较符号
	Right    string   `json:"right" yaml:"right"`
}

func (c *Compare) Parse(gisk *Gisk) (bool bool, err error) {
	left, err := GetValueByTrait(gisk, c.Left)
	if err != nil {
		return
	}
	right, err := GetValueByTrait(gisk, c.Right)
	if err != nil {
		return
	}
	return Operation(left, c.Operator, right)
}

type Rule struct {
	Key         string              `json:"key" yaml:"key"`                //唯一标识
	Name        string              `json:"name" yaml:"name"`              //名称
	Desc        string              `json:"desc" yaml:"desc"`              //描述
	Version     string              `json:"version" yaml:"version"`        //版本
	Parallel    bool                `json:"parallel" yaml:"parallel"`      //是否并发执行
	Compares    map[string]*Compare `json:"compares" yaml:"compares"`      //比较
	Expression  string              `json:"expression"  yaml:"expression"` //计算公式
	ActionTure  []RawMessage        `json:"action_true" yaml:"action_true"`
	ActionFalse []RawMessage        `json:"action_false" yaml:"action_false"`
}

func (r *Rule) Parse(gisk *Gisk) (bool, error) {
	exp := splitExpression(r.Expression)
	RPN, err := InfixToRPN(exp)
	if err != nil {
		return false, err
	}
	resultMap := make(map[string]bool)
	if r.Parallel {
		resultMap, err = r.parallelParseRule(gisk)
		if err != nil {
			return false, err
		}
	}
	res, err := r.evalRPN(gisk, RPN, resultMap)
	if err != nil {
		return false, err
	}

	// 执行动作
	var actions []RawMessage
	if res {
		actions = r.ActionTure
	} else {
		actions = r.ActionFalse
	}

	for _, action := range actions {
		//先获取动作类型
		var actionType ActionType
		switch gisk.DslFormat {
		case JSON:
			err = json.Unmarshal(action, &actionType)
			if err != nil {
				return false, err
			}
			actionStruct, ok := getAction(actionType.ActionType)
			if !ok {
				return false, fmt.Errorf("action type %s not found", actionType.ActionType)
			}
			err = json.Unmarshal(action, &actionStruct)
			if err != nil {
				return false, err
			}
			err = actionStruct.Parse(gisk)
			if err != nil {
				return false, err
			}

		case YAML:
			err = yaml.Unmarshal(action, &actionType)
			if err != nil {
				return false, err
			}
			actionStruct, ok := getAction(actionType.ActionType)
			if !ok {
				return false, fmt.Errorf("action type %s not found", actionType.ActionType)
			}
			err = yaml.Unmarshal(action, &actionStruct)
			if err != nil {
				return false, err
			}
			err = actionStruct.Parse(gisk)
			if err != nil {
				return false, err
			}
		default:
			return false, errors.New("dsl format not support")
		}
	}
	return res, nil
}

// parallelParseRule 并行执行规则获取结果
func (r *Rule) parallelParseRule(gisk *Gisk) (map[string]bool, error) {
	resultMap := make(map[string]bool)
	var resMap sync.Map
	type resType struct {
		res bool
		err error
	}
	// 使用 WaitGroup 等待所有协程执行完毕
	var wg sync.WaitGroup
	wg.Add(len(r.Compares))

	for name, c := range r.Compares {
		go func(name string, c *Compare) {
			defer wg.Done()
			//执行规则，获取规则执行结果
			res, err := c.Parse(gisk)
			if err != nil {
				resMap.Store(name, resType{false, err})
			} else {
				resMap.Store(name, resType{res, nil})
			}
		}(name, c)
	}
	wg.Wait()
	var err error
	// 遍历结果
	resMap.Range(func(key, value interface{}) bool {
		index := key.(string)
		result := value.(resType)

		if result.err != nil {
			err = result.err
			return false
		}
		resultMap[index] = result.res
		return true
	})

	return resultMap, err
}

// splitExpression 将计算公式拆分成数组
func splitExpression(expression string) []string {
	expression = strings.ReplaceAll(expression, " ", "") //删除空格字符
	re := regexp.MustCompile(`[()]|&&|\|\||[^()&|]+`)
	parts := re.FindAllString(expression, -1)
	return parts
}

// InfixToRPN 将算公式转换为逆波兰表达式
func InfixToRPN(parts []string) ([]string, error) {
	var rpn []string                               // 初始化逆波兰表达式的输出结果切片
	var stack []string                             // 初始化操作符栈
	precedence := map[string]int{"||": 1, "&&": 2} // 定义逻辑运算符及其优先级

	for _, part := range parts {
		if op, ok := precedence[part]; ok {
			// 当栈非空且栈顶元素优先级大于等于当前运算符，或栈顶元素不为左括号时，将栈顶元素出栈并添加到rpn中
			for len(stack) > 0 && precedence[stack[len(stack)-1]] >= op && stack[len(stack)-1] != "(" {
				rpn = append(rpn, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, part) // 当前运算符入栈
		} else if part == "(" {
			stack = append(stack, part)
		} else if part == ")" {
			for stack[len(stack)-1] != "(" {
				if len(stack) == 0 {
					return nil, errors.New("没有匹配左括号的意外右括号")
				}
				rpn = append(rpn, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1] // 弹出左括号
		} else {
			rpn = append(rpn, part) // 直接将该部分添加到逆波兰表达式结果切片中
		}
	}

	// 清空栈并将剩余的所有元素依次添加到逆波兰表达式结果切片中
	for len(stack) > 0 {
		rpn = append(rpn, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return rpn, nil // 返回逆波兰表达式结果切片
}

// 计算逆波兰表达式
func (r *Rule) evalRPN(gisk *Gisk, RPN []string, resultMap map[string]bool) (bool, error) {
	stack := make([]bool, 0)
	for _, token := range RPN {
		if token == "&&" || token == "||" {
			operand2 := stack[len(stack)-1]
			operand1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2] // 出栈两次
			result := false
			switch token {
			case "&&":
				result = operand1 && operand2
			case "||":
				result = operand1 || operand2
			}
			stack = append(stack, result) // 将结果入栈
		} else {
			if v, ok := resultMap[token]; ok {
				stack = append(stack, v)
			} else {
				rule, ok := r.Compares[token]
				if !ok {
					panic("没有找到规则:" + token)
				}
				val, err := rule.Parse(gisk)
				if err != nil {
					return false, err
				}
				stack = append(stack, val)
			}
		}
	}
	return stack[0], nil
}

func GetRule(gisk *Gisk, key string, version string) (*Rule, error) {
	dsl, _ := gisk.DslGetter.GetDsl(RULE, key, version)
	var r Rule
	if gisk.DslFormat == JSON {
		err := json.Unmarshal([]byte(dsl), &r)
		if err != nil {
			return nil, err
		}
	} else if gisk.DslFormat == YAML {
		err := yaml.Unmarshal([]byte(dsl), &r)
		if err != nil {
			return nil, err
		}
	}
	return &r, nil
}
