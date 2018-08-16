package flowengine

type ActionNode struct {
	actionFunc func(*FlowContext)
}

func (n *ActionNode) Run(c *FlowContext) {
	n.actionFunc(c)
}

func NewActionNode(actionFunc func(*FlowContext)) *ActionNode {
	return &ActionNode{
		actionFunc: actionFunc,
	}
}
