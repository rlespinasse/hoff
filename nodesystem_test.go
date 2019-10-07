package hoff

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/internal/nodelink"
	"github.com/rlespinasse/hoff/internal/utils"
	"github.com/rlespinasse/hoff/node"
)

var (
	someActionNode, _         = node.NewAction("someActionNode", func(*node.Context) error { return nil })
	anotherActionNode, _      = node.NewAction("anotherActionNode", func(*node.Context) error { return nil })
	alwaysTrueDecisionNode, _ = node.NewDecision("alwaysTrueDecisionNode", func(*node.Context) (bool, error) { return true, nil })
)

func Test_NodeSystem_IsValid(t *testing.T) {
	testCases := []struct {
		name                string
		givenNodes          []node.Node
		givenNodesJoinModes map[node.Node]JoinMode
		givenLinks          []nodelink.NodeLink
		expectedNodeSystem  *NodeSystem
		expectedErrors      []error
	}{
		{
			name: "Can have no nodes",
			expectedNodeSystem: &NodeSystem{
				nodes:          []node.Node{},
				nodesJoinModes: map[node.Node]JoinMode{},
				links:          []nodelink.NodeLink{},
			},
		},
		{
			name: "Can have one action node",
			givenNodes: []node.Node{
				someActionNode,
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []node.Node{
					someActionNode,
				},
				nodesJoinModes: map[node.Node]JoinMode{},
				links:          []nodelink.NodeLink{},
			},
		},
		{
			name: "Can't have the same node more than one time",
			givenNodes: []node.Node{
				someActionNode,
				someActionNode,
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []node.Node{
					someActionNode,
					someActionNode,
				},
				nodesJoinModes: map[node.Node]JoinMode{},
				links:          []nodelink.NodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have multiple instances (2) of the same node: %+v", someActionNode),
			},
		},
		{
			name: "Can't have decision node without link to another node",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []node.Node{
					alwaysTrueDecisionNode,
				},
				nodesJoinModes: map[node.Node]JoinMode{},
				links:          []nodelink.NodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have decision node without link from it: %+v", alwaysTrueDecisionNode),
			},
		},
		{
			name: "Can have a link between decision node and action node",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []node.Node{
					alwaysTrueDecisionNode,
					someActionNode,
				},
				nodesJoinModes: map[node.Node]JoinMode{},
				links: []nodelink.NodeLink{
					nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
				},
			},
		},
		{
			name: "Can have 2 action nodes link together",
			givenNodes: []node.Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.New(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []node.Node{
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[node.Node]JoinMode{},
				links: []nodelink.NodeLink{
					nodelink.New(someActionNode, anotherActionNode),
				},
			},
		},
		{
			name: "Can't add empty 'from' on branch link",
			givenLinks: []nodelink.NodeLink{
				{To: someActionNode},
			},
			expectedNodeSystem: &NodeSystem{
				nodes:          []node.Node{},
				nodesJoinModes: map[node.Node]JoinMode{},
				links:          []nodelink.NodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have missing 'from' attribute"),
			},
		},
		{
			name: "Can't add empty 'to' on branch link",
			givenLinks: []nodelink.NodeLink{
				{From: someActionNode},
			},
			expectedNodeSystem: &NodeSystem{
				nodes:          []node.Node{},
				nodesJoinModes: map[node.Node]JoinMode{},
				links:          []nodelink.NodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have missing 'to' attribute"),
			},
		},
		{
			name: "Can't add link with the node on 'from' and 'to'",
			givenLinks: []nodelink.NodeLink{
				nodelink.New(someActionNode, someActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes:          []node.Node{},
				nodesJoinModes: map[node.Node]JoinMode{},
				links:          []nodelink.NodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have link on from and to the same node"),
			},
		},
		{
			name: "Can have fork links with the same 'from'",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
				nodelink.NewOnBranch(alwaysTrueDecisionNode, anotherActionNode, true),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []node.Node{
					alwaysTrueDecisionNode,
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[node.Node]JoinMode{},
				links: []nodelink.NodeLink{
					nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
					nodelink.NewOnBranch(alwaysTrueDecisionNode, anotherActionNode, true),
				},
			},
		},
		{
			name: "Can have join links with the same 'to'",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
				someActionNode,
				anotherActionNode,
			},
			givenNodesJoinModes: map[node.Node]JoinMode{
				anotherActionNode: JoinAnd,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(alwaysTrueDecisionNode, anotherActionNode, true),
				nodelink.New(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []node.Node{
					alwaysTrueDecisionNode,
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[node.Node]JoinMode{
					anotherActionNode: JoinAnd,
				},
				links: []nodelink.NodeLink{
					nodelink.NewOnBranch(alwaysTrueDecisionNode, anotherActionNode, true),
					nodelink.New(someActionNode, anotherActionNode),
				},
			},
		},
		{
			name: "Can have merge links with the same 'to'",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
				someActionNode,
				anotherActionNode,
			},
			givenNodesJoinModes: map[node.Node]JoinMode{
				anotherActionNode: JoinOr,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(alwaysTrueDecisionNode, anotherActionNode, true),
				nodelink.New(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []node.Node{
					alwaysTrueDecisionNode,
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[node.Node]JoinMode{
					anotherActionNode: JoinOr,
				},
				links: []nodelink.NodeLink{
					nodelink.NewOnBranch(alwaysTrueDecisionNode, anotherActionNode, true),
					nodelink.New(someActionNode, anotherActionNode),
				},
			},
		},
		{
			name: "Can't hava a link with branch who is not needed",
			givenNodes: []node.Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(someActionNode, anotherActionNode, true),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []node.Node{
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[node.Node]JoinMode{},
				links:          []nodelink.NodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have not needed branch"),
			},
		},
		{
			name: "Can't have a link with a missing branch",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.New(alwaysTrueDecisionNode, someActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []node.Node{
					alwaysTrueDecisionNode,
					someActionNode,
				},
				nodesJoinModes: map[node.Node]JoinMode{},
				links:          []nodelink.NodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have missing branch"),
				fmt.Errorf("can't have decision node without link from it: %+v", alwaysTrueDecisionNode),
			},
		},
		{
			name: "Can't have a link with an undeclared node as 'to'",
			givenNodes: []node.Node{
				someActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.New(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []node.Node{
					someActionNode,
				},
				nodesJoinModes: map[node.Node]JoinMode{},
				links: []nodelink.NodeLink{
					nodelink.New(someActionNode, anotherActionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have undeclared node '%+v' as 'to' in branch link %+v", anotherActionNode, nodelink.New(someActionNode, anotherActionNode)),
			},
		},
		{
			name: "Can't have a link with an undeclared node as 'from'",
			givenNodes: []node.Node{
				anotherActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.New(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []node.Node{
					anotherActionNode,
				},
				nodesJoinModes: map[node.Node]JoinMode{},
				links: []nodelink.NodeLink{
					nodelink.New(someActionNode, anotherActionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have undeclared node '%+v' as 'from' in branch link %+v", someActionNode, nodelink.New(someActionNode, anotherActionNode)),
			},
		},
		{
			name: "Can't have cycle between 2 links",
			givenNodes: []node.Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.New(someActionNode, anotherActionNode),
				nodelink.New(anotherActionNode, someActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []node.Node{
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[node.Node]JoinMode{},
				links: []nodelink.NodeLink{
					nodelink.New(someActionNode, anotherActionNode),
					nodelink.New(anotherActionNode, someActionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("Can't have cycle in links between nodes: %+v", []nodelink.NodeLink{
					nodelink.New(someActionNode, anotherActionNode),
					nodelink.New(anotherActionNode, someActionNode),
				}),
			},
		},
		{
			name: "Can't have cycle between some links",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
				nodelink.New(someActionNode, anotherActionNode),
				nodelink.New(anotherActionNode, someActionNode),
			},
			givenNodesJoinModes: map[node.Node]JoinMode{
				someActionNode: JoinAnd,
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []node.Node{
					alwaysTrueDecisionNode,
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[node.Node]JoinMode{
					someActionNode: JoinAnd,
				},
				links: []nodelink.NodeLink{
					nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
					nodelink.New(someActionNode, anotherActionNode),
					nodelink.New(anotherActionNode, someActionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("Can't have cycle in links between nodes: %+v", []nodelink.NodeLink{
					nodelink.New(someActionNode, anotherActionNode),
					nodelink.New(anotherActionNode, someActionNode),
				}),
			},
		},
		{
			name: "Can't have cycle between some links with branch",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
				nodelink.NewOnBranch(alwaysTrueDecisionNode, anotherActionNode, false),
				nodelink.New(anotherActionNode, alwaysTrueDecisionNode),
			},
			expectedNodeSystem: &NodeSystem{
				nodes: []node.Node{
					alwaysTrueDecisionNode,
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[node.Node]JoinMode{},
				links: []nodelink.NodeLink{
					nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
					nodelink.NewOnBranch(alwaysTrueDecisionNode, anotherActionNode, false),
					nodelink.New(anotherActionNode, alwaysTrueDecisionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("Can't have cycle in links between nodes: %+v", []nodelink.NodeLink{
					nodelink.NewOnBranch(alwaysTrueDecisionNode, anotherActionNode, false),
					nodelink.New(anotherActionNode, alwaysTrueDecisionNode),
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
		givenNodes                         []node.Node
		givenNodesJoinModes                map[node.Node]JoinMode
		givenLinks                         []nodelink.NodeLink
		givenNodesAfterActivation          []node.Node
		givenNodesJoinModesAfterActivation map[node.Node]JoinMode
		givenLinksAfterActivation          []nodelink.NodeLink
		expectedActivatation               bool
		expectedInitialNodes               []node.Node
		expectedFollowingNodesTree         map[node.Node]map[*bool][]node.Node
		expectedAncestorsNodesTree         map[node.Node]map[*bool][]node.Node
		expectedErrors                     []error
	}{
		{
			name:                       "Can activate an empty validated system",
			expectedActivatation:       true,
			expectedInitialNodes:       []node.Node{},
			expectedFollowingNodesTree: map[node.Node]map[*bool][]node.Node{},
			expectedAncestorsNodesTree: map[node.Node]map[*bool][]node.Node{},
		},
		{
			name: "Can't activate an unvalidated system",
			givenLinks: []nodelink.NodeLink{
				nodelink.New(someActionNode, anotherActionNode),
			},
			expectedActivatation: false,
			expectedErrors: []error{
				errors.New("can't activate a unvalidated node system"),
			},
		},
		{
			name:                       "Can activate an one node validated system",
			givenNodes:                 []node.Node{someActionNode},
			expectedActivatation:       true,
			expectedInitialNodes:       []node.Node{someActionNode},
			expectedFollowingNodesTree: map[node.Node]map[*bool][]node.Node{},
		},
		{
			name: "Can activate an no needed branch node link validated system",
			givenNodes: []node.Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.New(someActionNode, anotherActionNode),
			},
			expectedActivatation: true,
			expectedInitialNodes: []node.Node{
				someActionNode,
			},
			expectedFollowingNodesTree: map[node.Node]map[*bool][]node.Node{
				someActionNode: {
					nil: {anotherActionNode},
				},
			},
		},
		{
			name: "Can activate an orphan action node validated system",
			givenNodes: []node.Node{
				someActionNode,
				anotherActionNode,
				alwaysTrueDecisionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(alwaysTrueDecisionNode, anotherActionNode, true),
			},
			expectedActivatation: true,
			expectedInitialNodes: []node.Node{
				someActionNode,
				alwaysTrueDecisionNode,
			},
			expectedFollowingNodesTree: map[node.Node]map[*bool][]node.Node{
				alwaysTrueDecisionNode: {
					utils.BoolPointer(true): {anotherActionNode},
				},
			},
		},
		{
			name: "Can't add any node after activation",
			givenNodesAfterActivation: []node.Node{
				someActionNode,
			},
			expectedActivatation: true,
			expectedErrors: []error{
				fmt.Errorf("can't add node, node system is freeze due to activation"),
			},
		},
		{
			name: "Can't add any node join mode after activation",
			givenNodesJoinModesAfterActivation: map[node.Node]JoinMode{
				someActionNode: JoinNone,
			},
			expectedActivatation: true,
			expectedErrors: []error{
				fmt.Errorf("can't add node join mode, node system is freeze due to activation"),
			},
		},
		{
			name: "Can't add any branch link after activation",
			givenNodes: []node.Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinksAfterActivation: []nodelink.NodeLink{
				nodelink.New(someActionNode, anotherActionNode),
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
			if testCase.expectedInitialNodes != nil && !cmp.Equal(system.InitialNodes(), testCase.expectedInitialNodes, node.NodeComparator) {
				t.Errorf("initial nodes - got: %+v, want: %+v", system.InitialNodes(), testCase.expectedInitialNodes)
			}
			if testCase.expectedFollowingNodesTree != nil && !cmp.Equal(system.followingNodesTree, testCase.expectedFollowingNodesTree, node.NodeComparator) {
				t.Errorf("following node tree - got: %#v, want: %#v", system.followingNodesTree, testCase.expectedFollowingNodesTree)
			}
			if testCase.expectedAncestorsNodesTree != nil && !cmp.Equal(system.ancestorsNodesTree, testCase.expectedAncestorsNodesTree, node.NodeComparator) {
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
		givenNodes             []node.Node
		givenLinks             []nodelink.NodeLink
		givenNode              node.Node
		givenBranch            *bool
		expectedFollowingNodes []node.Node
		expectedError          error
	}{
		{
			name: "Can't follow on an unactivated system",
			givenLinks: []nodelink.NodeLink{
				nodelink.New(someActionNode, anotherActionNode),
			},
			expectedError: errors.New("can't follow a node if system is not activated"),
		},
		{
			name: "Can follow 'from' on link",
			givenNodes: []node.Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.New(someActionNode, anotherActionNode),
			},
			givenNode:              someActionNode,
			expectedFollowingNodes: []node.Node{anotherActionNode},
		},
		{
			name: "Can follow 'to' on link",
			givenNodes: []node.Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.New(someActionNode, anotherActionNode),
			},
			givenNode:              anotherActionNode,
			expectedFollowingNodes: nil,
		},
		{
			name: "Can follow 'from' on branch link",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenNode:              alwaysTrueDecisionNode,
			givenBranch:            utils.BoolPointer(true),
			expectedFollowingNodes: []node.Node{someActionNode},
		},
		{
			name: "Can follow 'to' on branch link",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenNode:              someActionNode,
			expectedFollowingNodes: nil,
		},
		{
			name: "Can't follow 'from' on branch link but without passing the branch",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenNode:              alwaysTrueDecisionNode,
			expectedFollowingNodes: nil,
		},
		{
			name: "Can't follow 'from' on branch link but without passing the right branch",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
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
			if !cmp.Equal(nodes, testCase.expectedFollowingNodes, node.NodeComparator) {
				t.Errorf("following nodes - got: %+v, want: %+v", nodes, testCase.expectedFollowingNodes)
			}
		})
	}
}

func Test_NodeSystem_Ancestors(t *testing.T) {
	testCases := []struct {
		name                  string
		givenNodes            []node.Node
		givenLinks            []nodelink.NodeLink
		givenNode             node.Node
		givenBranch           *bool
		expectedAncestorNodes []node.Node
		expectedError         error
	}{
		{
			name: "Can't get ancestors on an unactivated system",
			givenLinks: []nodelink.NodeLink{
				nodelink.New(someActionNode, anotherActionNode),
			},
			expectedError: errors.New("can't get ancestors of a node if system is not activated"),
		},
		{
			name: "Can have ancestors of 'from' on link",
			givenNodes: []node.Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.New(someActionNode, anotherActionNode),
			},
			givenNode:             someActionNode,
			expectedAncestorNodes: nil,
		},
		{
			name: "Can have ancestors of 'to' on link",
			givenNodes: []node.Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.New(someActionNode, anotherActionNode),
			},
			givenNode:             anotherActionNode,
			expectedAncestorNodes: []node.Node{someActionNode},
		},
		{
			name: "Can have ancestors of 'from' on branch link",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenNode:             alwaysTrueDecisionNode,
			expectedAncestorNodes: nil,
		},
		{
			name: "Can have ancestors of 'to' on branch link",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenNode:             someActionNode,
			givenBranch:           utils.BoolPointer(true),
			expectedAncestorNodes: []node.Node{alwaysTrueDecisionNode},
		},
		{
			name: "Can have ancestors 'from' on branch link but without passing the branch",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenNode:             someActionNode,
			expectedAncestorNodes: nil,
		},
		{
			name: "Can't have ancestors 'from' on branch link but without passing the right branch",
			givenNodes: []node.Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
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
			if !cmp.Equal(nodes, testCase.expectedAncestorNodes, node.NodeComparator) {
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
	action1, _ := node.NewAction("action1", func(c *node.Context) error {
		return nil
	})
	decision2, _ := node.NewDecision("decision2", func(c *node.Context) (bool, error) {
		return true, nil
	})
	decision3, _ := node.NewDecision("decision3", func(c *node.Context) (bool, error) {
		return true, nil
	})
	action4, _ := node.NewAction("action4", func(c *node.Context) error {
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
	trigger, _ := node.NewDecision("trigger", func(c *node.Context) (bool, error) {
		return true, nil
	})
	a1, _ := node.NewAction("a1", func(c *node.Context) error {
		return nil
	})
	a2, _ := node.NewAction("a2", func(c *node.Context) error {
		return nil
	})
	a3, _ := node.NewAction("a3", func(c *node.Context) error {
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
		fmt.Errorf("Can't have cycle in links between nodes: %+v", []nodelink.NodeLink{
			nodelink.New(a2, a3),
			nodelink.New(a3, a2),
		}),
	}

	if !cmp.Equal(errs, expectedErrors, utils.ErrorComparator) {
		t.Errorf("errors - got: %+v, want: %+v", errs, expectedErrors)
	}
}

func Test_multiple_cycles(t *testing.T) {
	trigger, _ := node.NewDecision("trigger", func(c *node.Context) (bool, error) {
		return true, nil
	})
	a1, _ := node.NewAction("a1", func(c *node.Context) error {
		return nil
	})
	a2, _ := node.NewAction("a2", func(c *node.Context) error {
		return nil
	})
	a3, _ := node.NewAction("a3", func(c *node.Context) error {
		return nil
	})
	a4, _ := node.NewAction("a4", func(c *node.Context) error {
		return nil
	})
	a5, _ := node.NewAction("a5", func(c *node.Context) error {
		return nil
	})
	a6, _ := node.NewAction("a6", func(c *node.Context) error {
		return nil
	})
	a7, _ := node.NewAction("a7", func(c *node.Context) error {
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
		fmt.Errorf("Can't have cycle in links between nodes: %+v", []nodelink.NodeLink{
			nodelink.New(a2, a3),
			nodelink.New(a3, a2),
		}),
		fmt.Errorf("Can't have cycle in links between nodes: %+v", []nodelink.NodeLink{
			nodelink.New(a5, a6),
			nodelink.New(a6, a7),
			nodelink.New(a7, a5),
		}),
	}

	if !cmp.Equal(errs, expectedErrors, utils.ErrorComparator) {
		t.Errorf("errors - got: %+v, want: %+v", errs, expectedErrors)
	}
}

func loadNodeSystem(system *NodeSystem, nodes []node.Node, nodesJoinModes map[node.Node]JoinMode, links []nodelink.NodeLink) []error {
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
