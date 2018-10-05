/*
Package node define action, and decision nodes.

You can also create your own node if the struct respect the Node interface.
A Node use a Context during computation as input data and output data.

NOTE: you will also find a "go-cmp" comparator for Node.
*/
package node

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
