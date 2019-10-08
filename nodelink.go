package hoff

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/internal/utils"
	"github.com/rlespinasse/hoff/node"
)

var (
	// nodeLinkComparator is a google/go-cmp comparator of Node Links
	nodeLinkComparator = cmp.Comparer(func(x, y nodeLink) bool {
		return cmp.Equal(x.From, y.From, node.NodeComparator) && cmp.Equal(x.To, y.To, node.NodeComparator) && cmp.Equal(x.Branch, y.Branch)
	})
)

// nodeLink store all information needed to represent a link in the node system
type nodeLink struct {
	From   node.Node
	To     node.Node
	Branch *bool
}

// newNodeLink create a new link from a node to another node
func newNodeLink(from, to node.Node) nodeLink {
	return nodeLink{
		From: from,
		To:   to,
	}
}

// newNodeLinkOnBranch create a new link from a node (and his branch output) to another node
func newNodeLinkOnBranch(from, to node.Node, branch bool) nodeLink {
	return nodeLink{
		From:   from,
		To:     to,
		Branch: utils.BoolPointer(branch),
	}
}

// String print human-readable version of a node link
func (n nodeLink) String() string {
	branch := ""
	if n.Branch != nil {
		branch = fmt.Sprintf(" branch:%v", *n.Branch)
	}
	return fmt.Sprintf("{from:'%v' to:'%v'%v}", n.From, n.To, branch)
}
