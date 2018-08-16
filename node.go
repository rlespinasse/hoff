package flow

type Node interface {
	Run(c *Context) RunState
	AvailableBranches() []NodeBranch
}
