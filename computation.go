package namingishard

import (
	"errors"

	"github.com/google/go-cmp/cmp"
)

type computation struct {
	computation bool
	system      *NodeSystem
	context     *Context
	report      map[Node]ComputeState
}

func (x computation) Equal(y computation) bool {
	return cmp.Equal(x.computation, y.computation) && cmp.Equal(x.context, y.context) && cmp.Equal(x.system, y.system) && cmp.Equal(x.report, y.report)
}

func NewComputation(system *NodeSystem, context *Context) (*computation, error) {
	if system == nil {
		return nil, errors.New("must have a node system to work properly")
	}
	if !system.IsValidated() {
		return nil, errors.New("must have a validated node system to work properly")
	}
	if context == nil {
		return nil, errors.New("must have a context to work properly")
	}
	return &computation{
		computation: false,
		system:      system,
		context:     context,
	}, nil
}

func (cp *computation) Compute() error {
	cp.report = make(map[Node]ComputeState)
	for _, node := range cp.system.InitialNodes() {
		err := cp.computeNode(node, nil, nil)
		if err != nil {
			return err
		}
	}
	cp.computation = true
	return nil
}

func (cp *computation) computeNode(node, callingNode Node, branch *bool) error {
	if callingNode != nil {
		joinMode := cp.system.NodeJoinMode(node)
		if joinMode != JoinModeNone {
			linkedNodes := cp.system.nodesLinkedTo(node)
			passingNodesCount := 0
			for _, linkedNode := range linkedNodes {
				report, found := cp.report[linkedNode]
				if !found {
					return nil
				} else if report.value == pass && report.branch == branch {
					passingNodesCount++
				}
			}

			switch joinMode {
			case JoinModeAnd:
				if passingNodesCount != len(linkedNodes) {
					return nil
				}
			}
		}
	}
	state := node.Compute(cp.context)
	cp.report[node] = state
	if state.value == pass {
		nextNodes, _ := cp.system.follow(node, state.branch)
		for _, newNode := range nextNodes {
			err := cp.computeNode(newNode, node, state.branch)
			if err != nil {
				return err
			}
		}
	}
	return state.err
}

func (cp *computation) isDone() bool {
	return cp.computation
}
