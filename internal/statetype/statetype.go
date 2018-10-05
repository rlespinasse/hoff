package statetype

// StateType is the type of the computation state of a Node during a Computation
type StateType string

const (
	// ContinueState tell that the Node computation authorize
	// to compute the following nodes
	ContinueState StateType = "Continue"
	// StopState tell that the Node computation stop
	// the computation of the following nodes
	StopState = "Stop"
	// SkipState tell that the Node computation is skipped
	SkipState = "Skip"
	// AbortState tell the Node computation encounter an error
	// and abort the computation
	AbortState = "Abort"
)
