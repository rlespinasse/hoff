package flow

import "github.com/google/go-cmp/cmp"

var emptyContext = NewContext()

var pointerOfContextEqualOpts = cmp.Comparer(func(x, y *Context) bool {
	return (x == nil && y == nil) || (x != nil && y != nil && cmp.Equal(*x, *y, contextEqualOpts))
})

var contextEqualOpts = cmp.Comparer(func(x, y Context) bool {
	return cmp.Equal(x.data, y.data)
})
