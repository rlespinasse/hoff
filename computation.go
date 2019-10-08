package hoff

import (
	"errors"

	"github.com/rlespinasse/hoff/computestate"
	"github.com/rlespinasse/hoff/internal/utils"
	"github.com/rlespinasse/hoff/statetype"

	"github.com/google/go-cmp/cmp"
)

// Computation take a NodeSystem and compute a Context against it.
type Computation struct {
	System  *NodeSystem
	Context *Context
	Status  bool
	Report  map[Node]computestate.ComputeState
}

// NewComputation create a computation based on a valid, and activated NodeSystem and a Context.
func NewComputation(system *NodeSystem, context *Context) (*Computation, error) {
	if system == nil {
		return nil, errors.New("must have a node system to work properly")
	}
	if !system.IsActivated() {
		return nil, errors.New("must have an activated node system to work properly")
	}
	if context == nil {
		return nil, errors.New("must have a context to work properly")
	}
	return &Computation{
		Status:  false,
		System:  system,
		Context: context,
	}, nil
}

// Equal validate the two Computation are equals.
func (cp Computation) Equal(o Computation) bool {
	return cmp.Equal(cp.Status, o.Status) && cmp.Equal(cp.Context, o.Context) && cmp.Equal(cp.System, o.System) && cmp.Equal(cp.Report, o.Report)
}

// Compute run all nodes in the defined order to enhance the Context.
// At the end of the computation (Status at true), you can read the compute state
// of each node in the Report.
func (cp *Computation) Compute() error {
	cp.Report = make(map[Node]computestate.ComputeState)
	err := cp.computeNodes(cp.System.InitialNodes())
	if err != nil {
		return err
	}
	cp.Status = true
	return nil
}

func (cp *Computation) computeNodes(nodes []Node) error {
	for _, node := range nodes {
		err := cp.computeNode(node)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cp *Computation) computeNode(node Node) error {
	order := cp.calculateComputeOrder(node)

	switch order {
	case dontRunIt, alreadyRunOnce:
		return nil
	case skipIt:
		cp.Report[node] = computestate.Skip()
	case computeIt:
		state := node.Compute(cp.Context)
		cp.Report[node] = state
		if state.Value == statetype.AbortState {
			return state.Error
		}
	}

	return cp.computeFollowingNodes(node, nodeBranches(node)...)
}

func (cp *Computation) computeFollowingNodes(node Node, branches ...*bool) error {
	for _, branch := range branches {
		nextNodes, _ := cp.System.Follow(node, branch)
		err := cp.computeNodes(nextNodes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cp *Computation) calculateComputeOrder(node Node) computeOrder {
	if _, ok := cp.Report[node]; ok {
		return alreadyRunOnce
	}

	ancestorsCount, ancestorsComputed, ancestorsWithContinueState := cp.ansectorsComputationStatistics(node)
	if ancestorsCount != ancestorsComputed {
		return dontRunIt
	} else if ancestorsCount == 0 {
		return computeIt
	}

	joinMode := cp.System.JoinModeOfNode(node)
	switch joinMode {
	case JoinAnd:
		if ancestorsCount == ancestorsWithContinueState {
			return computeIt
		}
	case JoinOr:
		if ancestorsWithContinueState > 0 {
			return computeIt
		}
	case JoinNone:
		if ancestorsWithContinueState == 1 {
			return computeIt
		}
	}
	return skipIt
}

func (cp *Computation) ansectorsComputationStatistics(node Node) (int, int, int) {
	ancestorsCount, ancestorsComputed, ancestorsWithContinueState := cp.ansectorsComputationStatisticsOnBranch(node, nil)

	ancestorsCountOnBranchTrue, ancestorsComputedOnBranchTrue, ancestorsWithContinueStateOnBranchTrue := cp.ansectorsComputationStatisticsOnBranch(node, utils.BoolPointer(true))
	ancestorsCount += ancestorsCountOnBranchTrue
	ancestorsComputed += ancestorsComputedOnBranchTrue
	ancestorsWithContinueState += ancestorsWithContinueStateOnBranchTrue

	ancestorsCountOnBranchFalse, ancestorsComputedOnBranchFalse, ancestorsWithContinueStateOnBranchFalse := cp.ansectorsComputationStatisticsOnBranch(node, utils.BoolPointer(false))
	ancestorsCount += ancestorsCountOnBranchFalse
	ancestorsComputed += ancestorsComputedOnBranchFalse
	ancestorsWithContinueState += ancestorsWithContinueStateOnBranchFalse

	return ancestorsCount, ancestorsComputed, ancestorsWithContinueState
}

func (cp *Computation) ansectorsComputationStatisticsOnBranch(node Node, branch *bool) (int, int, int) {
	linkedNodes, _ := cp.System.Ancestors(node, branch)
	computedNodes := 0
	nodesWithContinueState := 0
	for _, linkedNode := range linkedNodes {
		report, found := cp.Report[linkedNode]
		if found {
			computedNodes++
			if report.Value == statetype.ContinueState && report.Branch == branch {
				nodesWithContinueState++
			}
		}
	}
	return len(linkedNodes), computedNodes, nodesWithContinueState
}

type computeOrder string

const (
	computeIt      computeOrder = "compute_it"
	skipIt                      = "skip_it"
	dontRunIt                   = "dont_run_it"
	alreadyRunOnce              = "already_run_once"
)

func nodeBranches(node Node) []*bool {
	if node.DecideCapability() {
		return []*bool{utils.BoolPointer(true), utils.BoolPointer(false)}
	}
	return []*bool{nil}
}
