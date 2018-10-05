package nodelink

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/internal/utils"
	"github.com/rlespinasse/hoff/node"
)

var (
	// NodeLinkComparator is a google/go-cmp comparator of Node Links
	NodeLinkComparator = cmp.Comparer(func(x, y NodeLink) bool {
		return cmp.Equal(x.From, y.From, node.NodeComparator) && cmp.Equal(x.To, y.To, node.NodeComparator) && cmp.Equal(x.Branch, y.Branch)
	})
)

// NodeLink store all information needed to represent a link in the node system
type NodeLink struct {
	From   node.Node
	To     node.Node
	Branch *bool
}

// New create a new link from a node to another node
func New(from, to node.Node) NodeLink {
	return NodeLink{
		From: from,
		To:   to,
	}
}

// NewOnBranch create a new link from a node (and his branch output) to another node
func NewOnBranch(from, to node.Node, branch bool) NodeLink {
	return NodeLink{
		From:   from,
		To:     to,
		Branch: utils.BoolPointer(branch),
	}
}

// String print human-readable version of a node link
func (n NodeLink) String() string {
	branch := ""
	if n.Branch != nil {
		branch = fmt.Sprintf(" branch:%v", *n.Branch)
	}
	return fmt.Sprintf("{from:'%v' to:'%v'%v}", n.From, n.To, branch)
}
