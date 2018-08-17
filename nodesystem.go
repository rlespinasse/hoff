package flow

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

type NodeBranchLink struct {
	From   Node
	To     Node
	Branch *string
}

type NodeSystem struct {
	validity bool
	nodes    []Node
	links    []NodeBranchLink
}

func NewNodeSystem() *NodeSystem {
	return &NodeSystem{
		validity: false,
		nodes:    make([]Node, 0),
		links:    make([]NodeBranchLink, 0),
	}
}

func (x *NodeSystem) Equal(y *NodeSystem) bool {
	return x.validity == y.validity && cmp.Equal(x.nodes, y.nodes, nodeEqualOpts) && cmp.Equal(x.links, y.links, nodeEqualOpts)
}

func (s *NodeSystem) AddNode(n Node) (bool, error) {
	if s.IsValidated() {
		return false, errors.New("can't add node, node system is freeze due to successful validation")
	}
	s.nodes = append(s.nodes, n)
	return true, nil
}

func (s *NodeSystem) AddBranchLink(n NodeBranchLink) (bool, error) {
	if s.IsValidated() {
		return false, errors.New("can't add branch link, node system is freeze due to successful validation")
	}
	s.links = append(s.links, n)
	return true, nil
}

func (s *NodeSystem) Validate() (bool, []error) {
	errors := make([]error, 0)
	errors = append(errors, checkForOrphanDecisionNode(s)...)
	errors = append(errors, validateNodeBranchLink(s)...)

	s.validity = len(errors) < 1
	if s.IsValidated() {
		return true, nil
	}
	return false, errors
}

func (s *NodeSystem) IsValidated() bool {
	return s.validity
}

func checkForOrphanDecisionNode(s *NodeSystem) []error {
	errors := make([]error, 0)
	for _, node := range s.nodes {
		if isDecisionNode(node) {
			noLink := true
			for _, link := range s.links {
				if link.From == node {
					noLink = false
					break
				}
			}
			if noLink {
				errors = append(errors, fmt.Errorf("orphan decision node: %+v", node))
			}
		}
	}
	return errors
}

func validateNodeBranchLink(s *NodeSystem) []error {
	errors := make([]error, 0)
	for _, link := range s.links {
		if link.From == nil {
			errors = append(errors, fmt.Errorf("missing 'From' attribute: %+v", link))
		} else if link.Branch == nil {
			if len(link.From.AvailableBranches()) > 0 {
				errors = append(errors, fmt.Errorf("missing branch from %+v, available branches %+v", link.From, link.From.AvailableBranches()))
			}
		} else if link.Branch != nil {
			if len(link.From.AvailableBranches()) > 0 {
				unknonwBranch := true
				for _, branch := range link.From.AvailableBranches() {
					if branch == *link.Branch {
						unknonwBranch = false
					}
				}
				if unknonwBranch {
					errors = append(errors, fmt.Errorf("unknown branch: '%v', from %+v, available branches %+v", *link.Branch, link.From, link.From.AvailableBranches()))
				}
			} else {
				errors = append(errors, fmt.Errorf("not needed branch: '%v', from %+v", *link.Branch, link.From))
			}
		}
		if link.To == nil {
			errors = append(errors, fmt.Errorf("missing 'To' attribute: %+v", link))
		}
	}
	return errors
}
