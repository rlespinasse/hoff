package hoff

import (
	"errors"
)

// DecisionNode is a type of Node who compute a function
// to take a decision based on Context.
type DecisionNode struct {
	name         string
	decisionFunc func(*Context) (bool, error)
}

func (n DecisionNode) String() string {
	return n.name
}

// Compute run the decision function and decide which compute state to return.
func (n *DecisionNode) Compute(c *Context) ComputeState {
	decision, err := n.decisionFunc(c)
	if err != nil {
		return NewAbortComputeState(err)
	}
	if decision {
		return NewContinueOnBranchComputeState(true)
	}
	return NewContinueOnBranchComputeState(false)
}

// DecideCapability is actived due to the fact that an decision is taked during compute.
func (n *DecisionNode) DecideCapability() bool {
	return true
}

// NewDecisionNode create a DecisionNode based on a name and a function to take the needed decision.
func NewDecisionNode(name string, decisionFunc func(*Context) (bool, error)) (*DecisionNode, error) {
	if decisionFunc == nil {
		return nil, errors.New("can't create decision node without function")
	}
	return &DecisionNode{name: name, decisionFunc: decisionFunc}, nil
}
