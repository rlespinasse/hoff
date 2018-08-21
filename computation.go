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
		err := cp.computeNode(node, noLink)
		if err != nil {
			return err
		}
	}
	cp.computation = true
	return nil
}

func (cp *computation) computeNode(node Node, linkKind linkKind) error {
	if linkKind == joinLink {
		linkedNodes := cp.system.nodesLinkedTo(node)
		for _, linkedNode := range linkedNodes {
			if _, found := cp.report[linkedNode]; !found {
				return nil
			}
		}
	} else if report, found := cp.report[node]; found {
		return report.err
	}
	state := node.Compute(cp.context)
	cp.report[node] = state
	if state.value == pass {
		nextNodes, kind, _ := cp.system.follow(node, state.branch)
		switch kind {
		case forkLink:
			for _, node := range nextNodes {
				err := cp.computeNode(node, forkLink)
				if err != nil {
					return err
				}
			}
		default:
			if nextNodes != nil {
				return cp.computeNode(nextNodes[0], kind)
			}
		}
	}
	return state.err
}

func (cp *computation) isDone() bool {
	return cp.computation
}
