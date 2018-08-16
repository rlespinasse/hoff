package flowengine

type FlowNode interface {
	Run(c *FlowContext)
}
