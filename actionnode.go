package namingishard

import (
	"errors"
)

type ActionNode struct {
	name       string
	actionFunc func(*Context) (bool, error)
}

func (n ActionNode) String() string {
	return n.name
}

func (n *ActionNode) Compute(c *Context) ComputeState {
	state, err := n.actionFunc(c)
	if err != nil {
		return ComputeStateAbort(err)
	}
	if state {
		return ComputeStateContinue()
	}
	return ComputeStateStop()
}

func (n *ActionNode) decideCapability() bool {
	return false
}

func NewActionNode(name string, actionFunc func(*Context) (bool, error)) (*ActionNode, error) {
	if actionFunc == nil {
		return nil, errors.New("can't create action node without function")
	}
	return &ActionNode{name: name, actionFunc: actionFunc}, nil
}
