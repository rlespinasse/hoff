package flow

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var someActionNode = NewActionNode(func(*Context) (bool, error) { return true, nil })
var anotherActionNode = NewActionNode(func(*Context) (bool, error) { return true, nil })
var alwaysTrueDecisionNode = NewDecisionNode(func(*Context) (bool, error) { return true, nil })

func Test_NodeSystem(t *testing.T) {
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
				fmt.Errorf("orphan decision node: %+v", alwaysTrueDecisionNode),
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
				fmt.Errorf("missing 'From' attribute: %+v", NodeBranchLink{}),
				fmt.Errorf("missing 'To' attribute: %+v", NodeBranchLink{}),
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
				fmt.Errorf("unknown branch: 'some_branch', from %+v, available branches %+v", alwaysTrueDecisionNode, alwaysTrueDecisionNode.AvailableBranches()),
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
				fmt.Errorf("not needed branch: 'not_needed_branch', from %+v", someActionNode),
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
				fmt.Errorf("missing branch from %+v, available branches %+v", alwaysTrueDecisionNode, alwaysTrueDecisionNode.AvailableBranches()),
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
				fmt.Errorf("undeclared node '%+v' as 'To' in branch link %+v", anotherActionNode, NodeBranchLink{
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
				fmt.Errorf("undeclared node '%+v' as 'From' in branch link %+v", someActionNode, NodeBranchLink{
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
				fmt.Errorf("orphan decision node: %+v", alwaysTrueDecisionNode),
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
				fmt.Errorf("undeclared node '%+v' as 'To' in branch link %+v", anotherActionNode, NodeBranchLink{
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
			if !cmp.Equal(errs, testCase.expectedErrors, errorEqualOpts) {
				t.Errorf("errors - got: %+v, want: %+v", errs, testCase.expectedErrors)
			}
			if !cmp.Equal(system, testCase.expectedNodeSystem) {
				t.Errorf("system - got: %+v, want: %+v", system, testCase.expectedNodeSystem)
			}
		})
	}
}

var errorEqualOpts = cmp.Comparer(func(x, y error) bool {
	if x == nil || y == nil {
		return x == nil && y == nil
	}
	return x.Error() == y.Error()
})
