package nodelink

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/rlespinasse/hoff/computestate"
	"github.com/rlespinasse/hoff/node"
)

type SomeNode struct{}

func (n *SomeNode) Compute(c *node.Context) computestate.ComputeState {
	c.Store("message", "SomeNode is passing")
	return computestate.Continue()
}

func (n *SomeNode) DecideCapability() bool {
	return false
}

func Test_New(t *testing.T) {
	givenFromNode := &SomeNode{}
	givenToNode := &SomeNode{}
	expectedString := "{from:'&{}' to:'&{}'}"

	link := New(givenFromNode, givenToNode)
	linkString := link.String()
	if linkString != expectedString {
		t.Errorf("got: %+v, want: %+v", linkString, expectedString)
	}
}

func Test_NewOnBranch(t *testing.T) {
	givenFromNode := &SomeNode{}
	givenToNode := &SomeNode{}
	givenBranch := true
	expectedString := "{from:'&{}' to:'&{}' branch:true}"

	link := NewOnBranch(givenFromNode, givenToNode, givenBranch)
	linkString := link.String()
	if linkString != expectedString {
		t.Errorf("got: %+v, want: %+v", linkString, expectedString)
	}
}

func Test_NodeLinkComparator_Equal(t *testing.T) {
	givenFromNode := &SomeNode{}
	givenToNode := &SomeNode{}
	givenBranch := true

	link := NewOnBranch(givenFromNode, givenToNode, givenBranch)
	anotherLink := NewOnBranch(givenFromNode, givenToNode, givenBranch)

	if !cmp.Equal(link, anotherLink, NodeLinkComparator) {
		t.Errorf("link: %+v and anotherLink: %+v must be equals", link, anotherLink)
	}
}

func Test_NodeLinkComparator_NotEqual(t *testing.T) {
	givenFromNode := &SomeNode{}
	givenToNode := &SomeNode{}
	givenBranch := true

	link := New(givenFromNode, givenToNode)
	anotherLink := NewOnBranch(givenFromNode, givenToNode, givenBranch)

	if cmp.Equal(link, anotherLink, NodeLinkComparator) {
		t.Errorf("link: %+v and anotherLink: %+v must not be equals", link, anotherLink)
	}
}
