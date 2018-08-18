package flow

import (
	"github.com/google/go-cmp/cmp"
)

var pointerOfComputationEqualOpts = cmp.Comparer(func(x, y *computation) bool {
	return (x == nil && y == nil) || (x != nil && y != nil && cmp.Equal(*x, *y, computationEqualOpts))
})

var computationEqualOpts = cmp.Comparer(func(x, y computation) bool {
	return x.computation == y.computation && cmp.Equal(x.context, y.context, pointerOfContextEqualOpts) && x.system == y.system && cmp.Equal(x.report, y.report)
})
