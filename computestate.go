package namingishard

import (
	"github.com/google/go-cmp/cmp"
)

const (
	pass State = "pass"
	stop       = "stop"
)

type State string

type ComputeState struct {
	value  State
	branch *bool
	err    error
}

func (x ComputeState) Equal(y ComputeState) bool {
	return cmp.Equal(x.value, y.value) && cmp.Equal(x.branch, y.branch) && cmp.Equal(x.err, y.err, equalOptionForError)
}

func ComputeStatePass() ComputeState {
	return ComputeState{
		value: pass,
	}
}

func ComputeStateBranchPass(branch bool) ComputeState {
	return ComputeState{
		value:  pass,
		branch: boolPointer(branch),
	}
}

func ComputeStateStopOnError(err error) ComputeState {
	return ComputeState{
		value: stop,
		err:   err,
	}
}
