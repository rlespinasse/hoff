package hoff

// JoinMode define the mode to join multiple Nodes (source) to the same linked Node (target)
type JoinMode string

const (
	// JoinAnd will force the system to have a ComputeState at Continue
	// for all defined Nodes in order to compute the linked Node.
	JoinAnd JoinMode = "and"
	// JoinOr will force the system to have at least on ComputeState at Continue
	// for all defined Nodes in order to compute the linked Node.
	JoinOr = "or"
	// JoinNone is the default JoinMode to define a mono link between two Nodes.
	JoinNone = "none"
)
