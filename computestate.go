package namingishard

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

const (
	Continue State = "Continue"
	Stop           = "Stop"
	Skip           = "Skip"
	Abort          = "Abort"
)

type State string

type ComputeState struct {
	value  State
	branch *bool
	err    error
}

func (c ComputeState) String() string {
	branch := ""
	if c.branch != nil {
		branch = fmt.Sprintf(" on %v", *c.branch)
	}
	err := ""
	if c.err != nil {
		err = fmt.Sprintf(" on %v", c.err)
	}
	return fmt.Sprintf("'%v%v%v'", c.value, branch, err)
}

func (x ComputeState) Equal(y ComputeState) bool {
	return cmp.Equal(x.value, y.value) && cmp.Equal(x.branch, y.branch) && cmp.Equal(x.err, y.err, equalOptionForError)
}

func ComputeStateContinue() ComputeState {
	return ComputeState{
		value: Continue,
	}
}

func ComputeStateContinueOnBranch(branch bool) ComputeState {
	return ComputeState{
		value:  Continue,
		branch: boolPointer(branch),
	}
}

func ComputeStateStop() ComputeState {
	return ComputeState{
		value: Stop,
	}
}

func ComputeStateSkip() ComputeState {
	return ComputeState{
		value: Skip,
	}
}

func ComputeStateAbort(err error) ComputeState {
	return ComputeState{
		value: Abort,
		err:   err,
	}
}
