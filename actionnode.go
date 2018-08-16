package flowengine

type ActionNode struct {
	actionFunc func(*FlowContext) (bool, error)
}

func (n *ActionNode) Run(c *FlowContext) RunState {
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

func NewActionNode(actionFunc func(*FlowContext) (bool, error)) *ActionNode {
	return &ActionNode{
		actionFunc: actionFunc,
	}
}
