package flows

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

func (n NodeLink) String() string {
	branch := ""
	if n.Branch != nil {
		branch = fmt.Sprintf(", Branch: %v", *n.Branch)
	}
	return fmt.Sprintf("{From: %v, To: %v%v}", n.From, n.To, branch)
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

	if n.From == n.To {
		return false, fmt.Errorf("can't have link on from and to the same node")
	}

	s.links = append(s.links, n)
	return true, nil
}

func (s *NodeSystem) Validate() (bool, []error) {
	errors := make([]error, 0)
	errors = append(errors, checkForOrphanMultiBranchesNode(s)...)
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

func checkForMultipleLinksWithSameTo(s *NodeSystem) []error {
	errors := make([]error, 0)
	sameToNodes := make(map[Node][]NodeLink)
	for i := 0; i < len(s.links); i++ {
		for j := 0; j < len(s.links); j++ {
			if i != j && s.links[i].To == s.links[j].To {
				sameToNodes[s.links[i].To] = append(sameToNodes[s.links[i].To], s.links[i])
			}
		}
	}
	for _, links := range sameToNodes {
		if len(links) > 1 {
			errors = append(errors, fmt.Errorf("Can't have multiple links with the same 'To': %+v", links))
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
			if node == l.From {
				notInCycle = false
				break
			}
		}
		if notInCycle {
			cycleNodeLinks := walkNodeLink(l, l, s.links, []NodeLink{})
			if cycleNodeLinks != nil {
				for _, link := range cycleNodeLinks {
					nodesInCycle = append(nodesInCycle, link.From)
				}
				errors = append(errors, fmt.Errorf("Can't have cycle between links: %+v", cycleNodeLinks))
			}
		}
	}
	return errors
}

func walkNodeLink(startingLink NodeLink, link NodeLink, links []NodeLink, walkedNodeLinks []NodeLink) []NodeLink {
	nodeLinks := append(walkedNodeLinks, link)
	if link.To == startingLink.From {
		return nodeLinks
	}
	for i := 0; i < len(links); i++ {
		if links[i] == link || links[i] == startingLink {
			continue
		}
		if link.To == links[i].From {
			return walkNodeLink(startingLink, links[i], links, nodeLinks)
		}
	}
	return nil
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
