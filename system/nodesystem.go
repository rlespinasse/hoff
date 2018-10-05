/*
Package system expose functions to create and manipulate a Node system.

NOTE: The NodeSystem, once correctly configure, need to be activate in order to work propertly.
*/
package system

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/internal/nodelink"
	"github.com/rlespinasse/hoff/node"
	"github.com/rlespinasse/hoff/system/joinmode"
)

// NodeSystem is a system to configure workflow between action nodes, or decision nodes.
// The nodes are linked between them by link and join mode options.
// An activated Node system will be walked throw Follow and Ancestors functions
type NodeSystem struct {
	activated      bool
	nodes          []node.Node
	nodesJoinModes map[node.Node]joinmode.JoinMode
	links          []nodelink.NodeLink

	initialNodes       []node.Node
	followingNodesTree map[node.Node]map[*bool][]node.Node
	ancestorsNodesTree map[node.Node]map[*bool][]node.Node
}

// New create an empty Node system
// who need to be valid and activated in order to be used.
func New() *NodeSystem {
	return &NodeSystem{
		activated:          false,
		nodes:              make([]node.Node, 0),
		links:              make([]nodelink.NodeLink, 0),
		nodesJoinModes:     make(map[node.Node]joinmode.JoinMode),
		initialNodes:       make([]node.Node, 0),
		followingNodesTree: make(map[node.Node]map[*bool][]node.Node),
		ancestorsNodesTree: make(map[node.Node]map[*bool][]node.Node),
	}
}

// Equal validate the two NodeSystem are equals.
func (s *NodeSystem) Equal(o *NodeSystem) bool {
	return cmp.Equal(s.activated, o.activated) && cmp.Equal(s.nodes, o.nodes, node.NodeComparator) && cmp.Equal(s.nodesJoinModes, o.nodesJoinModes) && cmp.Equal(s.links, o.links, nodelink.NodeLinkComparator)
}

// AddNode add a node to the system before activation.
func (s *NodeSystem) AddNode(n node.Node) (bool, error) {
	if s.activated {
		return false, errors.New("can't add node, node system is freeze due to activation")
	}
	s.nodes = append(s.nodes, n)
	return true, nil
}

// ConfigureJoinModeOnNode configure the join mode of a node into the system before activation.
func (s *NodeSystem) ConfigureJoinModeOnNode(n node.Node, m joinmode.JoinMode) (bool, error) {
	if s.activated {
		return false, errors.New("can't add node join mode, node system is freeze due to activation")
	}
	s.nodesJoinModes[n] = m
	return true, nil
}

// AddLink add a link from a node to another node into the system before activation.
func (s *NodeSystem) AddLink(from, to node.Node) (bool, error) {
	return s.addLink(from, to, nil)
}

// AddLinkOnBranch add a link from a node (on a specific branch) to another node into the system before activation.
func (s *NodeSystem) AddLinkOnBranch(from, to node.Node, branch bool) (bool, error) {
	return s.addLink(from, to, &branch)
}

// IsValid check if the configuration of the node system is valid based on checks.
// Check for decision node with any node links as from,
// check for cyclic redundancy in node links,
// check for undeclared node used in node links,
// check for multiple declaration of same node instance.
func (s *NodeSystem) IsValid() (bool, []error) {
	errors := make([]error, 0)
	errors = append(errors, checkForOrphanMultiBranchesNode(s)...)
	errors = append(errors, checkForCyclicRedundancyInnodeLinks(s)...)
	errors = append(errors, checkForUndeclaredNodeInnodeLink(s)...)
	errors = append(errors, checkForMultipleInstanceOfSameNode(s)...)

	if len(errors) == 0 {
		return true, nil
	}
	return false, errors
}

