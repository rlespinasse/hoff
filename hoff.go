/*
Package hoff is a library to define a Node workflow and compute data against it.

Create a node system and activate it:

	ns := system.New()
	ns.AddNode(some_action_node) // read the input_data in context and create a output_data in context
	ns.AddNode(decision_node) // check if the input_data is valid with some functionals rules
	ns.AddNode(another_action_node) // enhance output_data with some functionals tasks
	ns.ConfigureJoinModeOnNode(another_action_node, joinmode.AND)
	ns.AddLink(some_action_node, another_action_node)
	ns.AddLinkOnBranch(decision_node, another_action_node, true)
	errs := ns.Activate()
	if errs != nil {
		// error handling
	}

Create a computation and launch it:

	cxt := node.NewContextWithoutData()
	cxt.Store("input_info", input_info)
	cp := computation.New(ns, context)
	err := cp.Compute()
	if err != nil {
		// error handling
	}
	fmt.Printf("computation report: %+v", cp.Report)
	fmt.Printf("output_data: %+v", cp.Context.Read("output_data"))

*/
package hoff
