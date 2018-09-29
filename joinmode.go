package hoff

type JoinMode string

const (
	JoinModeAnd  JoinMode = "and"
	JoinModeOr            = "or"
	JoinModeNone          = "none"
)
