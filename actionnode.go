package flow

import (
	"errors"
)

type ActionNode struct {
	actionFunc func(*Context) (bool, error)
}

func (n *ActionNode) Compute(c *Context) ComputeState {
	pass, err := n.actionFunc(c)
	if err != nil {
		return ComputeStateStopOnError(err)
	}
	if pass {
		return ComputeStatePass()
	}
	return ComputeStateStop()
}

func (n *ActionNode) AvailableBranches() []string {
	return nil
}

func NewActionNode(actionFunc func(*Context) (bool, error)) (*ActionNode, error) {
	if actionFunc == nil {
		return nil, errors.New("can't create action node without function")
	}
	return &ActionNode{actionFunc: actionFunc}, nil
}
