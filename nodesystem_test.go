package hoff

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/internal/utils"
)

var (
	someActionNode, _         = NewActionNode("someActionNode", func(*Context) error { return nil })
	anotherActionNode, _      = NewActionNode("anotherActionNode", func(*Context) error { return nil })
	alwaysTrueDecisionNode, _ = NewDecisionNode("alwaysTrueDecisionNode", func(*Context) (bool, error) { return true, nil })
)

func Test_NodeSystem_IsValid(t *testing.T) {
	testCases := []struct {
		name                string
		givenNodes          []Node
		givenNodesJoinModes map[Node]JoinMode
		givenLinks          []nodeLink
		expectedNodeSystem  *NodeSystem
		expectedErrors      []error
	}{
		{
			name: "Can have no nodes",
			expectedNodeSystem: &NodeSystem{
				nodes:          []Node{},
				nodesJoinModes: map[Node]JoinMode{},
				links:          []nodeLink{},
			},
		},
		{
			name: "Can have one action node",
			givenNodes: []Node{
				someActionNode,
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []Node{
					someActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{},
				links:          []nodeLink{},
			},
		},
		{
			name: "Can't have the same node more than one time",
			givenNodes: []Node{
				someActionNode,
				someActionNode,
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []Node{
					someActionNode,
					someActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{},
				links:          []nodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have multiple instances (2) of the same node: %+v", someActionNode),
			},
		},
		{
			name: "Can't have decision node without link to another node",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []Node{
					alwaysTrueDecisionNode,
				},
				nodesJoinModes: map[Node]JoinMode{},
				links:          []nodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have decision node without link from it: %+v", alwaysTrueDecisionNode),
			},
		},
		{
			name: "Can have a link between decision node and action node",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{},
				links: []nodeLink{
					newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
				},
			},
		},
		{
			name: "Can have 2 action nodes link together",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLink(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []Node{
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{},
				links: []nodeLink{
					newNodeLink(someActionNode, anotherActionNode),
				},
			},
		},
		{
			name: "Can't add empty 'from' on branch link",
			givenLinks: []nodeLink{
				{To: someActionNode},
			},
			expectedNodeSystem: &NodeSystem{
				nodes:          []Node{},
				nodesJoinModes: map[Node]JoinMode{},
				links:          []nodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have missing 'from' attribute"),
			},
		},
		{
			name: "Can't add empty 'to' on branch link",
			givenLinks: []nodeLink{
				{From: someActionNode},
			},
			expectedNodeSystem: &NodeSystem{
				nodes:          []Node{},
				nodesJoinModes: map[Node]JoinMode{},
				links:          []nodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have missing 'to' attribute"),
			},
		},
		{
			name: "Can't add link with the node on 'from' and 'to'",
			givenLinks: []nodeLink{
				newNodeLink(someActionNode, someActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes:          []Node{},
				nodesJoinModes: map[Node]JoinMode{},
				links:          []nodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have link on from and to the same node"),
			},
		},
		{
			name: "Can have fork links with the same 'from'",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
				newNodeLinkOnBranch(alwaysTrueDecisionNode, anotherActionNode, true),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{},
				links: []nodeLink{
					newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
					newNodeLinkOnBranch(alwaysTrueDecisionNode, anotherActionNode, true),
				},
			},
		},
		{
			name: "Can have join links with the same 'to'",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
				anotherActionNode,
			},
			givenNodesJoinModes: map[Node]JoinMode{
				anotherActionNode: JoinAnd,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(alwaysTrueDecisionNode, anotherActionNode, true),
				newNodeLink(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{
					anotherActionNode: JoinAnd,
				},
				links: []nodeLink{
					newNodeLinkOnBranch(alwaysTrueDecisionNode, anotherActionNode, true),
					newNodeLink(someActionNode, anotherActionNode),
				},
			},
		},
		{
			name: "Can have merge links with the same 'to'",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
				anotherActionNode,
			},
			givenNodesJoinModes: map[Node]JoinMode{
				anotherActionNode: JoinOr,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(alwaysTrueDecisionNode, anotherActionNode, true),
				newNodeLink(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{
					anotherActionNode: JoinOr,
				},
				links: []nodeLink{
					newNodeLinkOnBranch(alwaysTrueDecisionNode, anotherActionNode, true),
					newNodeLink(someActionNode, anotherActionNode),
				},
			},
		},
		{
			name: "Can't hava a link with branch who is not needed",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(someActionNode, anotherActionNode, true),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []Node{
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{},
				links:          []nodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have not needed branch"),
			},
		},
		{
			name: "Can't have a link with a missing branch",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLink(alwaysTrueDecisionNode, someActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{},
				links:          []nodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have missing branch"),
				fmt.Errorf("can't have decision node without link from it: %+v", alwaysTrueDecisionNode),
			},
		},
		{
			name: "Can't have a link with an undeclared node as 'to'",
			givenNodes: []Node{
				someActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLink(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []Node{
					someActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{},
				links: []nodeLink{
					newNodeLink(someActionNode, anotherActionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have undeclared node '%+v' as 'to' in branch link %+v", anotherActionNode, newNodeLink(someActionNode, anotherActionNode)),
			},
		},
		{
			name: "Can't have a link with an undeclared node as 'from'",
			givenNodes: []Node{
				anotherActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLink(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []Node{
					anotherActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{},
				links: []nodeLink{
					newNodeLink(someActionNode, anotherActionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have undeclared node '%+v' as 'from' in branch link %+v", someActionNode, newNodeLink(someActionNode, anotherActionNode)),
			},
		},
		{
			name: "Can't have cycle between 2 links",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLink(someActionNode, anotherActionNode),
				newNodeLink(anotherActionNode, someActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []Node{
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{},
				links: []nodeLink{
					newNodeLink(someActionNode, anotherActionNode),
					newNodeLink(anotherActionNode, someActionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("Can't have cycle in links between nodes: %+v", []nodeLink{
					newNodeLink(someActionNode, anotherActionNode),
					newNodeLink(anotherActionNode, someActionNode),
				}),
			},
		},
		{
			name: "Can't have cycle between some links",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
				newNodeLink(someActionNode, anotherActionNode),
				newNodeLink(anotherActionNode, someActionNode),
			},
			givenNodesJoinModes: map[Node]JoinMode{
				someActionNode: JoinAnd,
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{
					someActionNode: JoinAnd,
				},
				links: []nodeLink{
					newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
					newNodeLink(someActionNode, anotherActionNode),
					newNodeLink(anotherActionNode, someActionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("Can't have cycle in links between nodes: %+v", []nodeLink{
					newNodeLink(someActionNode, anotherActionNode),
					newNodeLink(anotherActionNode, someActionNode),
				}),
			},
		},
		{
			name: "Can't have cycle between some links with branch",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
				newNodeLinkOnBranch(alwaysTrueDecisionNode, anotherActionNode, false),
				newNodeLink(anotherActionNode, alwaysTrueDecisionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{},
				links: []nodeLink{
					newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
					newNodeLinkOnBranch(alwaysTrueDecisionNode, anotherActionNode, false),
					newNodeLink(anotherActionNode, alwaysTrueDecisionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("Can't have cycle in links between nodes: %+v", []nodeLink{
					newNodeLinkOnBranch(alwaysTrueDecisionNode, anotherActionNode, false),
					newNodeLink(anotherActionNode, alwaysTrueDecisionNode),
				}),
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			system := NewNodeSystem()
			errs := loadNodeSystem(system, testCase.givenNodes, testCase.givenNodesJoinModes, testCase.givenLinks)

			_, validityErrs := system.IsValid()
			errs = append(errs, validityErrs...)

			if !cmp.Equal(errs, testCase.expectedErrors, utils.ErrorComparator) {
				t.Errorf("errors - got: %+v, want: %+v", errs, testCase.expectedErrors)
			}
			if !cmp.Equal(system, testCase.expectedNodeSystem) {
				t.Errorf("system - got: %+v, want: %+v", system, testCase.expectedNodeSystem)
			}
		})
	}
}

func Test_NodeSystem_Activate(t *testing.T) {
	testCases := []struct {
		name                               string
		givenNodes                         []Node
		givenNodesJoinModes                map[Node]JoinMode
		givenLinks                         []nodeLink
		givenNodesAfterActivation          []Node
		givenNodesJoinModesAfterActivation map[Node]JoinMode
		givenLinksAfterActivation          []nodeLink
		expectedActivatation               bool
		expectedInitialNodes               []Node
		expectedFollowingNodesTree         map[Node]map[*bool][]Node
		expectedAncestorsNodesTree         map[Node]map[*bool][]Node
		expectedErrors                     []error
	}{
		{
			name:                       "Can activate an empty validated system",
			expectedActivatation:       true,
			expectedInitialNodes:       []Node{},
			expectedFollowingNodesTree: map[Node]map[*bool][]Node{},
			expectedAncestorsNodesTree: map[Node]map[*bool][]Node{},
		},
		{
			name: "Can't activate an unvalidated system",
			givenLinks: []nodeLink{
				newNodeLink(someActionNode, anotherActionNode),
			},
			expectedActivatation: false,
			expectedErrors: []error{
				errors.New("can't activate a unvalidated node system"),
			},
		},
		{
			name:                       "Can activate an one node validated system",
			givenNodes:                 []Node{someActionNode},
			expectedActivatation:       true,
			expectedInitialNodes:       []Node{someActionNode},
			expectedFollowingNodesTree: map[Node]map[*bool][]Node{},
		},
		{
			name: "Can activate an no needed branch node link validated system",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLink(someActionNode, anotherActionNode),
			},
			expectedActivatation: true,
			expectedInitialNodes: []Node{
				someActionNode,
			},
			expectedFollowingNodesTree: map[Node]map[*bool][]Node{
				someActionNode: {
					nil: {anotherActionNode},
				},
			},
		},
		{
			name: "Can activate an orphan action node validated system",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
				alwaysTrueDecisionNode,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(alwaysTrueDecisionNode, anotherActionNode, true),
			},
			expectedActivatation: true,
			expectedInitialNodes: []Node{
				someActionNode,
				alwaysTrueDecisionNode,
			},
			expectedFollowingNodesTree: map[Node]map[*bool][]Node{
				alwaysTrueDecisionNode: {
					utils.BoolPointer(true): {anotherActionNode},
				},
			},
		},
		{
			name: "Can't add any node after activation",
			givenNodesAfterActivation: []Node{
				someActionNode,
			},
			expectedActivatation: true,
			expectedErrors: []error{
				fmt.Errorf("can't add node, node system is freeze due to activation"),
			},
		},
		{
			name: "Can't add any node join mode after activation",
			givenNodesJoinModesAfterActivation: map[Node]JoinMode{
				someActionNode: JoinNone,
			},
			expectedActivatation: true,
			expectedErrors: []error{
				fmt.Errorf("can't add node join mode, node system is freeze due to activation"),
			},
		},
		{
			name: "Can't add any branch link after activation",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinksAfterActivation: []nodeLink{
				newNodeLink(someActionNode, anotherActionNode),
			},
			expectedActivatation: true,
			expectedErrors: []error{
				fmt.Errorf("can't add branch link, node system is freeze due to activation"),
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			system := NewNodeSystem()
			errs := loadNodeSystem(system, testCase.givenNodes, testCase.givenNodesJoinModes, testCase.givenLinks)

			err := system.Activate()
			if err != nil {
				errs = append(errs, err)
			}

			afterErrs := loadNodeSystem(system, testCase.givenNodesAfterActivation, testCase.givenNodesJoinModesAfterActivation, testCase.givenLinksAfterActivation)
			if afterErrs != nil {
				errs = append(errs, afterErrs...)
			}

			if testCase.givenNodesAfterActivation != nil || testCase.givenNodesJoinModesAfterActivation != nil || testCase.givenLinksAfterActivation != nil {
				_, validityErrs := system.IsValid()
				errs = append(errs, validityErrs...)
			}

			if system.IsActivated() != testCase.expectedActivatation {
				t.Errorf("activation - got: %+v, want: %+v", system.activated, testCase.expectedActivatation)
			}
			if !cmp.Equal(errs, testCase.expectedErrors, utils.ErrorComparator) {
				t.Errorf("errors - got: %+v, want: %+v", errs, testCase.expectedErrors)
			}
			if testCase.expectedInitialNodes != nil && !cmp.Equal(system.InitialNodes(), testCase.expectedInitialNodes, NodeComparator) {
				t.Errorf("initial nodes - got: %+v, want: %+v", system.InitialNodes(), testCase.expectedInitialNodes)
			}
			if testCase.expectedFollowingNodesTree != nil && !cmp.Equal(system.followingNodesTree, testCase.expectedFollowingNodesTree, NodeComparator) {
				t.Errorf("following node tree - got: %#v, want: %#v", system.followingNodesTree, testCase.expectedFollowingNodesTree)
			}
			if testCase.expectedAncestorsNodesTree != nil && !cmp.Equal(system.ancestorsNodesTree, testCase.expectedAncestorsNodesTree, NodeComparator) {
				t.Errorf("ancestors node tree - got: %#v, want: %#v", system.ancestorsNodesTree, testCase.expectedAncestorsNodesTree)
			}
		})
	}
}

func Test_NodeSystem_multiple_activatation(t *testing.T) {
	system := NewNodeSystem()
	system.Activate()
	err := system.Activate()
	if err != nil {
		t.Errorf("Can be activate multiple times without errors, but got: %+v", err)
	}
}

func Test_NodeSystem_Follow(t *testing.T) {
	testCases := []struct {
		name                   string
		givenNodes             []Node
		givenLinks             []nodeLink
		givenNode              Node
		givenBranch            *bool
		expectedFollowingNodes []Node
		expectedError          error
	}{
		{
			name: "Can't follow on an unactivated system",
			givenLinks: []nodeLink{
				newNodeLink(someActionNode, anotherActionNode),
			},
			expectedError: errors.New("can't follow a node if system is not activated"),
		},
		{
			name: "Can follow 'from' on link",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLink(someActionNode, anotherActionNode),
			},
			givenNode:              someActionNode,
			expectedFollowingNodes: []Node{anotherActionNode},
		},
		{
			name: "Can follow 'to' on link",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLink(someActionNode, anotherActionNode),
			},
			givenNode:              anotherActionNode,
			expectedFollowingNodes: nil,
		},
		{
			name: "Can follow 'from' on branch link",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenNode:              alwaysTrueDecisionNode,
			givenBranch:            utils.BoolPointer(true),
			expectedFollowingNodes: []Node{someActionNode},
		},
		{
			name: "Can follow 'to' on branch link",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenNode:              someActionNode,
			expectedFollowingNodes: nil,
		},
		{
			name: "Can't follow 'from' on branch link but without passing the branch",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenNode:              alwaysTrueDecisionNode,
			expectedFollowingNodes: nil,
		},
		{
			name: "Can't follow 'from' on branch link but without passing the right branch",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenNode:              alwaysTrueDecisionNode,
			givenBranch:            utils.BoolPointer(false),
			expectedFollowingNodes: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			system := NewNodeSystem()
			for _, node := range testCase.givenNodes {
				system.AddNode(node)
			}
			for _, link := range testCase.givenLinks {
				if link.Branch == nil {
					system.AddLink(link.From, link.To)
				} else {
					system.AddLinkOnBranch(link.From, link.To, *link.Branch)
				}
			}
			system.Activate()
			nodes, err := system.Follow(testCase.givenNode, testCase.givenBranch)

			if !cmp.Equal(err, testCase.expectedError, utils.ErrorComparator) {
				t.Errorf("error - got: %+v, want: %+v", err, testCase.expectedError)
			}
			if !cmp.Equal(nodes, testCase.expectedFollowingNodes, NodeComparator) {
				t.Errorf("following nodes - got: %+v, want: %+v", nodes, testCase.expectedFollowingNodes)
			}
		})
	}
}

func Test_NodeSystem_Ancestors(t *testing.T) {
	testCases := []struct {
		name                  string
		givenNodes            []Node
		givenLinks            []nodeLink
		givenNode             Node
		givenBranch           *bool
		expectedAncestorNodes []Node
		expectedError         error
	}{
		{
			name: "Can't get ancestors on an unactivated system",
			givenLinks: []nodeLink{
				newNodeLink(someActionNode, anotherActionNode),
			},
			expectedError: errors.New("can't get ancestors of a node if system is not activated"),
		},
		{
			name: "Can have ancestors of 'from' on link",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLink(someActionNode, anotherActionNode),
			},
			givenNode:             someActionNode,
			expectedAncestorNodes: nil,
		},
		{
			name: "Can have ancestors of 'to' on link",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLink(someActionNode, anotherActionNode),
			},
			givenNode:             anotherActionNode,
			expectedAncestorNodes: []Node{someActionNode},
		},
		{
			name: "Can have ancestors of 'from' on branch link",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenNode:             alwaysTrueDecisionNode,
			expectedAncestorNodes: nil,
		},
		{
			name: "Can have ancestors of 'to' on branch link",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenNode:             someActionNode,
			givenBranch:           utils.BoolPointer(true),
			expectedAncestorNodes: []Node{alwaysTrueDecisionNode},
		},
		{
			name: "Can have ancestors 'from' on branch link but without passing the branch",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenNode:             someActionNode,
			expectedAncestorNodes: nil,
		},
		{
			name: "Can't have ancestors 'from' on branch link but without passing the right branch",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenNode:             someActionNode,
			givenBranch:           utils.BoolPointer(false),
			expectedAncestorNodes: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			system := NewNodeSystem()
			loadNodeSystem(system, testCase.givenNodes, nil, testCase.givenLinks)

			system.Activate()
			nodes, err := system.Ancestors(testCase.givenNode, testCase.givenBranch)

			if !cmp.Equal(err, testCase.expectedError, utils.ErrorComparator) {
				t.Errorf("error - got: %+v, want: %+v", err, testCase.expectedError)
			}
			if !cmp.Equal(nodes, testCase.expectedAncestorNodes, NodeComparator) {
				t.Errorf("ancestor nodes - got: %+v, want: %+v", nodes, testCase.expectedAncestorNodes)
			}
		})
	}
}

func Test_JoinModeOfNode_found(t *testing.T) {
	givenNode := someActionNode
	givenJoinMode := JoinAnd

	system := NewNodeSystem()
	system.AddNode(givenNode)
	system.ConfigureJoinModeOnNode(givenNode, givenJoinMode)
	system.Activate()

	storedJoinMode := system.JoinModeOfNode(givenNode)

	if givenJoinMode != storedJoinMode {
		t.Errorf("got: %+v, want: %+v", storedJoinMode, givenJoinMode)
	}
}

func Test_JoinModeOfNode_notfound(t *testing.T) {
	givenNode := someActionNode

	system := NewNodeSystem()
	system.AddNode(givenNode)
	system.Activate()

	storedJoinMode := system.JoinModeOfNode(givenNode)

	if JoinNone != storedJoinMode {
		t.Errorf("got: %+v, want: %+v", storedJoinMode, JoinNone)
	}
}

func Test_Github_Issue_10(t *testing.T) {
	action1, _ := NewActionNode("action1", func(c *Context) error {
		return nil
	})
	decision2, _ := NewDecisionNode("decision2", func(c *Context) (bool, error) {
		return true, nil
	})
	decision3, _ := NewDecisionNode("decision3", func(c *Context) (bool, error) {
		return true, nil
	})
	action4, _ := NewActionNode("action4", func(c *Context) error {
		return nil
	})

	ns := NewNodeSystem()
	ns.AddNode(action1)
	ns.AddNode(decision2)
	ns.AddNode(decision3)
	ns.AddNode(action4)
	ns.AddLink(action1, decision2)
	ns.AddLink(action1, decision3)
	ns.AddLinkOnBranch(decision2, action4, false)
	ns.AddLinkOnBranch(decision3, action4, true)
	ns.ConfigureJoinModeOnNode(action4, JoinNone)
	_, errs := ns.IsValid()

	expectedErrors := []error{
		errors.New("can't have multiple links (2) to the same node: action4 without join mode"),
	}

	if !cmp.Equal(errs, expectedErrors, utils.ErrorComparator) {
		t.Errorf("errors - got: %+v, want: %+v", errs, expectedErrors)
	}
}

func Test_Github_Issue_16(t *testing.T) {
	trigger, _ := NewDecisionNode("trigger", func(c *Context) (bool, error) {
		return true, nil
	})
	a1, _ := NewActionNode("a1", func(c *Context) error {
		return nil
	})
	a2, _ := NewActionNode("a2", func(c *Context) error {
		return nil
	})
	a3, _ := NewActionNode("a3", func(c *Context) error {
		return nil
	})

	ns := NewNodeSystem()
	ns.AddNode(trigger)
	ns.AddNode(a1)
	ns.AddNode(a2)
	ns.AddNode(a3)
	ns.AddLinkOnBranch(trigger, a1, true)
	ns.AddLink(a1, a2)
	ns.AddLink(a2, a3)
	ns.AddLink(a3, a2)
	ns.ConfigureJoinModeOnNode(a2, JoinAnd)
	_, errs := ns.IsValid()

	expectedErrors := []error{
		fmt.Errorf("Can't have cycle in links between nodes: %+v", []nodeLink{
			newNodeLink(a2, a3),
			newNodeLink(a3, a2),
		}),
	}

	if !cmp.Equal(errs, expectedErrors, utils.ErrorComparator) {
		t.Errorf("errors - got: %+v, want: %+v", errs, expectedErrors)
	}
}

func Test_multiple_cycles(t *testing.T) {
	trigger, _ := NewDecisionNode("trigger", func(c *Context) (bool, error) {
		return true, nil
	})
	a1, _ := NewActionNode("a1", func(c *Context) error {
		return nil
	})
	a2, _ := NewActionNode("a2", func(c *Context) error {
		return nil
	})
	a3, _ := NewActionNode("a3", func(c *Context) error {
		return nil
	})
	a4, _ := NewActionNode("a4", func(c *Context) error {
		return nil
	})
	a5, _ := NewActionNode("a5", func(c *Context) error {
		return nil
	})
	a6, _ := NewActionNode("a6", func(c *Context) error {
		return nil
	})
	a7, _ := NewActionNode("a7", func(c *Context) error {
		return nil
	})

	ns := NewNodeSystem()
	ns.AddNode(trigger)
	ns.AddNode(a1)
	ns.AddNode(a2)
	ns.AddNode(a3)
	ns.AddNode(a4)
	ns.AddNode(a5)
	ns.AddNode(a6)
	ns.AddNode(a7)

	// Cycle between a2 and a3
	ns.AddLinkOnBranch(trigger, a1, true)
	ns.AddLink(a1, a2)
	ns.AddLink(a2, a3)
	ns.AddLink(a3, a2)
	ns.ConfigureJoinModeOnNode(a2, JoinAnd)

	// Cycle between a5, a6, and a7
	ns.AddLinkOnBranch(trigger, a4, false)
	ns.AddLink(a4, a5)
	ns.AddLink(a5, a6)
	ns.AddLink(a6, a7)
	ns.AddLink(a7, a5)
	ns.ConfigureJoinModeOnNode(a5, JoinAnd)
	_, errs := ns.IsValid()

	expectedErrors := []error{
		fmt.Errorf("Can't have cycle in links between nodes: %+v", []nodeLink{
			newNodeLink(a2, a3),
			newNodeLink(a3, a2),
		}),
		fmt.Errorf("Can't have cycle in links between nodes: %+v", []nodeLink{
			newNodeLink(a5, a6),
			newNodeLink(a6, a7),
			newNodeLink(a7, a5),
		}),
	}

	if !cmp.Equal(errs, expectedErrors, utils.ErrorComparator) {
		t.Errorf("errors - got: %+v, want: %+v", errs, expectedErrors)
	}
}

func loadNodeSystem(system *NodeSystem, nodes []Node, nodesJoinModes map[Node]JoinMode, links []nodeLink) []error {
	var errs []error
	for _, node := range nodes {
		_, err := system.AddNode(node)
		if err != nil {
			errs = append(errs, err)
		}
	}
	for node, mode := range nodesJoinModes {
		_, err := system.ConfigureJoinModeOnNode(node, mode)
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, link := range links {
		if link.Branch == nil {
			_, err := system.AddLink(link.From, link.To)
			if err != nil {
				errs = append(errs, err)
			}
		} else {
			_, err := system.AddLinkOnBranch(link.From, link.To, *link.Branch)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}
