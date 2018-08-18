package flow

import (
	"github.com/google/go-cmp/cmp"
)

const (
	pass State = "pass"
	stop       = "stop"
	fail       = "fail"
)

type State string

type ComputeState struct {
	value  State
	branch *string
	err    error
}

func (x ComputeState) Equal(y ComputeState) bool {
	return x.value == y.value && cmp.Equal(x.branch, y.branch) && cmp.Equal(x.err, y.err, errorEqualOpts)
}

func ComputeStatePass() ComputeState {
	return ComputeState{
		value: pass,
	}
}

func ComputeStateBranchPass(branch string) ComputeState {
	return ComputeState{
		value:  pass,
		branch: &branch,
	}
}

func ComputeStateStop() ComputeState {
	return ComputeState{
		value: stop,
	}
}

func ComputeStateBranchStop(branch string) ComputeState {
	return ComputeState{
		value:  stop,
		branch: &branch,
	}
}

func ComputeStateFail(err error) ComputeState {
	return ComputeState{
		value: fail,
		err:   err,
	}
}
