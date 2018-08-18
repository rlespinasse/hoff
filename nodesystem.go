package flow

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

const (
	noBranchKey = "*"
)

type NodeLink struct {
	From   Node
	To     Node
	Branch *string
}

type NodeSystem struct {
	active   bool
	validity bool
	nodes    []Node
	links    []NodeLink

	initialNodes []Node
	nodesTree    map[Node]map[string]Node
}

func NewNodeSystem() *NodeSystem {
	return &NodeSystem{
		validity: false,
		nodes:    make([]Node, 0),
		links:    make([]NodeLink, 0),
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

func (s *NodeSystem) AddLink(n NodeLink) (bool, error) {
	if s.IsValidated() {
		return false, errors.New("can't add branch link, node system is freeze due to successful validation")
	}

	if n.From == nil {
		return false, fmt.Errorf("can't have missing 'From' attribute")
	}

	if n.Branch == nil {
		if len(n.From.AvailableBranches()) > 0 {
			return false, fmt.Errorf("can't have missing branch")
		}
	} else {
		if n.From.AvailableBranches() == nil {
			return false, fmt.Errorf("can't have not needed branch")
		}

		if len(n.From.AvailableBranches()) > 0 {
			unknonwBranch := true
			for _, branch := range n.From.AvailableBranches() {
				if branch == *n.Branch {
					unknonwBranch = false
				}
			}
			if unknonwBranch {
				return false, fmt.Errorf("can't have unknown branch")
			}
		}
	}

	if n.To == nil {
		return false, fmt.Errorf("can't have missing 'To' attribute")
	}

	s.links = append(s.links, n)
	return true, nil
}

func (s *NodeSystem) Validate() (bool, []error) {
	errors := make([]error, 0)
	errors = append(errors, checkForOrphanMultiBranchesNode(s)...)
	errors = append(errors, checkForUndeclaredNodeInNodeLink(s)...)
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

func checkForOrphanMultiBranchesNode(s *NodeSystem) []error {
	errors := make([]error, 0)
	for _, node := range s.nodes {
		if len(node.AvailableBranches()) > 0 {
			noLink := true
			for _, link := range s.links {
				if link.From == node {
					noLink = false
					break
				}
			}
			if noLink {
				errors = append(errors, fmt.Errorf("can't have orphan multi-branches node: %+v", node))
			}
		}
	}
	return errors
}

func checkForUndeclaredNodeInNodeLink(s *NodeSystem) []error {
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
