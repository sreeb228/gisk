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
- 评分卡
- 条件分流
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
## 未完待续。。。