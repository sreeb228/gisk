package gisk

import (
	"fmt"
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
		f.nodes = make(map[string]NodeInterface)
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
			case BranchFlowNode:
				var node branchFlowNode
				if err := gisk.Unmarshal(nodeRaw, &node); err != nil {
					f.errInitNodes = err
					return
				}
				f.nodes[node.NodeKey] = &node
			case ActionFlowNode:
				var node actionFlowNode
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
	NodeKey    string       `json:"node_key" yaml:"node_key"`               //节点key
	NodeType   FlowNodeType `json:"node_type" yaml:"node_type"`             //节点类型
	EleType    ElementType  `json:"element_type" yaml:"element_type"`       //元素类型
	EleKey     string       `json:"element_key" yaml:"element_key"`         //元素key
	EleVersion string       `json:"element_version" yaml:"element_version"` //元素版本
	NextNode   string       `json:"next_node" yaml:"next_node"`             //下一个节点
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

// 分流节点
type branchFlowNode struct {
	NodeKey  string       `json:"node_key" yaml:"node_key"`   //节点key
	NodeType FlowNodeType `json:"node_type" yaml:"node_type"` //节点类型
	Left     string       `json:"left" yaml:"left"`           //左侧
	Branches []struct {
		Operator Operator `json:"operator" yaml:"operator"`   // 比较符号
		Right    string   `json:"right" yaml:"right"`         // 右侧
		NextNode string   `json:"next_node" yaml:"next_node"` // 下一个节点
	} `json:"branches" yaml:"branches"` // 分支
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

// 动作节点
type actionFlowNode struct {
	NodeKey  string       `json:"node_key" yaml:"node_key"`   //节点key
	NodeType FlowNodeType `json:"node_type" yaml:"node_type"` //节点类型
	Actions  []RawMessage `json:"actions" yaml:"actions"`     //动作
	NextNode string       `json:"next_node" yaml:"next_node"` //下一个节点
}

func (node *actionFlowNode) Parse(gisk *Gisk, flow *Flow) ([]NodeInterface, error) {
	for _, action := range node.Actions {
		//先获取动作类型
		var actionType ActionType
		err := gisk.Unmarshal(action, &actionType)
		if err != nil {
			return nil, err
		}
		actionStruct, ok := getAction(actionType.ActionType)
		if !ok {
			return nil, fmt.Errorf("action type %s not found", actionType.ActionType)
		}
		err = gisk.Unmarshal(action, &actionStruct)
		if err != nil {
			return nil, err
		}
		err = actionStruct.Parse(gisk)
		if err != nil {
			return nil, err
		}
	}
	nextNodes := make([]NodeInterface, 0)
	if node.NextNode != "" {
		nextNode, err := flow.GetNode(gisk, node.NextNode)
		if err != nil {
			return nil, err
		}
		if nextNode == nil {
			return nil, errors.New("next node not found")
		}
		nextNodes = append(nextNodes, nextNode)
	}
	return nextNodes, nil
}

func GetFlow(gisk *Gisk, key string, version string) (*Flow, error) {
	dsl, _ := gisk.DslGetter.GetDsl(FLOW, key, version)
	var f Flow
	err := gisk.Unmarshal([]byte(dsl), &f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}
