package namingishard

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

type NodeSystem struct {
	active         bool
	validity       bool
	nodes          []Node
	nodesJoinModes map[Node]JoinMode
	links          []NodeLink

	initialNodes []Node
	linkTree     map[Node]map[*bool][]Node
	linkedTo     map[Node][]Node
}

func NewNodeSystem() *NodeSystem {
	return &NodeSystem{
		validity:       false,
		nodes:          make([]Node, 0),
		nodesJoinModes: make(map[Node]JoinMode),
		links:          make([]NodeLink, 0),
	}
}

func (x *NodeSystem) Equal(y *NodeSystem) bool {
	return cmp.Equal(x.validity, y.validity) && cmp.Equal(x.nodes, y.nodes, equalOptionForNode) && cmp.Equal(x.links, y.links, equalOptionForNodeLink)
}

func (s *NodeSystem) AddNode(n Node) (bool, error) {
	if s.IsValidated() {
		return false, errors.New("can't add node, node system is freeze due to successful validation")
	}
	s.nodes = append(s.nodes, n)
	return true, nil
}

func (s *NodeSystem) AddNodeJoinMode(n Node, m JoinMode) (bool, error) {
	if s.IsValidated() {
		return false, errors.New("can't add node join mode, node system is freeze due to successful validation")
	}
	s.nodesJoinModes[n] = m
	return true, nil
}

func (s *NodeSystem) AddLink(n NodeLink) (bool, error) {
	if s.IsValidated() {
		return false, errors.New("can't add branch link, node system is freeze due to successful validation")
	}

	if n.from == nil {
		return false, fmt.Errorf("can't have missing 'from' attribute")
	}

	if n.branch == nil && n.from.decideCapability() {
		return false, fmt.Errorf("can't have missing branch")
	}

	if n.branch != nil && !n.from.decideCapability() {
		return false, fmt.Errorf("can't have not needed branch")
	}

	if n.to == nil {
		return false, fmt.Errorf("can't have missing 'to' attribute")
	}

	if n.from == n.to {
		return false, fmt.Errorf("can't have link on from and to the same node")
	}

	s.links = append(s.links, n)
	return true, nil
}

func (s *NodeSystem) Validate() (bool, []error) {
	errors := make([]error, 0)
	errors = append(errors, checkForOrphanMultiBranchesNode(s)...)
	errors = append(errors, checkForCyclicRedundancyInNodeLinks(s)...)
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
	linkTree := make(map[Node]map[*bool][]Node)
	linkedTo := make(map[Node][]Node)

	toNodes := make([]Node, 0)
	for _, link := range s.links {
		linkBranchTree, foundNode := linkTree[link.from]

		if !foundNode {
			linkTree[link.from] = make(map[*bool][]Node)
			linkBranchTree, _ = linkTree[link.from]
		}

		linkBranchTree[link.branch] = append(linkBranchTree[link.branch], link.to)
		linkedTo[link.to] = append(linkedTo[link.to], link.from)
		toNodes = append(toNodes, link.to)
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
	s.linkTree = linkTree
	s.linkedTo = linkedTo

	s.active = true
	return nil
}

func (s *NodeSystem) NodeJoinMode(n Node) JoinMode {
	mode, foundMode := s.nodesJoinModes[n]
	if foundMode {
		return mode
	}
	return JoinModeNone
}

func (s *NodeSystem) follow(n Node, branch *bool) ([]Node, error) {
	if !s.active {
		return nil, errors.New("can't follow a node if system is not activated")
	}
	links, foundLinks := s.linkTree[n]
	if foundLinks {
		nodes, foundNodes := links[branch]
		if foundNodes {
			return nodes, nil
		}
	}
	return nil, nil
}

func (s *NodeSystem) nodesLinkedTo(n Node) []Node {
	return s.linkedTo[n]
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
		if node.decideCapability() {
			noLink := true
			for _, link := range s.links {
				if link.from == node {
					noLink = false
					break
				}
			}
			if noLink {
				errors = append(errors, fmt.Errorf("can't have decision node without link from it: %+v", node))
			}
		}
	}
	return errors
}

func checkForCyclicRedundancyInNodeLinks(s *NodeSystem) []error {
	errors := make([]error, 0)
	nodesInCycle := make([]Node, 0)
	for i := 0; i < len(s.links); i++ {
		l := s.links[i]
		notInCycle := true
		for _, node := range nodesInCycle {
			if node == l.from {
				notInCycle = false
				break
			}
		}
		if notInCycle {
			cycleNodeLinks := walkNodeLink(l, l, s.links, []NodeLink{})
			if cycleNodeLinks != nil {
				for _, link := range cycleNodeLinks {
					nodesInCycle = append(nodesInCycle, link.from)
				}
				errors = append(errors, fmt.Errorf("Can't have cycle between links: %+v", cycleNodeLinks))
			}
		}
	}
	return errors
}

func walkNodeLink(startingLink NodeLink, link NodeLink, links []NodeLink, walkedNodeLinks []NodeLink) []NodeLink {
	nodeLinks := append(walkedNodeLinks, link)
	if link.to == startingLink.from {
		return nodeLinks
	}
	for i := 0; i < len(links); i++ {
		if links[i] == link || links[i] == startingLink {
			continue
		}
		if link.to == links[i].from {
			return walkNodeLink(startingLink, links[i], links, nodeLinks)
		}
	}
	return nil
}

func checkForUndeclaredNodeInNodeLink(s *NodeSystem) []error {
	errors := make([]error, 0)
	for _, link := range s.links {
		if link.from != nil && !s.haveNode(link.from) {
			errors = append(errors, fmt.Errorf("can't have undeclared node '%+v' as 'from' in branch link %+v", link.from, link))
		}
		if link.to != nil && !s.haveNode(link.to) {
			errors = append(errors, fmt.Errorf("can't have undeclared node '%+v' as 'to' in branch link %+v", link.to, link))
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
