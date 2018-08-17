package flow

import (
	"github.com/google/go-cmp/cmp"
)

var emptySystem = NewNodeSystem()

var someActionNode = NewActionNode(func(*Context) (bool, error) { return true, nil })
var anotherActionNode = NewActionNode(func(*Context) (bool, error) { return true, nil })
var alwaysTrueDecisionNode = NewDecisionNode(func(*Context) (bool, error) { return true, nil })

var errorEqualOpts = cmp.Comparer(func(x, y error) bool {
	if x == nil || y == nil {
		return x == nil && y == nil
	}
	return x.Error() == y.Error()
})
