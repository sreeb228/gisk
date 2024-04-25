## gisk风控策略引擎

gisk是独立的即插即用的轻量级决策引擎，支持json和yaml格式DSL。支持自定义运算符，自定义函数，自定义动作开放灵活的扩展适用更多的业务场景。

### 功能列表
- 基础值
    - 变量
    - 输入值
    - 函数
- 决策（规则）
- 决策集
- 决策流
- 分流决策流
- 赋值决策流
- 评分卡
- 支持数据类型：number, string, bool, array, map
- 支持运算符号（支持自定义，覆盖默认运算符）：eq, neq, gt, lt, gte, lte, in, notIn, like, notLike
- 支持函数（支持自定义，函数支持嵌套）：内置函数rand、sum
- 支持决策动作（支持自定义）：内置动作 赋值,访问url
- 支持串行并行执行
- DSL支持历史版本控制（提供获取接口，实现不同介质DSL储存。内置文件存储）
- DSL支持json和yaml格式


## 快速开始
- 环境准备

  *go version go1.2+*


- 安装

```shell
go get -u -v gitee.com/sreeb/gisk
```
- 基础使用  

```go
package main

import (
  "fmt"
  "gitee.com/sreeb/gisk"
)

func main() {
  elementType := gisk.RULES //规则
  rulesKey := "rules1" //规则唯一key
  version := "1" //规则版本

  g := gisk.New() //创建gisk实例
  g.SetDslFormat(gisk.JSON) //设置dsl格式
  err := g.Parse(elementType, rulesKey, version) //解析规则
  if err != nil {
    //错误处理
  }
  //获取所有被初始化的变量值
  variates := g.GetVariates()
  fmt.Println(variates)
}

```

- 注册dsl获取接口

*dsl接口获取器,可以根据提供的类型和key和版本获取dsl字符串。 dsl接口获取器需要实现`gisk.DslGetterInterface`接口，返回dsl字符串。*
```go
type DslGetterInterface interface {
	GetDsl(elementType ElementType, key string, version string) (string, error)
}
```
示例：
```go
type fileDslGetter struct {
}

func (getter *fileDslGetter) GetDsl(elementType gisk.ElementType, key string, version string) (string, error) {
	path := "./dsl/" + elementType + "_" + key + "_" + version + ".json"
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return  string(bytes), nil
}

func main() {
    //创建gisk实例
    g := gisk.New()
    //设置dsl获取器
    g.SetDslGetter(&fileDslGetter{})
    //...
}
```


- 注册比较符

*系统实现默认的比较符 `eq`, `neq`, `gt`, `lt`, `gte`, `lte`,`in`,`notIn`,`like`,`notLike`如果需要自定义比较符，可以调用RegisterOperation方法。自定义比较符优先级高于系统内置的比较符，可以注册和系统同名比较符实现复写。*

```go
package main

import (
	"errors"
	"gitee.com/sreeb/gisk"
	"regexp"
)

func main() {
	//自定义正则匹配比较符
	gisk.RegisterOperation("reg", func(left gisk.Value, operator gisk.Operator, right gisk.Value) (result bool, err error) {
		if left.ValueType != gisk.STRING || right.ValueType != gisk.STRING {
			err = errors.New("left and right must be string")
			return
		}
		// 正则表达式匹配
		reg := regexp.MustCompile(right.Value.(string))
		if reg.MatchString(left.Value.(string)) {
			result = true
		}
		return
	})

	g := gisk.New()
	//...
}
```
- 注册函数

系统实现默认的函数 `rand`，`sum`，如果需要自定义函数，可以调用RegisterFunction方法。自定义函数优先级高于系统内置的函数，可以注册和系统同名函数实现复写。*
```go
package main

import (
  "gitee.com/sreeb/gisk"
  "github.com/gogf/gf/v2/util/gconv"
)

func main() {
    //注册求和函数
    gisk.RegisterFunc("sum", func(parameters ...gisk.Value) (gisk.Value, error) {
      var v float64
      for _, parameter := range parameters {
        v += gconv.Float64(parameter.Value)
      }
      return gisk.Value{ValueType: gisk.NUMBER, Value: v}, nil
    })
	
    g := gisk.New()
    //...
}
```
## DSL语法

### 基础值

*基础值包含变量，输入值，函数。在dsl中用字符串表示。基础值*

#### 变量

表示法：`variate_age_1`

解释：`variate`表示为变量类型，`age`表示变量唯一key（key可以有下划线），`1`表示变量版本号。

变量DSL：
```json
{
    "key": "age",
    "name": "年龄",
    "desc": "用户年龄",
    "version": "1",
    "value_type": "number",
    "default": 20,
    "is_input": false
}
```
变量结构体：
```go
type Variate struct {
	Key       string      `json:"key" yaml:"key"`               //唯一标识
	Name      string      `json:"name" yaml:"name"`             //变量名称（前端页面使用，不涉及逻辑）
	Desc      string      `json:"desc" yaml:"desc"`             //描述（前端页面使用，不涉及逻辑）
	Version   string      `json:"version" yaml:"version"`       //版本
	ValueType ValueType   `json:"value_type" yaml:"value_type"` //值类型（number：转成float64, string, bool, array:以,分割转成array, map：json字符串转成map）
	Default   interface{} `json:"default" yaml:"default"`       //默认值
	IsInput   bool        `json:"is_input" yaml:"is_input"`     //是否要从输入值中匹配(true时会从输入值中匹配相同key的值进行赋值)
}
```

