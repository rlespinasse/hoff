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

type RunState struct {
	value  State
	branch *string
	err    error
}

func RunStatePass() RunState {
	return RunState{
		value: pass,
	}
}

func RunStateBranchPass(branch string) RunState {
	return RunState{
		value:  pass,
		branch: &branch,
	}
}

func RunStateStop() RunState {
	return RunState{
		value: stop,
	}
}

func RunStateBranchStop(branch string) RunState {
	return RunState{
		value:  stop,
		branch: &branch,
	}
}

func RunStateFail(err error) RunState {
	return RunState{
		value: fail,
		err:   err,
	}
}
