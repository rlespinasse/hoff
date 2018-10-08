package node

import (
	"errors"

	"github.com/rlespinasse/hoff/computestate"
)

// ActionNode is a type of Node who compute a function
// to realize some actions based on Context.
type ActionNode struct {
	name       string
	actionFunc func(*Context) error
}

func (n ActionNode) String() string {
	return n.name
}

// Compute run the action function and decide which compute state to return.
func (n *ActionNode) Compute(c *Context) computestate.ComputeState {
	err := n.actionFunc(c)
	if err != nil {
		return computestate.Abort(err)
	}
	return computestate.Continue()
}

// DecideCapability is desactived due to the fact that an action don't take a decision.
func (n *ActionNode) DecideCapability() bool {
	return false
}

// NewAction create a ActionNode based on a name and a function to realize the needed action.
func NewAction(name string, actionFunc func(*Context) error) (*ActionNode, error) {
	if actionFunc == nil {
		return nil, errors.New("can't create action node without function")
	}
	return &ActionNode{name: name, actionFunc: actionFunc}, nil
}
