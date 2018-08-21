package namingishard

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

const (
	noBranchKey = "*"
)

type NodeSystem struct {
	active   bool
	validity bool
	nodes    []Node
	links    []NodeLink

	initialNodes []Node
	kindTree     map[Node]map[string]linkKind
	linkTree     map[Node]map[string][]Node
	linkedTo     map[Node][]Node
}

func NewNodeSystem() *NodeSystem {
	return &NodeSystem{
		validity: false,
		nodes:    make([]Node, 0),
		links:    make([]NodeLink, 0),
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

func (s *NodeSystem) AddLink(n NodeLink) (bool, error) {
	if s.IsValidated() {
		return false, errors.New("can't add branch link, node system is freeze due to successful validation")
	}

	if n.from == nil {
		return false, fmt.Errorf("can't have missing 'from' attribute")
	}

	if n.branch == nil {
		if len(n.from.AvailableBranches()) > 0 {
			return false, fmt.Errorf("can't have missing branch")
		}
	} else {
		if n.from.AvailableBranches() == nil {
			return false, fmt.Errorf("can't have not needed branch")
		}

		if len(n.from.AvailableBranches()) > 0 {
			unknonwBranch := true
			for _, branch := range n.from.AvailableBranches() {
				if branch == *n.branch {
					unknonwBranch = false
				}
			}
			if unknonwBranch {
				return false, fmt.Errorf("can't have unknown branch")
			}
		}
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
	errors = append(errors, checkForMultipleLinksWithSameFrom(s)...)
	errors = append(errors, checkForMultipleLinksWithSameTo(s)...)
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
	linkTree := make(map[Node]map[string][]Node)
	kindTree := make(map[Node]map[string]linkKind)
	linkedTo := make(map[Node][]Node)

	toNodes := make([]Node, 0)
	for _, link := range s.links {
		branch := noBranchKey
		if link.branch != nil {
			branch = *link.branch
		}
		linkBranchTree, foundNode := linkTree[link.from]
		kindBranchTree, _ := kindTree[link.from]
		if !foundNode {
			linkTree[link.from] = make(map[string][]Node)
			kindTree[link.from] = make(map[string]linkKind)
			linkBranchTree, _ = linkTree[link.from]
			kindBranchTree, _ = kindTree[link.from]
		}
		linkBranchTree[branch] = append(linkBranchTree[branch], link.to)
		kindBranchTree[branch] = link.kind
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
	s.kindTree = kindTree
	s.linkedTo = linkedTo

	s.active = true
	return nil
}

func (s *NodeSystem) follow(n Node, b *string) ([]Node, linkKind, error) {
	if !s.active {
		return nil, noLink, errors.New("can't follow a node if system is not activated")
	}
	links, foundLinks := s.linkTree[n]
	if foundLinks {
		kinds, _ := s.kindTree[n]
		branch := noBranchKey
		if b != nil {
			branch = *b
		}
		nodes, foundNodes := links[branch]
		if foundNodes {
			kind, _ := kinds[branch]
			return nodes, kind, nil
		}
	}
	return nil, noLink, nil
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
		if len(node.AvailableBranches()) > 0 {
			noLink := true
			for _, link := range s.links {
				if link.from == node {
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

func checkForMultipleLinksWithSameFrom(s *NodeSystem) []error {
	errors := make([]error, 0)
	sameFromNodes := make(map[Node][]NodeLink)
	for i := 0; i < len(s.links); i++ {
		for j := 0; j < len(s.links); j++ {
			if i != j && s.links[i].from == s.links[j].from {
				sameFromNodes[s.links[i].from] = append(sameFromNodes[s.links[i].from], s.links[i])
			}
		}
	}
	for _, links := range sameFromNodes {
		if len(links) > 1 {
			kind := links[0].kind
			for _, link := range links {
				if kind != link.kind {
					errors = append(errors, fmt.Errorf("Can't have mixed kind links with the same 'from': %+v", links))
					break
				}
			}
		}
	}
	return errors
}

func checkForMultipleLinksWithSameTo(s *NodeSystem) []error {
	errors := make([]error, 0)
	sameToNodes := make(map[Node][]NodeLink)
	for i := 0; i < len(s.links); i++ {
		for j := 0; j < len(s.links); j++ {
			if i != j && s.links[i].to == s.links[j].to {
				sameToNodes[s.links[i].to] = append(sameToNodes[s.links[i].to], s.links[i])
			}
		}
	}
	for _, links := range sameToNodes {
		if len(links) > 1 {
			kind := links[0].kind
			for _, link := range links {
				if kind != link.kind {
					errors = append(errors, fmt.Errorf("Can't have mixed kind links with the same 'to': %+v", links))
					break
				}
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
