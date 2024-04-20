package gisk

import (
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

		errChan := make(chan error)
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
		for err := range errChan {
			if err != nil {
				return err
			}
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
