package hoff

import "github.com/google/go-cmp/cmp"

type Node interface {
	Compute(c *Context) ComputeState
	decideCapability() bool
}

var equalOptionForNode = cmp.Comparer(func(x, y Node) bool {
	return x == y
})
