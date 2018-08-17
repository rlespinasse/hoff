package flow

const (
	pass State = iota + 1
	stop
	fail
)

type State int

func (state State) String() string {
	names := [...]string{
		"pass",
		"stop",
		"fail",
	}
	if state < pass || state > fail {
		return "unknown"
	}
	return names[state-1]
}

type ComputeState struct {
	value  State
	branch *string
	err    error
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