#### 输入值
表示法：`input_1_number`

解释：`input`表示为输入值类型，`1`表示输入值，`number`表示输入值数据类型。

输入值没有DSL

#### 函数

表示法：`func_rand(input_1_number,func_sum(input_10_number,input_20_string))`

解释：函数支持嵌套，函数参数只能是基础值。上述表达式解释为：rand(1,sum(10,20))，rand函数两个参数1和sum函数，sum函数两个参数10和20。 注册函数需要实现 `type Func func(parameters ...Value) (Value, error)`类型。

函数无DSL

## 规则（决策）

*规则是基于多个比较条件通过逻辑运算符进行组合，最终得到一个布尔值。可以根据最终布尔值进行后续动作执行，赋值，访问url，发送消息，连接数据库等等操作。规则支持括号运算，支持串行并行执行*

规则dsl

```json
{
  "key": "rule1",   //唯一key
  "name": "用户筛选", //规则名称
  "desc": "用户筛选", //规则描述
  "version": "1",  //规则版本
  "parallel": true, //是否并行执行（并行执行会先并发获取比较条件组所有结果，再用规则表达式进行逻辑运算）
  "compares": {     
    "年龄小于40": {
      "left": "variate_年龄_1", //比较条件左边
      "operator": "lt", //比较条件运算符
      "right": "input_40_number" //比较条件右边
    },
    "性别女": {
      "left": "variate_性别_1",
      "operator": "eq",
      "right": "input_女_string"
    },
    "身高大于等于170": {
      "left": "variate_身高_1",
      "operator": "gte",
      "right": "input_170_number"
    }
  },  //比较条件组
  "expression": "年龄小于40 && (性别女 || 身高大于等于170)", //规则表达式
  "action_true": [
    {
      "action_type": "assignment",  //赋值操作
      "variate": "variate_命中结果_1", //赋值变量
      "value": "input_true_bool" //赋值值
    },
    {
      "action_type": "geturl", //访问url
      "url": "https://xxx.com" //url
    }
  ],  //true执行动作
  "action_false": [
    {
      "action_type": "assignment", //赋值操作
      "variate": "variate_命中结果_1", //赋值变量
      "value": "input_false_bool" //赋值值
    }
  ] //false执行动作
}

```
上述规则解释为：如果年龄小于40且性别为女或者身高大于等于170，则命中，否则未命中。命中时执行动作：变量命中结果赋值为true时，访问url。未命中时，赋值变量命中结果为false。


规则结构体：
```go
type Rule struct {
	Key         string              `json:"key" yaml:"key"`                //唯一标识
	Name        string              `json:"name" yaml:"name"`              //名称
	Desc        string              `json:"desc" yaml:"desc"`              //描述
	Version     string              `json:"version" yaml:"version"`        //版本
	Parallel    bool                `json:"parallel" yaml:"parallel"`      //是否并发执行
	Compares    map[string]*Compare `json:"compares" yaml:"compares"`      //比较
	Expression  string              `json:"expression"  yaml:"expression"` //计算公式
	ActionTure  []RawMessage        `json:"action_true" yaml:"action_true"` //命中执行动作
	ActionFalse []RawMessage        `json:"action_false" yaml:"action_false"` //未命中执行动作
}
type Compare struct {
  Left     string   `json:"left" yaml:"left"`
  Operator Operator `json:"operator" yaml:"operator"` // 比较符号(可自定义注册)
  Right    string   `json:"right" yaml:"right"`
}
```



## 规则集（决策集）

* 规则集是多个规则的集合，支持串行并行执行和中断。并行模式下规则的先后顺序和中断不生效*

规则集dsl:
```json
{
  "key": "ruleset", //唯一key
  "name": "决策集1", //规则集名称
  "desc": "决策集1", //规则集描述
  "version": "1", //规则集版本
  "parallel": false, //是否并行执行（串行执行：顺序执行规则，中断后不执行后续规则。 并行执行：并发执行规则，不考虑顺序和中断）
  "rules": [ 
    {
      "rule_key": "rule1", //规则key
      "rule_version": "1", //规则版本
      "break_mode": "hit_break" //中断模式:hit_break命中中断，miss_break未命中中断。中断表示不执行后续的规则
    },
    {
      "rule_key": "rule2",
      "rule_version": "1",
      "break_mode": "miss_break"
    },
    {
      "rule_key": "rule3",
      "rule_version": "1",
      "break_mode": "miss_break"
    }
  ]//规则集规则
}
```

规则集结构体：
```go
type Ruleset struct {
	Key      string         `json:"key" yaml:"key"`           //唯一标识
	Name     string         `json:"name" yaml:"name"`         //名称
	Desc     string         `json:"desc" yaml:"desc"`         //描述
	Version  string         `yaml:"version" json:"version"`   //版本
	Parallel bool           `json:"parallel" yaml:"parallel"` //是否并发执行
	Rules    []*rulesetRule `json:"rules" yaml:"rules"`       //规则集规则
}

type rulesetRule struct {
	RuleKey     string    `json:"rule_key" yaml:"rule_key"` //规则key
	RuleVersion string    `json:"rule_version" yaml:"rule_version"` //规则版本
	BreakMode   BreakMode `json:"break_mode" yaml:"break_mode"` //中断模式
}
```


## 决策流
决策流支持普通决策流，分流决策流和赋值决策流（规则树）

未完待续。。
