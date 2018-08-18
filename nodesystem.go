package flow

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

const (
	noBranchKey = "*"
)

type NodeBranchLink struct {
	From   Node
	To     Node
	Branch *string
}

type NodeSystem struct {
	active   bool
	validity bool
	nodes    []Node
	links    []NodeBranchLink

	initialNodes []Node
	nodesTree    map[Node]map[string]Node
}

func NewNodeSystem() *NodeSystem {
	return &NodeSystem{
		validity: false,
		nodes:    make([]Node, 0),
		links:    make([]NodeBranchLink, 0),
	}
}

func (x *NodeSystem) Equal(y *NodeSystem) bool {
	return cmp.Equal(x.validity, y.validity) && cmp.Equal(x.nodes, y.nodes, equalOptionForNode) && cmp.Equal(x.links, y.links, equalOptionForNode)
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
	errors = append(errors, checkForMissingFromInNodeBranchLink(s)...)
	errors = append(errors, checkForMissingToInNodeBranchLink(s)...)
	errors = append(errors, checkForMissingBranchInNodeBranchLink(s)...)
	errors = append(errors, checkForNotNeededBranchInNodeBranchLink(s)...)
	errors = append(errors, checkForUnknownBranchInNodeBranchLink(s)...)
	errors = append(errors, checkForUndeclaredNodeInNodeBranchLink(s)...)
	errors = append(errors, checkForMultipleInstanceOfSameNode(s)...)

	s.validity = len(errors) < 1
	if s.IsValidated() {
		return true, nil
	}
	return false, errors
}

func (s *NodeSystem) activate() error {
	if s.active {
		return nil
	}
	if !s.validity {
		return errors.New("can't activate a unvalidated node system")
	}

	initialNodes := make([]Node, 0)
	nodesTree := make(map[Node]map[string]Node)

	toNodes := make([]Node, 0)
	for _, link := range s.links {
		branch := noBranchKey
		if link.Branch != nil {
			branch = *link.Branch
		}
		fromNodeTree, foundFromNode := nodesTree[link.From]
		if !foundFromNode {
			nodesTree[link.From] = make(map[string]Node)
			fromNodeTree, _ = nodesTree[link.From]
		}
		fromNodeTree[branch] = link.To
		toNodes = append(toNodes, link.To)
	}

	for _, node := range s.nodes {
		isInitialNode := true
		for _, toNode := range toNodes {
			if node == toNode {
				isInitialNode = false
				break
			}
		}
		if isInitialNode {
			initialNodes = append(initialNodes, node)
		}
	}

	s.initialNodes = initialNodes
	s.nodesTree = nodesTree

	s.active = true
	return nil
}

func (s *NodeSystem) follow(n Node, b *string) (Node, error) {
	if !s.active {
		return nil, errors.New("can't follow a node if system is not activated")
	}
	links, foundLinks := s.nodesTree[n]
	if foundLinks {
		branch := noBranchKey
		if b != nil {
			branch = *b
		}
		node, foundNodes := links[branch]
		if foundNodes {
			return node, nil
		}
	}
	return nil, nil
}

func (s *NodeSystem) InitialNodes() []Node {
	return s.initialNodes
}

func (s *NodeSystem) IsValidated() bool {
	return s.validity
}

func (s *NodeSystem) haveNode(n Node) bool {
	for _, node := range s.nodes {
		if node == n {
			return true
		}
	}
	return false
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
				errors = append(errors, fmt.Errorf("can't have orphan decision node: %+v", node))
			}
		}
	}
	return errors
}

func checkForMissingFromInNodeBranchLink(s *NodeSystem) []error {
	errors := make([]error, 0)
	for _, link := range s.links {
		if link.From == nil {
			errors = append(errors, fmt.Errorf("can't have missing 'From' attribute: %+v", link))
		}
	}
	return errors
}

func checkForMissingToInNodeBranchLink(s *NodeSystem) []error {
	errors := make([]error, 0)
	for _, link := range s.links {
		if link.To == nil {
			errors = append(errors, fmt.Errorf("can't have missing 'To' attribute: %+v", link))
		}
	}
	return errors
}

func checkForMissingBranchInNodeBranchLink(s *NodeSystem) []error {
	errors := make([]error, 0)
	for _, link := range s.links {
		if link.From != nil && link.Branch == nil && len(link.From.AvailableBranches()) > 0 {
			errors = append(errors, fmt.Errorf("can't have missing branch from %+v, available branches %+v", link.From, link.From.AvailableBranches()))
		}
	}
	return errors
}

func checkForUnknownBranchInNodeBranchLink(s *NodeSystem) []error {
	errors := make([]error, 0)
	for _, link := range s.links {
		if link.From != nil && link.Branch != nil && len(link.From.AvailableBranches()) > 0 {
			unknonwBranch := true
			for _, branch := range link.From.AvailableBranches() {
				if branch == *link.Branch {
					unknonwBranch = false
				}
			}
			if unknonwBranch {
				errors = append(errors, fmt.Errorf("can't have unknown branch: '%v', from %+v, available branches %+v", *link.Branch, link.From, link.From.AvailableBranches()))
			}
		}
	}
	return errors
}

func checkForNotNeededBranchInNodeBranchLink(s *NodeSystem) []error {
	errors := make([]error, 0)
	for _, link := range s.links {
		if link.From != nil && link.Branch != nil && link.From.AvailableBranches() == nil {
			errors = append(errors, fmt.Errorf("can't have not needed branch: '%v', from %+v", *link.Branch, link.From))
		}
	}
	return errors
}

func checkForUndeclaredNodeInNodeBranchLink(s *NodeSystem) []error {
	errors := make([]error, 0)
	for _, link := range s.links {
		if link.From != nil && !s.haveNode(link.From) {
			errors = append(errors, fmt.Errorf("can't have undeclared node '%+v' as 'From' in branch link %+v", link.From, link))
		}
		if link.To != nil && !s.haveNode(link.To) {
			errors = append(errors, fmt.Errorf("can't have undeclared node '%+v' as 'To' in branch link %+v", link.To, link))
		}
	}
	return errors
}

func checkForMultipleInstanceOfSameNode(s *NodeSystem) []error {
	errors := make([]error, 0)
	count := make(map[Node]int)
	for i := 0; i < len(s.nodes); i++ {
		for j := 0; j < len(s.nodes); j++ {
			if i != j && s.nodes[i] == s.nodes[j] {
				count[s.nodes[i]]++
			}
		}
	}
	for n, c := range count {
		if c > 1 {
			errors = append(errors, fmt.Errorf("can't have multiple instances (%v) of the same node: %+v", c, n))
		}
	}
	return errors
}
