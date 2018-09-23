package namingishard

import (
	"errors"
)

type ActionNode struct {
	actionFunc func(*Context) error
}

func (n *ActionNode) Compute(c *Context) ComputeState {
	err := n.actionFunc(c)
	if err != nil {
		return ComputeStateStopOnError(err)
	}
	return ComputeStatePass()
}

func (n *ActionNode) decideCapability() bool {
	return false
}

func NewActionNode(actionFunc func(*Context) error) (*ActionNode, error) {
	if actionFunc == nil {
		return nil, errors.New("can't create action node without function")
	}
	return &ActionNode{actionFunc: actionFunc}, nil
}
