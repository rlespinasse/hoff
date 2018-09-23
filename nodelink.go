package namingishard

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

type JoinMode string

const (
	JoinModeAnd  JoinMode = "and"
	JoinModeOr            = "or"
	JoinModeNone          = "none"
)

type NodeLink struct {
	from   Node
	to     Node
	branch *bool
}

func NewLink(from, to Node) NodeLink {
	return NodeLink{
		from: from,
		to:   to,
	}
}

func NewBranchLink(from, to Node, branch bool) NodeLink {
	return NodeLink{
		from:   from,
		to:     to,
		branch: boolPointer(branch),
	}
}

func (n NodeLink) String() string {
	branch := ""
	if n.branch != nil {
		branch = fmt.Sprintf(", branch: %v", *n.branch)
	}
	return fmt.Sprintf("{from: '%v', to: '%v'%v}", n.from, n.to, branch)
}

var equalOptionForNodeLink = cmp.Comparer(func(x, y NodeLink) bool {
	return cmp.Equal(x.from, y.from, equalOptionForNode) && cmp.Equal(x.to, y.to, equalOptionForNode) && cmp.Equal(x.branch, y.branch)
})
