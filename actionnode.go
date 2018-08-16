package flowengine

type ActionNode struct {
	action func(*FlowContext)
}

func (n *ActionNode) Run(c *FlowContext) {
	n.action(c)
}

func NewActionNode(action func(*FlowContext)) *ActionNode {
	return &ActionNode{
		action: action,
	}
}
