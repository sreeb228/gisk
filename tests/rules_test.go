package tests

import (
	"encoding/json"
	"fmt"
	"gitee.com/sreeb/gisk"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

type getUrl struct {
	gisk.ActionType
	Url string `json:"url"`
}

func (u *getUrl) Parse(gisk *gisk.Gisk) error {
	// 发送GET请求
	resp, err := http.Get(u.Url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 打印响应体内容
	bodyStr := string(bodyBytes)
	fmt.Println("Response Body:", bodyStr)
	return nil
}

func TestRules_Parse(t *testing.T) {
	// 注册请求url动作
	gisk.RegisterAction("geturl", &getUrl{})

	g := gisk.New()
	g.SetDslGetter(&fileDslGetter{})

	bytes, _ := os.ReadFile("./dsl/rules.json")
	var rules gisk.Rules
	json.Unmarshal(bytes, &rules)
	_, err := rules.Parse(g)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(g.GetVariates())

}
