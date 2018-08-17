package flow

import (
	"errors"
)

type computation struct {
	computation bool
	system      *NodeSystem
	context     *Context
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
