package system

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/internal/nodelink"
	"github.com/rlespinasse/hoff/internal/utils"
	"github.com/rlespinasse/hoff/node"
	"github.com/rlespinasse/hoff/system/joinmode"
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
		givenNodesJoinModes map[node.Node]joinmode.JoinMode
		givenLinks          []nodelink.NodeLink
		expectedNodeSystem  *NodeSystem
		expectedErrors      []error
	}{
		{
			name: "Can have no nodes",
			expectedNodeSystem: &NodeSystem{
				nodes:          []node.Node{},
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
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
			givenNodesJoinModes: map[node.Node]joinmode.JoinMode{
				anotherActionNode: joinmode.AND,
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{
					anotherActionNode: joinmode.AND,
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
			givenNodesJoinModes: map[node.Node]joinmode.JoinMode{
				anotherActionNode: joinmode.OR,
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{
					anotherActionNode: joinmode.OR,
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
				links: []nodelink.NodeLink{
					nodelink.New(someActionNode, anotherActionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have undeclared node '%+v' as 'from' in branch link %+v", someActionNode, nodelink.New(someActionNode, anotherActionNode)),
			},
		},
		{
			name: "Can't have cycle between some links",
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
				links: []nodelink.NodeLink{
					nodelink.New(someActionNode, anotherActionNode),
					nodelink.New(anotherActionNode, someActionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("Can't have cycle between links: %+v", []nodelink.NodeLink{
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
				nodesJoinModes: map[node.Node]joinmode.JoinMode{},
				links: []nodelink.NodeLink{
					nodelink.NewOnBranch(alwaysTrueDecisionNode, someActionNode, true),
					nodelink.NewOnBranch(alwaysTrueDecisionNode, anotherActionNode, false),
					nodelink.New(anotherActionNode, alwaysTrueDecisionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("Can't have cycle between links: %+v", []nodelink.NodeLink{
					nodelink.NewOnBranch(alwaysTrueDecisionNode, anotherActionNode, false),
					nodelink.New(anotherActionNode, alwaysTrueDecisionNode),
				}),
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			system := New()
			errs := load(system, testCase.givenNodes, testCase.givenNodesJoinModes, testCase.givenLinks)

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
		givenNodesJoinModes                map[node.Node]joinmode.JoinMode
		givenLinks                         []nodelink.NodeLink
		givenNodesAfterActivation          []node.Node
		givenNodesJoinModesAfterActivation map[node.Node]joinmode.JoinMode
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
			givenNodesJoinModesAfterActivation: map[node.Node]joinmode.JoinMode{
				someActionNode: joinmode.NONE,
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
			system := New()
			errs := load(system, testCase.givenNodes, testCase.givenNodesJoinModes, testCase.givenLinks)

			err := system.Activate()
			if err != nil {
				errs = append(errs, err)
			}

			afterErrs := load(system, testCase.givenNodesAfterActivation, testCase.givenNodesJoinModesAfterActivation, testCase.givenLinksAfterActivation)
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
	system := New()
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
			system := New()
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
			system := New()
			load(system, testCase.givenNodes, nil, testCase.givenLinks)

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
	givenJoinMode := joinmode.AND

	system := New()
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

	system := New()
	system.AddNode(givenNode)
	system.Activate()

	storedJoinMode := system.JoinModeOfNode(givenNode)

	if joinmode.NONE != storedJoinMode {
		t.Errorf("got: %+v, want: %+v", storedJoinMode, joinmode.NONE)
	}
}

func load(system *NodeSystem, nodes []node.Node, nodesJoinModes map[node.Node]joinmode.JoinMode, links []nodelink.NodeLink) []error {
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
