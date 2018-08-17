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
		name                   string
		givenNodes             []Node
		givenBranchLinks       []NodeBranchLink
		expectedNodeSystem     *NodeSystem
		expectedValidityErrors []error
	}{
		{
			"No nodes",
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
			validity, errs := system.Validate()

			if validity != testCase.expectedNodeSystem.validity {
				t.Errorf("validity - got: %+v, want: %+v", validity, testCase.expectedNodeSystem.validity)
			}
			if !cmp.Equal(errs, testCase.expectedValidityErrors, errorEqualOpts) {
				t.Errorf("validity errors - got: %+v, want: %+v", errs, testCase.expectedValidityErrors)
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
