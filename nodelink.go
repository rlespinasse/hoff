package flows

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

type NodeLink struct {
	from   Node
	to     Node
	branch *string
}

func NewLink(from, to Node) NodeLink {
	return NodeLink{
		from: from,
		to:   to,
	}
}

func NewBranchLink(from, to Node, branch string) NodeLink {
	return NodeLink{
		from:   from,
		to:     to,
		branch: stringPointer(branch),
	}
}

func (n NodeLink) String() string {
	branch := ""
	if n.branch != nil {
		branch = fmt.Sprintf(", branch: %v", *n.branch)
	}
	return fmt.Sprintf("{from: %v, to: %v%v}", n.from, n.to, branch)
}

var equalOptionForNodeLink = cmp.Comparer(func(x, y NodeLink) bool {
	return cmp.Equal(x.from, y.from, equalOptionForNode) && cmp.Equal(x.to, y.to, equalOptionForNode) && cmp.Equal(x.branch, y.branch)
})
