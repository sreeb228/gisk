package gisk

import (
	"github.com/pkg/errors"
	"sync"
)

type Flow struct {
	Key     string       `json:"key" yaml:"key"`         //唯一标识
	Name    string       `json:"name" yaml:"name"`       //名称
	Desc    string       `json:"desc" yaml:"desc"`       //描述
	Version string       `yaml:"version" json:"version"` //版本
	Nodes   []RawMessage `json:"nodes" yaml:"nodes"`     //节点

	nodes         map[string]NodeInterface //保存解析过的节点
	initNodesOnce sync.Once
	errInitNodes  error
}

func (f *Flow) Parse(gisk *Gisk) error {
	startNode, err := f.getStartNode(gisk)
	if err != nil {
		return err
	}

	var nextNodes []NodeInterface
	queue := []NodeInterface{startNode}

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		nextNodes, err = currentNode.Parse(gisk, f)
		if err != nil {
			return err
		}
		queue = append(queue, nextNodes...)
	}
	return nil
}

type flowNodeType struct {
	NodeType FlowNodeType `json:"node_type" yaml:"node_type"`
}

// 获取开始节点
func (f *Flow) getStartNode(gisk *Gisk) (NodeInterface, error) {
	startNode, err := f.GetNode(gisk, FlowStartNodeKey)
	if err != nil {
		return nil, err
	}

	if startNode == nil {
		return nil, errors.New("start node not found")
	}
	return startNode, nil
}

func (f *Flow) GetNode(gisk *Gisk, key string) (NodeInterface, error) {
	f.initNodesOnce.Do(func() {
		for _, nodeRaw := range f.Nodes {
			var nodeType flowNodeType
			if err := gisk.Unmarshal(nodeRaw, &nodeType); err != nil {
				f.errInitNodes = err
				return
			}
			switch nodeType.NodeType {
			case GeneralFlowNode:
				var node generalFlowNode
				if err := gisk.Unmarshal(nodeRaw, &node); err != nil {
					f.errInitNodes = err
					return
				}
				f.nodes[node.NodeKey] = &node
			}
		}
	})

	if f.errInitNodes != nil {
		return nil, f.errInitNodes
	}
	if res, ok := f.nodes[key]; ok {
		return res, nil
	}
	return nil, nil
}

type NodeInterface interface {
	Parse(gisk *Gisk, flow *Flow) ([]NodeInterface, error)
}

// 普通节点
type generalFlowNode struct {
	NodeKey    string       `json:"node_key" yaml:"node_key"`
	NodeType   FlowNodeType `json:"node_type" yaml:"node_type"`
	EleType    ElementType  `json:"element_type" yaml:"element_type"`
	EleKey     string       `json:"element_key" yaml:"element_key"`
	EleVersion string       `json:"element_version" yaml:"element_version"`
	NextNode   string       `json:"next_node" yaml:"next_node"`
}

func (node *generalFlowNode) Parse(gisk *Gisk, flow *Flow) ([]NodeInterface, error) {
	switch node.EleType {
	case RULE:
		rule, err := GetRule(gisk, node.EleKey, node.EleVersion)
		if err != nil {
			return nil, err
		}
		_, err = rule.Parse(gisk)
		if err != nil {
			return nil, err
		}
	case RULESET:
		ruleset, err := GetRuleset(gisk, node.EleKey, node.EleVersion)
		if err != nil {
			return nil, err
		}
		err = ruleset.Parse(gisk)
		if err != nil {
			return nil, err
		}
	}
	if node.NextNode == "" {
		return nil, nil
	}

	nextNode, err := flow.GetNode(gisk, node.NextNode)
	if err != nil {
		return nil, err
	}
	if nextNode == nil {
		return nil, errors.New("next node not found")
	}
	return append([]NodeInterface{}, nextNode), nil
}

// 分支节点
type branchFlowNode struct {
	NodeKey  string       `json:"node_key" yaml:"node_key"`
	NodeType FlowNodeType `json:"node_type" yaml:"node_type"`
	Left     string       `json:"left" yaml:"left"`
	Branches []struct {
		Operator Operator `json:"operator" yaml:"operator"` // 比较符号
		Right    string   `json:"right" yaml:"right"`
		NextNode string   `json:"next_node" yaml:"next_node"`
	} `json:"branches" yaml:"branches"`
}

func (node *branchFlowNode) Parse(gisk *Gisk, flow *Flow) ([]NodeInterface, error) {
	left, err := GetValueByTrait(gisk, node.Left)
	if err != nil {
		return nil, err
	}
	nextNodes := make([]NodeInterface, 0)
	for _, branch := range node.Branches {
		right, err := GetValueByTrait(gisk, branch.Right)
		if err != nil {
			return nil, err
		}
		res, err := Operation(left, branch.Operator, right)
		if err != nil {
			return nil, err
		}
		if res {
			if branch.NextNode != "" {
				nextNode, err := flow.GetNode(gisk, branch.NextNode)
				if err != nil {
					return nil, err
				}
				if nextNode == nil {
					return nil, errors.New("next node not found")
				}
				nextNodes = append(nextNodes, nextNode)
			}
		}
	}
	return nextNodes, nil
}
