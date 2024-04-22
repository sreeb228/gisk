package gisk

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"sync"
)

type Ruleset struct {
	Key      string         `json:"key" yaml:"key"`           //唯一标识
	Name     string         `json:"name" yaml:"name"`         //名称
	Desc     string         `json:"desc" yaml:"desc"`         //描述
	Version  string         `yaml:"version" json:"version"`   //版本
	Parallel bool           `json:"parallel" yaml:"parallel"` //是否并发执行
	Rules    []*rulesetRule `json:"rules" yaml:"rules"`
}

type rulesetRule struct {
	RuleKey     string    `json:"rule_key" yaml:"rule_key"`
	RuleVersion string    `json:"rule_version" yaml:"rule_version"`
	BreakMode   BreakMode `json:"break_mode" yaml:"break_mode"`
}

func (r *Ruleset) Parse(gisk *Gisk) error {
	//并发执行，并行执行不在处理执行顺序和中断模式
	if r.Parallel {
		// 使用 WaitGroup 等待所有协程执行完毕
		var wg sync.WaitGroup
		wg.Add(len(r.Rules))

		errChan := make(chan error, len(r.Rules))
		defer close(errChan)
		for _, rule := range r.Rules {
			go func(rulesetRule *rulesetRule) {
				defer wg.Done()
				r, err := GetRule(gisk, rulesetRule.RuleKey, rulesetRule.RuleVersion)
				if err != nil {
					errChan <- err
					return
				}
				_, err = r.Parse(gisk)
				if err != nil {
					errChan <- err
				}
			}(rule)
		}
		wg.Wait()

		//这里如果使用for range 会导致报错，因为errChan已经不在其他协程中，永远不会写入进去数据了，
		//range到chan中没有数据就会panic，可以close后再循环，但是这样可能在close前发生panic导致内存泄漏
		if len(errChan) > 0 {
			return <-errChan
		}
	} else {
		for _, rule := range r.Rules {
			r, err := GetRule(gisk, rule.RuleKey, rule.RuleVersion)
			if err != nil {
				return err
			}
			res, err := r.Parse(gisk)
			if err != nil {
				return err
			}

			if res && rule.BreakMode == HitBreak {
				return nil
			} else if !res && rule.BreakMode == MissBreak {
				return nil
			}
		}
	}
	return nil
}

func GetRuleset(gisk *Gisk, key string, version string) (*Ruleset, error) {
	dsl, _ := gisk.DslGetter.GetDsl(RULESET, key, version)
	var r Ruleset
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
