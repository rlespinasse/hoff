/*
Package hoff is a library to define a Node workflow and compute data against it.

Create a node system and activate it:

	ns := hoff.NewNodeSystem()
	ns.AddNode(some_action_node) // read the input_data in context and create a output_data in context
	ns.AddNode(decision_node) // check if the input_data is valid with some functionals rules
	ns.AddNode(another_action_node) // enhance output_data with some functionals tasks
	ns.ConfigureJoinModeOnNode(another_action_node, hoff.JoinAnd)
	ns.AddLink(some_action_node, another_action_node)
	ns.AddLinkOnBranch(decision_node, another_action_node, true)
	errs := ns.Activate()
	if errs != nil {
		// error handling
	}

Create a single computation and launch it:

	cxt := hoff.NewContextWithoutData()
	cxt.Store("input_info", input_info)
	cp := hoff.NewComputation(ns, context)
	err := cp.Compute()
	if err != nil {
		// error handling
	}
	fmt.Printf("computation report: %+v", cp.Report)
	fmt.Printf("output_data: %+v", cp.Context.Read("output_data"))

Create an engine and run multiple computations:

	eng := hoff.NewEngine(hoff.SequentialComputation)
	eng.ConfigureNodeSystem(ns)

	cr1 := eng.Compute(input_info)
	fmt.Printf("computation error: %+v", cr1.Error)
	fmt.Printf("computation report: %+v", cr1.Report)
	fmt.Printf("computed data: %+v", cr1.Data)

	cr2 := eng.Compute(another_aninput_info)
	fmt.Printf("computation error: %+v", cr2.Error)
	fmt.Printf("computation report: %+v", cr2.Report)
	fmt.Printf("computed data: %+v", cr2.Data)

*/
package hoff
