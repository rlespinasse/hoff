/*
Package computestate expose utility functions to create a ComputeState object.

NOTE: possible state values are availables in the "statetype" package
*/
package computestate

import (
	"fmt"

	"github.com/rlespinasse/hoff/internal/utils"
	"github.com/rlespinasse/hoff/statetype"
)

// ComputeState hold the result of a Node computation
type ComputeState struct {
	Value  statetype.StateType
	Branch *bool
	Error  error
}

// String print human-readable version of a compute state
func (cs ComputeState) String() string {
	branch := ""
	if cs.Branch != nil {
		branch = fmt.Sprintf(" on %v", *cs.Branch)
	}
	err := ""
	if cs.Error != nil {
		err = fmt.Sprintf(" on %v", cs.Error)
	}
	return fmt.Sprintf("'%v%v%v'", cs.Value, branch, err)
}

// Continue generate a computation state to continue to following nodes
func Continue() ComputeState {
	return ComputeState{
		Value: statetype.ContinueState,
	}
}

// ContinueOnBranch generate a computation state to continue to following nodes
// on a branch taken by an Decision Node (DecideCapability at true)
func ContinueOnBranch(branch bool) ComputeState {
	return ComputeState{
		Value:  statetype.ContinueState,
		Branch: utils.BoolPointer(branch),
	}
}

// Skip generate a computation state to specify
// that the Node computation have been skipped
func Skip() ComputeState {
	return ComputeState{
		Value: statetype.SkipState,
	}
}

// Abort generate a computation state to throw an unexpected error
func Abort(err error) ComputeState {
	return ComputeState{
		Value: statetype.AbortState,
		Error: err,
	}
}
