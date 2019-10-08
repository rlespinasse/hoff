package hoff

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_newNodeLink(t *testing.T) {
	givenFromNode := &SomeNode{}
	givenToNode := &SomeNode{}
	expectedString := "{from:'&{}' to:'&{}'}"

	link := newNodeLink(givenFromNode, givenToNode)
	linkString := link.String()
	if linkString != expectedString {
		t.Errorf("got: %+v, want: %+v", linkString, expectedString)
	}
}

func Test_newNodeLinkOnBranch(t *testing.T) {
	givenFromNode := &SomeNode{}
	givenToNode := &SomeNode{}
	givenBranch := true
	expectedString := "{from:'&{}' to:'&{}' branch:true}"

	link := newNodeLinkOnBranch(givenFromNode, givenToNode, givenBranch)
	linkString := link.String()
	if linkString != expectedString {
		t.Errorf("got: %+v, want: %+v", linkString, expectedString)
	}
}

func Test_nodeLinkComparator_Equal(t *testing.T) {
	givenFromNode := &SomeNode{}
	givenToNode := &SomeNode{}
	givenBranch := true

	link := newNodeLinkOnBranch(givenFromNode, givenToNode, givenBranch)
	anotherLink := newNodeLinkOnBranch(givenFromNode, givenToNode, givenBranch)

	if !cmp.Equal(link, anotherLink, nodeLinkComparator) {
		t.Errorf("link: %+v and anotherLink: %+v must be equals", link, anotherLink)
	}
}

func Test_nodeLinkComparator_NotEqual(t *testing.T) {
	givenFromNode := &SomeNode{}
	givenToNode := &SomeNode{}
	givenBranch := true

	link := newNodeLink(givenFromNode, givenToNode)
	anotherLink := newNodeLinkOnBranch(givenFromNode, givenToNode, givenBranch)

	if cmp.Equal(link, anotherLink, nodeLinkComparator) {
		t.Errorf("link: %+v and anotherLink: %+v must not be equals", link, anotherLink)
	}
}
