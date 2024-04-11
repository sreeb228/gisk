# DSL定义详解
## 输入值
输入值是由风控人员配置规则时手动输入的值。

结构体：
```golang
type Input struct {
    valueType valueType   //值类型
    value     interface{} //值
}
```
在DSL元素中表示： 格式：`input_值_值类型`,示例：`input_1_string`,解释：输入值为1，类型为string

## 变量

```go
type Variate struct {
    Key       string      `json:"key" yaml:"key"`               //唯一标识
    Name      string      `json:"name" yaml:"name"`             //变量名称
    Desc      string      `json:"desc" yaml:"desc"`             //描述
    Version   string      `json:"version" yaml:"version"`       //版本
    ValueType ValueType   `json:"value_type" yaml:"value_type"` //值类型
    Default   interface{} `json:"default" yaml:"default"`       //默认值
    IsInput   bool        `json:"is_input" yaml:"is_input"`     //是否要从输入值中匹配
}
```
在DSL元素中表示： 格式：variate_变量key_版本,示例：variate_age_1,解释：变量key为age，版本为1。

从DSL获取到variate_age_1时，会优先从当前DSL中获取变量dsl，如果没有则从dslgetter中获取.

json:
```json
{
  "key": "age",
  "name": "年龄",
  "desc": "年龄",
  "version": "1",
  "value_type": "number",
  "default": 0,
  "is_input": true
}
```
yaml:
```yaml
key: age
name: 年龄
desc: 年龄
version: "1"
value_type: number
default: 0
is_input: true
```