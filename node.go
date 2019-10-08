package hoff

import (
	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/computestate"
)

// Node define
type Node interface {
	// Compute a node based on a context
	Compute(c *Context) computestate.ComputeState
	// DecideCapability tell if the Node can decide.
	// This impact the compute state by adding a branch to the state.
	DecideCapability() bool
}

var (
	// NodeComparator is a google/go-cmp comparator of Node
	NodeComparator = cmp.Comparer(func(x, y Node) bool {
		return x == y
	})
)
