package flow

import "github.com/google/go-cmp/cmp"

type Node interface {
	Compute(c *Context) ComputeState
	AvailableBranches() []string
}

var equalOptionForNode = cmp.Comparer(func(x, y Node) bool {
	return x == y
})
