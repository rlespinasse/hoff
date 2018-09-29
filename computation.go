package hoff

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

func (cp *computation) computeNode(node, callingNode Node, callingBranch *bool) error {
	skipIt := false
	if callingNode != nil {
		joinMode := cp.system.NodeJoinMode(node)
		if joinMode != JoinModeNone {
			linkedNodes := cp.system.nodesLinkedTo(node)
			passingNodesCount := 0
			for _, linkedNode := range linkedNodes {
				report, found := cp.report[linkedNode]
				if !found {
					return nil
				} else if report.value == Continue && report.branch == callingBranch {
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

		callingState := cp.report[callingNode]
		if !(callingState.value == Continue && callingState.branch == callingBranch) {
			skipIt = true
		}
	}

	if !skipIt {
		state := node.Compute(cp.context)
		cp.report[node] = state

		if state.value == Abort {
			return state.err
		}

	} else {
		cp.report[node] = ComputeStateSkip()
	}

	if node.decideCapability() {
		err := cp.computeNextNodes(node, truePointer)
		if err != nil {
			return err
		}
		err = cp.computeNextNodes(node, falsePointer)
		if err != nil {
			return err
		}
	} else {
		err := cp.computeNextNodes(node, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cp *computation) computeNextNodes(callingNode Node, followBranch *bool) error {
	nextNodes, _ := cp.system.follow(callingNode, followBranch)
	for _, node := range nextNodes {
		err := cp.computeNode(node, callingNode, followBranch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cp *computation) isDone() bool {
	return cp.computation
}
