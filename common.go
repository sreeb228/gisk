package gisk

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"regexp"
	"strings"
)

// ConvertValueType 转换值数据类型
func ConvertValueType(value interface{}, valueType ValueType) (Value interface{}, err error) {
	switch valueType {
	case NUMBER:
		Value = gconv.Float64(value)
		return
	case STRING:
		Value = gconv.String(value)
		return
	case ARRAY:
		str := gconv.String(value)
		Value = strings.Split(str, ",")
		return
	case MAP:
		var v map[string]interface{}
		str := gconv.String(value)
		err = json.Unmarshal([]byte(str), &v)
		if err != nil {
			return nil, err
		}
		Value = v
		return
	case BOOL:
		Value = gconv.Bool(value)
		return
	default:
		return nil, errors.New("value type not support")
	}
}

// GetValueByTrait 获取特征数据
func GetValueByTrait(gisk *Gisk, key string) (Value, error) {
	key = strings.ReplaceAll(string(key), " ", "") //删除空格字符
	keys := strings.Split(string(key), "_")

	switch keys[0] {
	case string(INPUT): //特征格式 input_1_string
		if len(keys) < 3 {
			return Value{}, errors.New("input trait error")
		}
		valueType := keys[len(keys)-1]
		value := strings.Join(keys[1:len(keys)-1], "_")
		s := &Input{
			ValueType: ValueType(valueType),
			Value:     value,
		}
		return s.Parse(gisk)
	case string(VARIATE): //特征格式 variate_age_1
		if len(keys) < 3 {
			return Value{}, errors.New("variate trait error")
		}
		k := strings.Join(keys[1:len(keys)-1], "_")
		version := keys[len(keys)-1]
		dsl, e := gisk.DslGetter.GetDsl(VARIATE, k, version)
		if e != nil {
			return Value{}, e
		}
		var variate Variate
		var err error
		switch gisk.DslFormat {
		case JSON:
			err = json.Unmarshal([]byte(dsl), &variate)
		case YAML:
			err = yaml.Unmarshal([]byte(dsl), &variate)
		default:
			return Value{}, errors.New("failed to unmarshal DSL: " + err.Error())
		}
		return variate.Parse(gisk)
	case string(FUNC): //特征格式 func_rand(func_sum(input_1_number,input_2_number),input_100_number)函数可以嵌套函数
		re := regexp.MustCompile(`^func_(.+?)\((.+)?\)$`)
		matches := re.FindStringSubmatch(key)
		if len(matches) < 3 {
			return Value{}, errors.New("func trait error")
		}
		funcName := matches[1]
		f, ok := getFunc(funcName)
		if !ok {
			return Value{}, fmt.Errorf("func %s not fond", funcName)
		}

		str := matches[2]
		if str != "" {
			var value []Value
			parameters := parseFuncParameter(str)
			for _, v := range parameters {
				d, err := GetValueByTrait(gisk, v)
				if err != nil {
					return Value{}, err
				}
				value = append(value, d)
			}
			return f(value...)
		} else {
			return f()
		}
	}
	return Value{}, errors.Wrap(errors.New("value trait not support"), "")
}

func parseFuncParameter(str string) []string {
	var res []string
	//左右括号数量
	var lNum int
	var rNum int
	var temStr string
	for _, char := range str {
		s := string(char)
		if lNum == 0 {
			switch s {
			case "(":
				lNum++
				temStr += s
				break
			case ",":
				if temStr != "" {
					res = append(res, temStr)
					temStr = ""
				}
				break
			default:
				temStr += s
			}
		} else if lNum > 0 {
			temStr += s
			switch s {
			case "(":
				lNum++
				break
			case ")":
				rNum++
				if lNum == rNum {
					res = append(res, temStr)
					temStr = ""
					lNum = 0
					rNum = 0
				}
				break
			}
		}
	}
	if temStr != "" {
		res = append(res, temStr)
	}
	return res
}
