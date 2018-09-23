package namingishard

import (
	"errors"
)

type DecisionNode struct {
	decisionFunc func(*Context) (bool, error)
}

func (n *DecisionNode) Compute(c *Context) ComputeState {
	decision, err := n.decisionFunc(c)
	if err != nil {
		return ComputeStateStopOnError(err)
	}
	if decision {
		return ComputeStateBranchPass(true)
	}
	return ComputeStateBranchPass(false)
}

func (n *DecisionNode) decideCapability() bool {
	return true
}

func NewDecisionNode(decisionFunc func(*Context) (bool, error)) (*DecisionNode, error) {
	if decisionFunc == nil {
		return nil, errors.New("can't create decision node without function")
	}
	return &DecisionNode{decisionFunc: decisionFunc}, nil
}
