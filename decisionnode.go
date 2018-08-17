package flow

import "reflect"

type DecisionNode struct {
	decisionFunc func(*Context) (bool, error)
}

func (n *DecisionNode) Compute(c *Context) ComputeState {
	decision, err := n.decisionFunc(c)
	if err != nil {
		return ComputeStateFail(err)
	}
	if decision {
		return ComputeStateBranchPass("true")
	}
	return ComputeStateBranchPass("false")
}

func (n *DecisionNode) AvailableBranches() []string {
	return []string{"true", "false"}
}

func NewDecisionNode(decisionFunc func(*Context) (bool, error)) *DecisionNode {
	return &DecisionNode{
		decisionFunc: decisionFunc,
	}
}

func isDecisionNode(n Node) bool {
	return "*flow.DecisionNode" == reflect.TypeOf(n).String()
}
