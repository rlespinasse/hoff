package flow

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
		err := cp.computeNode(node)
		if err != nil {
			return err
		}
	}
	cp.computation = true
	return nil
}

func (cp *computation) computeNode(node Node) error {
	state := node.Compute(cp.context)
	cp.report[node] = state
	switch state.value {
	case pass:
		nextNode, _ := cp.system.follow(node, state.branch)
		if nextNode == nil {
			return nil
		}
		return cp.computeNode(nextNode)
	case stop:
		if state.err != nil {
			return state.err
		}
	}
	return nil
}

func (cp *computation) isDone() bool {
	return cp.computation
}
