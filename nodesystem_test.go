package flow

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var someActionNode, _ = NewActionNode(func(*Context) (bool, error) { return true, nil })
var anotherActionNode, _ = NewActionNode(func(*Context) (bool, error) { return true, nil })
var alwaysTrueDecisionNode, _ = NewDecisionNode(func(*Context) (bool, error) { return true, nil })

func Test_NodeSystem_Validate(t *testing.T) {
	testCases := []struct {
		name                            string
		givenNodes                      []Node
		givenBranchLinks                []NodeBranchLink
		givenNodesAfterValidation       []Node
		givenBranchLinksAfterValidation []NodeBranchLink
		expectedNodeSystem              *NodeSystem
		expectedErrors                  []error
	}{
		{
			name: "Can have no nodes",
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes:    []Node{},
				links:    []NodeBranchLink{},
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
				links: []NodeBranchLink{},
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
				links: []NodeBranchLink{},
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
				links: []NodeBranchLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have orphan decision node: %+v", alwaysTrueDecisionNode),
			},
		},
		{
			name: "Can have a link between decision node and action node",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From:   alwaysTrueDecisionNode,
					To:     someActionNode,
					Branch: ptrOfString("true"),
				},
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
				},
				links: []NodeBranchLink{
					NodeBranchLink{
						From:   alwaysTrueDecisionNode,
						To:     someActionNode,
						Branch: ptrOfString("true"),
					},
				},
			},
		},
		{
			name: "Can have 2 action nodes link together",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: someActionNode,
					To:   anotherActionNode,
				},
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes: []Node{
					someActionNode,
					anotherActionNode,
				},
				links: []NodeBranchLink{
					NodeBranchLink{
						From: someActionNode,
						To:   anotherActionNode,
					},
				},
			},
		},
		{
			name: "Can't add empty link",
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{},
			},
			expectedNodeSystem: &NodeSystem{
				validity: false,
				nodes:    []Node{},
				links: []NodeBranchLink{
					NodeBranchLink{},
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have missing 'From' attribute: %+v", NodeBranchLink{}),
				fmt.Errorf("can't have missing 'To' attribute: %+v", NodeBranchLink{}),
			},
		},
		{
			name: "Can't have a link with unknown branch",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From:   alwaysTrueDecisionNode,
					To:     someActionNode,
					Branch: ptrOfString("some_branch"),
				},
			},
			expectedNodeSystem: &NodeSystem{
				validity: false,
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
				},
				links: []NodeBranchLink{
					NodeBranchLink{
						From:   alwaysTrueDecisionNode,
						To:     someActionNode,
						Branch: ptrOfString("some_branch"),
					},
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have unknown branch: 'some_branch', from %+v, available branches %+v", alwaysTrueDecisionNode, alwaysTrueDecisionNode.AvailableBranches()),
			},
		},
		{
			name: "Can't hava a link with branch who is not needed",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From:   someActionNode,
					To:     anotherActionNode,
					Branch: ptrOfString("not_needed_branch"),
				},
			},
			expectedNodeSystem: &NodeSystem{
				validity: false,
				nodes: []Node{
					someActionNode,
					anotherActionNode,
				},
				links: []NodeBranchLink{
					NodeBranchLink{
						From:   someActionNode,
						To:     anotherActionNode,
						Branch: ptrOfString("not_needed_branch"),
					},
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have not needed branch: 'not_needed_branch', from %+v", someActionNode),
			},
		},
		{
			name: "Can't have a link with a missing branch",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: alwaysTrueDecisionNode,
					To:   someActionNode,
				},
			},
			expectedNodeSystem: &NodeSystem{
				validity: false,
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
				},
				links: []NodeBranchLink{
					NodeBranchLink{
						From: alwaysTrueDecisionNode,
						To:   someActionNode,
					},
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have missing branch from %+v, available branches %+v", alwaysTrueDecisionNode, alwaysTrueDecisionNode.AvailableBranches()),
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
				links:    []NodeBranchLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't add node, node system is freeze due to successful validation"),
			},
		},
		{
			name: "Can't add any branch link after successful validation",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenBranchLinksAfterValidation: []NodeBranchLink{
				NodeBranchLink{
					From: someActionNode,
					To:   anotherActionNode,
				},
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes: []Node{
					someActionNode,
					anotherActionNode,
				},
				links: []NodeBranchLink{},
			},
			expectedErrors: []error{
				fmt.Errorf("can't add branch link, node system is freeze due to successful validation"),
			},
		},
		{
			name: "Can't have a link with an undeclared node as 'To'",
			givenNodes: []Node{
				someActionNode,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: someActionNode,
					To:   anotherActionNode,
				},
			},
			expectedNodeSystem: &NodeSystem{
				validity: false,
				nodes: []Node{
					someActionNode,
				},
				links: []NodeBranchLink{
					NodeBranchLink{
						From: someActionNode,
						To:   anotherActionNode,
					},
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have undeclared node '%+v' as 'To' in branch link %+v", anotherActionNode, NodeBranchLink{
					From: someActionNode,
					To:   anotherActionNode,
				}),
			},
		},
		{
			name: "Can't have a link with an undeclared node as 'From'",
			givenNodes: []Node{
				anotherActionNode,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: someActionNode,
					To:   anotherActionNode,
				},
			},
			expectedNodeSystem: &NodeSystem{
				validity: false,
				nodes: []Node{
					anotherActionNode,
				},
				links: []NodeBranchLink{
					NodeBranchLink{
						From: someActionNode,
						To:   anotherActionNode,
					},
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have undeclared node '%+v' as 'From' in branch link %+v", someActionNode, NodeBranchLink{
					From: someActionNode,
					To:   anotherActionNode,
				}),
			},
		},
		{
			name: "Can add a branch link after a failed validation",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenBranchLinksAfterValidation: []NodeBranchLink{
				NodeBranchLink{
					From:   alwaysTrueDecisionNode,
					To:     someActionNode,
					Branch: ptrOfString("true"),
				},
			},
			expectedNodeSystem: &NodeSystem{
				validity: true,
				nodes: []Node{
					alwaysTrueDecisionNode,
					someActionNode,
				},
				links: []NodeBranchLink{
					NodeBranchLink{
						From:   alwaysTrueDecisionNode,
						To:     someActionNode,
						Branch: ptrOfString("true"),
					},
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have orphan decision node: %+v", alwaysTrueDecisionNode),
			},
		},
		{
			name: "Can add a node after a failed validation",
			givenNodes: []Node{
				someActionNode,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: someActionNode,
					To:   anotherActionNode,
				},
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
				links: []NodeBranchLink{
					NodeBranchLink{
						From: someActionNode,
						To:   anotherActionNode,
					},
				},
			},
			expectedErrors: []error{
				fmt.Errorf("can't have undeclared node '%+v' as 'To' in branch link %+v", anotherActionNode, NodeBranchLink{
					From: someActionNode,
					To:   anotherActionNode,
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
			for _, link := range testCase.givenBranchLinks {
				_, err := system.AddBranchLink(link)
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
			for _, link := range testCase.givenBranchLinksAfterValidation {
				_, err := system.AddBranchLink(link)
				if err != nil {
					errs = append(errs, err)
				}
			}
			if testCase.givenNodesAfterValidation != nil || testCase.givenBranchLinksAfterValidation != nil {
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
		givenBranchLinks     []NodeBranchLink
		expectedActivatation bool
		expectedInitialNodes []Node
		expectedNodeTree     map[Node]map[string]Node
		expectedError        error
	}{
		{
			name:                 "Can activate an empty validated system",
			expectedActivatation: true,
			expectedInitialNodes: []Node{},
			expectedNodeTree:     map[Node]map[string]Node{},
		},
		{
			name: "Can't activate an unvalidated system",
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: someActionNode,
					To:   anotherActionNode,
				},
			},
			expectedActivatation: false,
			expectedError:        errors.New("can't activate a unvalidated node system"),
		},
		{
			name:                 "Can activate an one node validated system",
			givenNodes:           []Node{someActionNode},
			expectedActivatation: true,
			expectedInitialNodes: []Node{someActionNode},
			expectedNodeTree:     map[Node]map[string]Node{},
		},
		{
			name: "Can activate an no needed branch node link validated system",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: someActionNode,
					To:   anotherActionNode,
				},
			},
			expectedActivatation: true,
			expectedInitialNodes: []Node{
				someActionNode,
			},
			expectedNodeTree: map[Node]map[string]Node{
				someActionNode: map[string]Node{
					"*": anotherActionNode,
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
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From:   alwaysTrueDecisionNode,
					To:     anotherActionNode,
					Branch: ptrOfString("true"),
				},
			},
			expectedActivatation: true,
			expectedInitialNodes: []Node{
				someActionNode,
				alwaysTrueDecisionNode,
			},
			expectedNodeTree: map[Node]map[string]Node{
				alwaysTrueDecisionNode: map[string]Node{
					"true": anotherActionNode,
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
			for _, link := range testCase.givenBranchLinks {
				system.AddBranchLink(link)
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
			if !cmp.Equal(system.nodesTree, testCase.expectedNodeTree, equalOptionForNode) {
				t.Errorf("node tree - got: %#v, want: %#v", system.nodesTree, testCase.expectedNodeTree)
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
		name                  string
		givenNodes            []Node
		givenBranchLinks      []NodeBranchLink
		givenFollowNode       Node
		givenFollowBranch     *string
		expectedFollowingNode Node
		expectedError         error
	}{
		{
			name: "Can't follow on an unactivated system",
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: someActionNode,
					To:   anotherActionNode,
				},
			},
			expectedError: errors.New("can't follow a node if system is not activated"),
		},
		{
			name: "Can follow 'From' on no branch link",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: someActionNode,
					To:   anotherActionNode,
				},
			},
			givenFollowNode:       someActionNode,
			expectedFollowingNode: anotherActionNode,
		},
		{
			name: "Can follow 'To' on no branch link",
			givenNodes: []Node{
				someActionNode,
				anotherActionNode,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: someActionNode,
					To:   anotherActionNode,
				},
			},
			givenFollowNode:       anotherActionNode,
			expectedFollowingNode: nil,
		},
		{
			name: "Can follow 'From' on branch link",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From:   alwaysTrueDecisionNode,
					To:     someActionNode,
					Branch: ptrOfString("true"),
				},
			},
			givenFollowNode:       alwaysTrueDecisionNode,
			givenFollowBranch:     ptrOfString("true"),
			expectedFollowingNode: someActionNode,
		},
		{
			name: "Can follow 'To' on branch link",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From:   alwaysTrueDecisionNode,
					To:     someActionNode,
					Branch: ptrOfString("true"),
				},
			},
			givenFollowNode:       someActionNode,
			expectedFollowingNode: nil,
		},
		{
			name: "Can't follow 'From' on branch link but without passing the branch",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From:   alwaysTrueDecisionNode,
					To:     someActionNode,
					Branch: ptrOfString("true"),
				},
			},
			givenFollowNode:       alwaysTrueDecisionNode,
			expectedFollowingNode: nil,
		},
		{
			name: "Can't follow 'From' on branch link but without passing the right branch",
			givenNodes: []Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From:   alwaysTrueDecisionNode,
					To:     someActionNode,
					Branch: ptrOfString("true"),
				},
			},
			givenFollowNode:       alwaysTrueDecisionNode,
			givenFollowBranch:     ptrOfString("false"),
			expectedFollowingNode: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			system := NewNodeSystem()
			for _, node := range testCase.givenNodes {
				system.AddNode(node)
			}
			for _, link := range testCase.givenBranchLinks {
				system.AddBranchLink(link)
			}
			system.Validate()
			system.activate()
			node, err := system.follow(testCase.givenFollowNode, testCase.givenFollowBranch)

			if !cmp.Equal(err, testCase.expectedError, equalOptionForError) {
				t.Errorf("error - got: %+v, want: %+v", err, testCase.expectedError)
			}
			if !cmp.Equal(node, testCase.expectedFollowingNode, equalOptionForNode) {
				t.Errorf("initial nodes - got: %+v, want: %+v", node, testCase.expectedFollowingNode)
			}
		})
	}
}