// Activate prepare the node system to be used.
// In order to activate it, the node system must be valid.
// Once activated, the initial nodes, following nodes, and ancestors nodes will be accessibles.
func (s *NodeSystem) Activate() error {
	if s.activated {
		return nil
	}

	validity, _ := s.IsValid()
	if !validity {
		return errors.New("can't activate a unvalidated node system")
	}

	initialNodes := make([]node.Node, 0)
	followingNodesTree := make(map[node.Node]map[*bool][]node.Node)
	ancestorsNodesTree := make(map[node.Node]map[*bool][]node.Node)

	toNodes := make([]node.Node, 0)
	for _, link := range s.links {
		followingNodesTreeOnBranch, foundNode := followingNodesTree[link.From]
		if !foundNode {
			followingNodesTree[link.From] = make(map[*bool][]node.Node)
			followingNodesTreeOnBranch = followingNodesTree[link.From]
		}
		followingNodesTreeOnBranch[link.Branch] = append(followingNodesTreeOnBranch[link.Branch], link.To)

		ancestorsNodesTreeOnBranch, foundNode := ancestorsNodesTree[link.To]
		if !foundNode {
			ancestorsNodesTree[link.To] = make(map[*bool][]node.Node)
			ancestorsNodesTreeOnBranch = ancestorsNodesTree[link.To]
		}
		ancestorsNodesTreeOnBranch[link.Branch] = append(ancestorsNodesTreeOnBranch[link.Branch], link.From)

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
	s.followingNodesTree = followingNodesTree
	s.ancestorsNodesTree = ancestorsNodesTree

	s.activated = true
	return nil
}

// JoinModeOfNode get the configured join mode of a node
func (s *NodeSystem) JoinModeOfNode(n node.Node) joinmode.JoinMode {
	mode, foundMode := s.nodesJoinModes[n]
	if foundMode {
		return mode
	}
	return joinmode.NONE
}

// InitialNodes get the initial nodes
func (s *NodeSystem) InitialNodes() []node.Node {
	return s.initialNodes
}

// IsActivated give the activation state of the node system.
// Only true if the node system is valid and have run the activate function without errors.
func (s *NodeSystem) IsActivated() bool {
	return s.activated
}

// Follow get the set of nodes accessible from a specific node and one of its branch after activation.
func (s *NodeSystem) Follow(n node.Node, branch *bool) ([]node.Node, error) {
	if !s.activated {
		return nil, errors.New("can't follow a node if system is not activated")
	}
	links, foundLinks := s.followingNodesTree[n]
	if foundLinks {
		nodes, foundNodes := links[branch]
		if foundNodes {
			return nodes, nil
		}
	}
	return nil, nil
}

// Ancestors get the set of nodes who access using one of their branch to a specific node after activation.
func (s *NodeSystem) Ancestors(n node.Node, branch *bool) ([]node.Node, error) {
	if !s.activated {
		return nil, errors.New("can't get ancestors of a node if system is not activated")
	}
	links, foundLinks := s.ancestorsNodesTree[n]
	if foundLinks {
		nodes, foundNodes := links[branch]
		if foundNodes {
			return nodes, nil
		}
	}
	return nil, nil
}

func (s *NodeSystem) addLink(from, to node.Node, branch *bool) (bool, error) {
	if s.activated {
		return false, errors.New("can't add branch link, node system is freeze due to activation")
	}

	if from == nil {
		return false, fmt.Errorf("can't have missing 'from' attribute")
	}

	if branch == nil && from.DecideCapability() {
		return false, fmt.Errorf("can't have missing branch")
	}

	if branch != nil && !from.DecideCapability() {
		return false, fmt.Errorf("can't have not needed branch")
	}

	if to == nil {
		return false, fmt.Errorf("can't have missing 'to' attribute")
	}

	if from == to {
		return false, fmt.Errorf("can't have link on from and to the same node")
	}

	if branch == nil {
		s.links = append(s.links, nodelink.New(from, to))
	} else {
		s.links = append(s.links, nodelink.NewOnBranch(from, to, *branch))
	}
	return true, nil
}

func (s *NodeSystem) haveNode(n node.Node) bool {
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
		if node.DecideCapability() {
			noLink := true
			for _, link := range s.links {
				if link.From == node {
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

func checkForCyclicRedundancyInnodeLinks(s *NodeSystem) []error {
	errors := make([]error, 0)
	nodesInCycle := make([]node.Node, 0)
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
			cyclenodeLinks := walkNodeLinks(l, l, s.links, []nodelink.NodeLink{})
			if cyclenodeLinks != nil {
				for _, link := range cyclenodeLinks {
					nodesInCycle = append(nodesInCycle, link.From)
				}
				errors = append(errors, fmt.Errorf("Can't have cycle between links: %+v", cyclenodeLinks))
			}
		}
	}
	return errors
}

func walkNodeLinks(startingLink, link nodelink.NodeLink, links []nodelink.NodeLink, walkedLinks []nodelink.NodeLink) []nodelink.NodeLink {
	nodeLinks := append(walkedLinks, link)
	if link.To == startingLink.From {
		return nodeLinks
	}
	for i := 0; i < len(links); i++ {
		if links[i] == link || links[i] == startingLink {
			continue
		}
		if link.To == links[i].From {
			return walkNodeLinks(startingLink, links[i], links, nodeLinks)
		}
	}
	return nil
}

func checkForUndeclaredNodeInnodeLink(s *NodeSystem) []error {
	errors := make([]error, 0)
	for _, link := range s.links {
		if link.From != nil && !s.haveNode(link.From) {
			errors = append(errors, fmt.Errorf("can't have undeclared node '%+v' as 'from' in branch link %+v", link.From, link))
		}
		if link.To != nil && !s.haveNode(link.To) {
			errors = append(errors, fmt.Errorf("can't have undeclared node '%+v' as 'to' in branch link %+v", link.To, link))
		}
	}
	return errors
}

func checkForMultipleInstanceOfSameNode(s *NodeSystem) []error {
	errors := make([]error, 0)
	count := make(map[node.Node]int)
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
