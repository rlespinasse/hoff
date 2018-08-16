package flow

type ActionNode struct {
	actionFunc func(*Context) (bool, error)
}

func (n *ActionNode) Run(c *Context) RunState {
	pass, err := n.actionFunc(c)
	if err != nil {
		return RunStateFail(err)
	}
	if pass {
		return RunStatePass()
	}
	return RunStateStop()
}

func (n *ActionNode) AvailableBranches() []NodeBranch {
	return AvailablesBranches()
}

func NewActionNode(actionFunc func(*Context) (bool, error)) *ActionNode {
	return &ActionNode{
		actionFunc: actionFunc,
	}
}
