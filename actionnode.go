package flow

type ActionNode struct {
	actionFunc func(*Context) (bool, error)
}

func (n *ActionNode) Compute(c *Context) ComputeState {
	pass, err := n.actionFunc(c)
	if err != nil {
		return ComputeStateFail(err)
	}
	if pass {
		return ComputeStatePass()
	}
	return ComputeStateStop()
}

func (n *ActionNode) AvailableBranches() []string {
	return nil
}

func NewActionNode(actionFunc func(*Context) (bool, error)) *ActionNode {
	return &ActionNode{
		actionFunc: actionFunc,
	}
}
