package flowengine

type DecisionNode struct {
	decisionFunc func(*FlowContext) bool
}

func (n *DecisionNode) Run(c *FlowContext) {
	n.decisionFunc(c)
}

func NewDecisionNode(decisionFunc func(*FlowContext) bool) *DecisionNode {
	return &DecisionNode{
		decisionFunc: decisionFunc,
	}
}
