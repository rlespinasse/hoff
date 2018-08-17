package flow

import "github.com/google/go-cmp/cmp"

type Node interface {
	Run(c *Context) RunState
	AvailableBranches() []string
}

var nodeEqualOpts = cmp.Comparer(func(x, y Node) bool {
	return x == y
})
