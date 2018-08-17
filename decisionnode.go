package flow

import "reflect"

type DecisionNode struct {
	decisionFunc func(*Context) (bool, error)
}

func (n *DecisionNode) Run(c *Context) RunState {
	decision, err := n.decisionFunc(c)
	if err != nil {
		return RunStateFail(err)
	}
	if decision {
		return RunStateBranchPass("true")
	}
	return RunStateBranchPass("false")
}

func (n *DecisionNode) AvailableBranches() []NodeBranch {
	return AvailablesBranches("true", "false")
}

func NewDecisionNode(decisionFunc func(*Context) (bool, error)) *DecisionNode {
	return &DecisionNode{
		decisionFunc: decisionFunc,
	}
}

func isDecisionNode(n Node) bool {
	return "*flow.DecisionNode" == reflect.TypeOf(n).String()
}
