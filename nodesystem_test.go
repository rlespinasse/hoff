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
			"No nodes",
			[]Node{},
			[]NodeBranchLink{},
			[]Node{},
			[]NodeBranchLink{},
			&NodeSystem{
				validity: true,
				nodes:    []Node{},
				links:    []NodeBranchLink{},
			},
			nil,
		},
		{
			"One action node",
			[]Node{
				someActionNode,
			},
			[]NodeBranchLink{},
			[]Node{},
			[]NodeBranchLink{},
			&NodeSystem{
				validity: true,
				nodes: []Node{
					someActionNode,
				},
				links: []NodeBranchLink{},
			},
			nil,
		},
		{
			"One decision node",
			[]Node{
				alwaysTrueDecisionNode,
			},
			[]NodeBranchLink{},
			[]Node{},
			[]NodeBranchLink{},
			&NodeSystem{
				validity: false,
				nodes: []Node{
					alwaysTrueDecisionNode,
				},
				links: []NodeBranchLink{},
			},
			[]error{
				fmt.Errorf("orphan decision node: %+v", alwaysTrueDecisionNode),
			},
		},
		{
			"Decision node and action node",
			[]Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			[]NodeBranchLink{
				NodeBranchLink{
					From:   alwaysTrueDecisionNode,
					To:     someActionNode,
					Branch: ptrOfString("true"),
				},
			},
			[]Node{},
			[]NodeBranchLink{},
			&NodeSystem{
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
			nil,
		},
		{
			"Two action nodes link together",
			[]Node{
				someActionNode,
				anotherActionNode,
			},
			[]NodeBranchLink{
				NodeBranchLink{
					From: someActionNode,
					To:   anotherActionNode,
				},
			},
			[]Node{},
			[]NodeBranchLink{},
			&NodeSystem{
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
			nil,
		},
		{
			"Can't declare empty link",
			[]Node{},
			[]NodeBranchLink{
				NodeBranchLink{},
			},
			[]Node{},
			[]NodeBranchLink{},
			&NodeSystem{
				validity: false,
				nodes:    []Node{},
				links: []NodeBranchLink{
					NodeBranchLink{},
				},
			},
			[]error{
				fmt.Errorf("missing 'From' attribute: %+v", NodeBranchLink{}),
				fmt.Errorf("missing 'To' attribute: %+v", NodeBranchLink{}),
			},
		},
		{
			"Link with unknown branch",
			[]Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			[]NodeBranchLink{
				NodeBranchLink{
					From:   alwaysTrueDecisionNode,
					To:     someActionNode,
					Branch: ptrOfString("some_branch"),
				},
			},
			[]Node{},
			[]NodeBranchLink{},
			&NodeSystem{
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
			[]error{
				fmt.Errorf("unknown branch: 'some_branch', from %+v, available branches %+v", alwaysTrueDecisionNode, alwaysTrueDecisionNode.AvailableBranches()),
			},
		},
		{
			"Link with branch who is not needed",
			[]Node{
				someActionNode,
				anotherActionNode,
			},
			[]NodeBranchLink{
				NodeBranchLink{
					From:   someActionNode,
					To:     anotherActionNode,
					Branch: ptrOfString("not_needed_branch"),
				},
			},
			[]Node{},
			[]NodeBranchLink{},
			&NodeSystem{
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
			[]error{
				fmt.Errorf("not needed branch: 'not_needed_branch', from %+v", someActionNode),
			},
		},
		{
			"Link with missing branch",
			[]Node{
				alwaysTrueDecisionNode,
				someActionNode,
			},
			[]NodeBranchLink{
				NodeBranchLink{
					From: alwaysTrueDecisionNode,
					To:   someActionNode,
				},
			},
			[]Node{},
			[]NodeBranchLink{},
			&NodeSystem{
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
			[]error{
				fmt.Errorf("missing branch from %+v, available branches %+v", alwaysTrueDecisionNode, alwaysTrueDecisionNode.AvailableBranches()),
			},
		},
		{
			"Must not add node after validation",
			[]Node{},
			[]NodeBranchLink{},
			[]Node{
				someActionNode,
			},
			[]NodeBranchLink{},
			&NodeSystem{
				validity: true,
				nodes:    []Node{},
				links:    []NodeBranchLink{},
			},
			[]error{
				fmt.Errorf("can't add node, node system is freeze due to successful validation"),
			},
		},
		{
			"Must not add branch link after validation",
			[]Node{
				someActionNode,
				anotherActionNode,
			},
			[]NodeBranchLink{},
			[]Node{},
			[]NodeBranchLink{
				NodeBranchLink{
					From: someActionNode,
					To:   anotherActionNode,
				},
			},
			&NodeSystem{
				validity: true,
				nodes: []Node{
					someActionNode,
					anotherActionNode,
				},
				links: []NodeBranchLink{},
			},
			[]error{
				fmt.Errorf("can't add branch link, node system is freeze due to successful validation"),
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
