package hoff

// StateType is the type of the computation state of a Node during a Computation
type StateType string

const (
	// ContinueState tell that the Node computation authorize
	// to compute the following nodes
	ContinueState StateType = "Continue"
	// SkipState tell that the Node computation is skipped
	SkipState = "Skip"
	// AbortState tell the Node computation encounter an error
	// and abort the computation
	AbortState = "Abort"
)
