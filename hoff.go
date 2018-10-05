/*
Package hoff is a library to define a Node workflow and compute data against it.

Create a node system and activate it:

	nodesystem := system.New()
	nodesystem.AddNode(some_action_node) // read the input_data in context and create a output_data in context
	nodesystem.AddNode(decision_node) // check if the input_data is valid with some functionals rules
	nodesystem.AddNode(another_action_node) // enhance output_data with some functionals tasks
	nodesystem.ConfigureJoinModeOnNode(another_action_node, joinmode.AND)
	nodesystem.AddLink(some_action_node, another_action_node)
	nodesystem.AddLinkOnBranch(decision_node, another_action_node, true)
	nodesystem.Activate()

Create a computation and launch it:

	context := node.NewContextWithoutData()
	context.Store("input_info", input_info)
	computation := computation.New(nodesystem, context)
	err := computation.Compute()
	if err != nil {
		// error handling
	}
	fmt.Printf("computation report: %+v", computation.Report)
	fmt.Printf("output_data: %+v", computation.Context.Read("output_data"))

*/
package hoff
