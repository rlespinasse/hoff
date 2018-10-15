/*
Package engine serve to run multiple computations against a node system.

Supported computations:

	sequential

*/
package engine

import (
	"errors"

	"github.com/rlespinasse/hoff/computation"
	"github.com/rlespinasse/hoff/computestate"
	"github.com/rlespinasse/hoff/node"

	"github.com/rlespinasse/hoff/system"
)

// Engine expose an engine to manage multiple computations based on a node system.
type Engine struct {
	mode   ComputationMode
	system *system.NodeSystem
}

// New create an engine with computation mode.
// Need to be configured with a node system
func New(mode ComputationMode) *Engine {
	return &Engine{
		mode: mode,
	}
}

// ConfigureNodeSystem add a node system to the engine (only once).
func (e *Engine) ConfigureNodeSystem(system *system.NodeSystem) error {
	if e.system != nil {
		return errors.New("node system already configured")
	}
	if !system.IsActivated() {
		return errors.New("node system need to be activated")
	}
	e.system = system
	return nil
}

// Compute run computation against node system with input data.
func (e *Engine) Compute(data map[string]interface{}) ComputationResult {
	if e.system == nil {
		return ComputationResult{
			Data:  data,
			Error: errors.New("need a configured node system"),
		}
	}

	cp, _ := computation.New(e.system, node.NewContext(data))

	err := cp.Compute()
	return ComputationResult{
		Data:   cp.Context.Data,
		Error:  err,
		Report: cp.Report,
	}
}

// ComputationResult store the result of a computation.
type ComputationResult struct {
	Error  error
	Data   map[string]interface{}
	Report map[node.Node]computestate.ComputeState
}
