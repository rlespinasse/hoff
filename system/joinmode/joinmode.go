package joinmode

// JoinMode define the mode to join multiple Nodes (source) to the same linked Node (target)
type JoinMode string

const (
	// AND will force the system to have a ComputeState at Continue
	// for all defined Nodes in order to compute the linked Node.
	AND JoinMode = "and"
	// OR will force the system to have at least on ComputeState at Continue
	// for all defined Nodes in order to compute the linked Node.
	OR = "or"
	// NONE is the default JoinMode to define a mono link between two Nodes.
	NONE = "none"
)
