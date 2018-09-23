package namingishard

import (
	"errors"
)

type DecisionNode struct {
	name         string
	decisionFunc func(*Context) (bool, error)
}

func (n DecisionNode) String() string {
	return n.name
}

func (n *DecisionNode) Compute(c *Context) ComputeState {
	decision, err := n.decisionFunc(c)
	if err != nil {
		return ComputeStateAbort(err)
	}
	if decision {
		return ComputeStateContinueOnBranch(true)
	}
	return ComputeStateContinueOnBranch(false)
}

func (n *DecisionNode) decideCapability() bool {
	return true
}

func NewDecisionNode(name string, decisionFunc func(*Context) (bool, error)) (*DecisionNode, error) {
	if decisionFunc == nil {
		return nil, errors.New("can't create decision node without function")
	}
	return &DecisionNode{name: name, decisionFunc: decisionFunc}, nil
}
