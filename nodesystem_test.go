package namingishard

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var someActionNode, _ = NewActionNode(func(*Context) error { return nil })
var anotherActionNode, _ = NewActionNode(func(*Context) error { return nil })
var someAnotherActionNode, _ = NewActionNode(func(*Context) error { return nil })
var alwaysTrueDecisionNode, _ = NewDecisionNode(func(*Context) (bool, error) { return true, nil })

func Test_NodeSystem_Validate(t *testing.T) {
	testCases := []struct {
		name                               string
		givenNodes                         []Node
		givenNodesJoinModes                map[Node]JoinMode
		givenLinks                         []NodeLink
		givenNodesAfterValidation          []Node
		givenNodesJoinModesAfterValidation map[Node]JoinMode
		givenLinksAfterValidation          []NodeLink
		expectedNodeSystem                 *NodeSystem
		expectedErrors                     []error
	}{
		{
			name: "Can have no nodes",
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes:    []Node{},
				links:    []NodeLink{},
			},
		},
		{
			name: "Can have one action node",
			givenNodes: []Node{
				someActionNode,
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes: []Node{
					someActionNode,
				},
				links: []NodeLink{},
			},
		},
		{
			name: "Can't have the same node more than one time",
			givenNodes: []Node{
				someActionNode,
				someActionNode,
			},
			expectedNodeSystem: &NodeSystem{
				validity: false,
				nodes: []Node{
					someActionNode,
					someActionNode,
				},
				links: []NodeLink{},
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
				validity: false,
				nodes: []Node{
					alwaysTrueDecisionNode,
				},
				links: []NodeLink{},
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
			givenLinks: []NodeLink{
				NewBranchLink(alwaysTrueDecisionNode, someActionNode, true),
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
				},
				links: []NodeLink{
					NewBranchLink(alwaysTrueDecisionNode, someActionNode, true),
				},
			},
		},
		{
			name: "Can have 2 action nodes link together",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []NodeLink{
				NewLink(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes: []Node{
					someActionNode,
					anotherActionNode,
				},
				links: []NodeLink{
					NewLink(someActionNode, anotherActionNode),
				},
			},
		},
		{
			name: "Can't add empty 'from' on branch link",
			givenLinks: []NodeLink{
				NodeLink{to: someActionNode},
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes:    []Node{},
				links:    []NodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have missing 'from' attribute"),
			},
		},
		{
			name: "Can't add empty 'to' on branch link",
			givenLinks: []NodeLink{
				NodeLink{from: someActionNode},
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes:    []Node{},
				links:    []NodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have missing 'to' attribute"),
			},
		},
		{
			name: "Can't add link with the node on 'from' and 'to'",
			givenLinks: []NodeLink{
				NewLink(someActionNode, someActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes:    []Node{},
				links:    []NodeLink{},
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
			givenLinks: []NodeLink{
				NewBranchLink(alwaysTrueDecisionNode, someActionNode, true),
				NewBranchLink(alwaysTrueDecisionNode, anotherActionNode, true),
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
					anotherActionNode,
				},
				links: []NodeLink{
					NewBranchLink(alwaysTrueDecisionNode, someActionNode, true),
					NewBranchLink(alwaysTrueDecisionNode, anotherActionNode, true),
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
				anotherActionNode: JoinModeAnd,
			},
			givenLinks: []NodeLink{
				NewBranchLink(alwaysTrueDecisionNode, anotherActionNode, true),
				NewLink(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{
					anotherActionNode: JoinModeAnd,
				},
				links: []NodeLink{
					NewBranchLink(alwaysTrueDecisionNode, anotherActionNode, true),
					NewLink(someActionNode, anotherActionNode),
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
				anotherActionNode: JoinModeOr,
			},
			givenLinks: []NodeLink{
				NewBranchLink(alwaysTrueDecisionNode, anotherActionNode, true),
				NewLink(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
					anotherActionNode,
				},
				nodesJoinModes: map[Node]JoinMode{
					anotherActionNode: JoinModeOr,
				},
				links: []NodeLink{
					NewBranchLink(alwaysTrueDecisionNode, anotherActionNode, true),
					NewLink(someActionNode, anotherActionNode),
				},
			},
		},
		{
			name: "Can't hava a link with branch who is not needed",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []NodeLink{
				NewBranchLink(someActionNode, anotherActionNode, true),
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes: []Node{
					someActionNode,
					anotherActionNode,
				},
				links: []NodeLink{},
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
			givenLinks: []NodeLink{
				NewLink(alwaysTrueDecisionNode, someActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				validity: false,
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
				},
				links: []NodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have missing branch"),
				fmt.Errorf("can't have decision node without link from it: %+v", alwaysTrueDecisionNode),
			},
		},
		{
			name: "Can't add any node after successful validation",
			givenNodesAfterValidation: []Node{
				someActionNode,
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes:    []Node{},
				links:    []NodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't add node, node system is freeze due to successful validation"),
			},
		},
		{
			name: "Can't add any node join mode after successful validation",
			givenNodesJoinModesAfterValidation: map[Node]JoinMode{
				someActionNode: JoinModeNone,
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes:    []Node{},
				links:    []NodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't add node join mode, node system is freeze due to successful validation"),
			},
		},
		{
			name: "Can't add any branch link after successful validation",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinksAfterValidation: []NodeLink{
				NewLink(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes: []Node{
					someActionNode,
					anotherActionNode,
				},
				links: []NodeLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't add branch link, node system is freeze due to successful validation"),
			},
		},
		{
			name: "Can't have a link with an undeclared node as 'to'",
			givenNodes: []Node{
				someActionNode,
			},
			givenLinks: []NodeLink{
				NewLink(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				validity: false,
				nodes: []Node{
					someActionNode,
				},
				links: []NodeLink{
					NewLink(someActionNode, anotherActionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have undeclared node '%+v' as 'to' in branch link %+v", anotherActionNode, NewLink(someActionNode, anotherActionNode)),
			},
		},
		{
			name: "Can't have a link with an undeclared node as 'from'",
			givenNodes: []Node{
				anotherActionNode,
			},
			givenLinks: []NodeLink{
				NewLink(someActionNode, anotherActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				validity: false,
				nodes: []Node{
					anotherActionNode,
				},
				links: []NodeLink{
					NewLink(someActionNode, anotherActionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have undeclared node '%+v' as 'from' in branch link %+v", someActionNode, NewLink(someActionNode, anotherActionNode)),
			},
		},
		{
			name: "Can add a branch link after a failed validation",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinksAfterValidation: []NodeLink{
				NewBranchLink(alwaysTrueDecisionNode, someActionNode, true),
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
				},
				links: []NodeLink{
					NewBranchLink(alwaysTrueDecisionNode, someActionNode, true),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have decision node without link from it: %+v", alwaysTrueDecisionNode),
			},
		},
		{
			name: "Can add a node after a failed validation",
			givenNodes: []Node{
				someActionNode,
			},
			givenLinks: []NodeLink{
				NewLink(someActionNode, anotherActionNode),
			},
			givenNodesAfterValidation: []Node{
				anotherActionNode,
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes: []Node{
					someActionNode,
					anotherActionNode,
				},
				links: []NodeLink{
					NewLink(someActionNode, anotherActionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have undeclared node '%+v' as 'to' in branch link %+v", anotherActionNode, NewLink(someActionNode, anotherActionNode)),
			},
		},
		{
			name: "Can't have cycle between some links",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []NodeLink{
				NewLink(someActionNode, anotherActionNode),
				NewLink(anotherActionNode, someActionNode),
			},
			expectedNodeSystem: &NodeSystem{
				validity: false,
				nodes: []Node{
					someActionNode,
					anotherActionNode,
				},
				links: []NodeLink{
					NewLink(someActionNode, anotherActionNode),
					NewLink(anotherActionNode, someActionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("Can't have cycle between links: %+v", []NodeLink{
					NewLink(someActionNode, anotherActionNode),
					NewLink(anotherActionNode, someActionNode),
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
			givenLinks: []NodeLink{
				NewBranchLink(alwaysTrueDecisionNode, someActionNode, true),
				NewBranchLink(alwaysTrueDecisionNode, anotherActionNode, false),
				NewLink(anotherActionNode, alwaysTrueDecisionNode),
			},

			expectedNodeSystem: &NodeSystem{
				validity: false,
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
					anotherActionNode,
				},
				links: []NodeLink{
					NewBranchLink(alwaysTrueDecisionNode, someActionNode, true),
					NewBranchLink(alwaysTrueDecisionNode, anotherActionNode, false),
					NewLink(anotherActionNode, alwaysTrueDecisionNode),
				},
			},
			expectedErrors: []error{
				fmt.Errorf("Can't have cycle between links: %+v", []NodeLink{
					NewBranchLink(alwaysTrueDecisionNode, anotherActionNode, false),
					NewLink(anotherActionNode, alwaysTrueDecisionNode),
				}),
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			system := NewNodeSystem()
			var errs []error
			for _, node := range testCase.givenNodes {
				_, err := system.AddNode(node)
				if err != nil {
					errs = append(errs, err)
				}
			}
			for node, mode := range testCase.givenNodesJoinModes {
				_, err := system.AddNodeJoinMode(node, mode)
				if err != nil {
					errs = append(errs, err)
				}
			}
			for _, link := range testCase.givenLinks {
				_, err := system.AddLink(link)
				if err != nil {
					errs = append(errs, err)
				}
			}
			validity, validityErrs := system.Validate()
			errs = append(errs, validityErrs...)

			for _, node := range testCase.givenNodesAfterValidation {
				_, err := system.AddNode(node)
				if err != nil {
					errs = append(errs, err)
				}
			}
			for node, mode := range testCase.givenNodesJoinModesAfterValidation {
				_, err := system.AddNodeJoinMode(node, mode)
				if err != nil {
					errs = append(errs, err)
				}
			}
			for _, link := range testCase.givenLinksAfterValidation {
				_, err := system.AddLink(link)
				if err != nil {
					errs = append(errs, err)
				}
			}
			if testCase.givenNodesAfterValidation != nil || testCase.givenLinksAfterValidation != nil {
				validity, validityErrs = system.Validate()
				errs = append(errs, validityErrs...)
			}

			if validity != testCase.expectedNodeSystem.validity {
				t.Errorf("validity - got: %+v, want: %+v", validity, testCase.expectedNodeSystem.validity)
			}
			if !cmp.Equal(errs, testCase.expectedErrors, equalOptionForError) {
				t.Errorf("errors - got: %+v, want: %+v", errs, testCase.expectedErrors)
			}
			if !cmp.Equal(system, testCase.expectedNodeSystem) {
				t.Errorf("system - got: %+v, want: %+v", system, testCase.expectedNodeSystem)
			}
		})
	}
}

func Test_NodeSystem_activate(t *testing.T) {
	testCases := []struct {
		name                 string
		givenNodes           []Node
		givenNodesJoinModes  map[Node]JoinMode
		givenLinks           []NodeLink
		expectedActivatation bool
		expectedInitialNodes []Node
		expectedLinkTree     map[Node]map[*bool][]Node
		expectedError        error
	}{
		{
			name:                 "Can activate an empty validated system",
			expectedActivatation: true,
			expectedInitialNodes: []Node{},
			expectedLinkTree:     map[Node]map[*bool][]Node{},
		},
		{
			name: "Can't activate an unvalidated system",
			givenLinks: []NodeLink{
				NewLink(someActionNode, anotherActionNode),
			},
			expectedActivatation: false,
			expectedError:        errors.New("can't activate a unvalidated node system"),
		},
		{
			name:                 "Can activate an one node validated system",
			givenNodes:           []Node{someActionNode},
			expectedActivatation: true,
			expectedInitialNodes: []Node{someActionNode},
			expectedLinkTree:     map[Node]map[*bool][]Node{},
		},
		{
			name: "Can activate an no needed branch node link validated system",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []NodeLink{
				NewLink(someActionNode, anotherActionNode),
			},
			expectedActivatation: true,
			expectedInitialNodes: []Node{
				someActionNode,
			},
			expectedLinkTree: map[Node]map[*bool][]Node{
				someActionNode: map[*bool][]Node{
					nil: []Node{anotherActionNode},
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
			givenLinks: []NodeLink{
				NewBranchLink(alwaysTrueDecisionNode, anotherActionNode, true),
			},
			expectedActivatation: true,
			expectedInitialNodes: []Node{
				someActionNode,
				alwaysTrueDecisionNode,
			},
			expectedLinkTree: map[Node]map[*bool][]Node{
				alwaysTrueDecisionNode: map[*bool][]Node{
					truePointer: []Node{anotherActionNode},
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			system := NewNodeSystem()
			for _, node := range testCase.givenNodes {
				system.AddNode(node)
			}
			for node, mode := range testCase.givenNodesJoinModes {
				system.AddNodeJoinMode(node, mode)
			}
			for _, link := range testCase.givenLinks {
				system.AddLink(link)
			}
			system.Validate()
			err := system.activate()
			if system.active != testCase.expectedActivatation {
				t.Errorf("activation - got: %+v, want: %+v", system.active, testCase.expectedActivatation)
			}
			if !cmp.Equal(err, testCase.expectedError, equalOptionForError) {
				t.Errorf("error - got: %+v, want: %+v", err, testCase.expectedError)
			}
			if !cmp.Equal(system.InitialNodes(), testCase.expectedInitialNodes, equalOptionForNode) {
				t.Errorf("initial nodes - got: %+v, want: %+v", system.InitialNodes(), testCase.expectedInitialNodes)
			}
			if !cmp.Equal(system.linkTree, testCase.expectedLinkTree, equalOptionForNode) {
				t.Errorf("node tree - got: %#v, want: %#v", system.linkTree, testCase.expectedLinkTree)
			}
		})
	}
}

func Test_NodeSystem_multiple_activatation(t *testing.T) {
	system := NewNodeSystem()
	system.Validate()
	system.activate()
	err := system.activate()
	if err != nil {
		t.Errorf("Can be activate multiple times without errors, but got: %+v", err)
	}
}

func Test_NodeSystem_follow(t *testing.T) {
	testCases := []struct {
		name                   string
		givenNodes             []Node
		givenLinks             []NodeLink
		givenFollowNode        Node
		givenFollowBranch      *bool
		expectedFollowingNodes []Node
		expectedError          error
	}{
		{
			name: "Can't follow on an unactivated system",
			givenLinks: []NodeLink{
				NewLink(someActionNode, anotherActionNode),
			},
			expectedError: errors.New("can't follow a node if system is not activated"),
		},
		{
			name: "Can follow 'from' on link",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []NodeLink{
				NewLink(someActionNode, anotherActionNode),
			},
			givenFollowNode:        someActionNode,
			expectedFollowingNodes: []Node{anotherActionNode},
		},
		{
			name: "Can follow 'to' on link",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenLinks: []NodeLink{
				NewLink(someActionNode, anotherActionNode),
			},
			givenFollowNode:        anotherActionNode,
			expectedFollowingNodes: nil,
		},
		{
			name: "Can follow 'from' on branch link",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []NodeLink{
				NewBranchLink(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenFollowNode:        alwaysTrueDecisionNode,
			givenFollowBranch:      boolPointer(true),
			expectedFollowingNodes: []Node{someActionNode},
		},
		{
			name: "Can follow 'to' on branch link",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []NodeLink{
				NewBranchLink(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenFollowNode:        someActionNode,
			expectedFollowingNodes: nil,
		},
		{
			name: "Can't follow 'from' on branch link but without passing the branch",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []NodeLink{
				NewBranchLink(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenFollowNode:        alwaysTrueDecisionNode,
			expectedFollowingNodes: nil,
		},
		{
			name: "Can't follow 'from' on branch link but without passing the right branch",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenLinks: []NodeLink{
				NewBranchLink(alwaysTrueDecisionNode, someActionNode, true),
			},
			givenFollowNode:        alwaysTrueDecisionNode,
			givenFollowBranch:      boolPointer(false),
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
				system.AddLink(link)
			}
			system.Validate()
			system.activate()
			nodes, err := system.follow(testCase.givenFollowNode, testCase.givenFollowBranch)

			if !cmp.Equal(err, testCase.expectedError, equalOptionForError) {
				t.Errorf("error - got: %+v, want: %+v", err, testCase.expectedError)
			}
			if !cmp.Equal(nodes, testCase.expectedFollowingNodes, equalOptionForNode) {
				t.Errorf("following nodes - got: %+v, want: %+v", nodes, testCase.expectedFollowingNodes)
			}
		})
	}
}
