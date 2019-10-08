package hoff

import (
	"fmt"

	"github.com/rlespinasse/hoff/internal/utils"
)

// ComputeState hold the result of a Node computation
type ComputeState struct {
	Value  StateType
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

// NewContinueComputeState generate a computation state to continue to following nodes
func NewContinueComputeState() ComputeState {
	return ComputeState{
		Value: ContinueState,
	}
}

// NewContinueOnBranchComputeState generate a computation state to continue to following nodes
// on a branch taken by an Decision Node (DecideCapability at true)
func NewContinueOnBranchComputeState(branch bool) ComputeState {
	return ComputeState{
		Value:  ContinueState,
		Branch: utils.BoolPointer(branch),
	}
}

// NewSkipComputeState generate a computation state to specify
// that the Node computation have been skipped
func NewSkipComputeState() ComputeState {
	return ComputeState{
		Value: SkipState,
	}
}

// NewAbortComputeState generate a computation state to throw an unexpected error
func NewAbortComputeState(err error) ComputeState {
	return ComputeState{
		Value: AbortState,
		Error: err,
	}
}
