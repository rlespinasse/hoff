package namingishard

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

type linkKind string

const (
	classicLink linkKind = "classic"
	forkLink             = "fork"
	joinLink             = "join"
	mergeLink            = "merge"
	noLink               = "none"
)

type NodeLink struct {
	from   Node
	to     Node
	branch *string
	kind   linkKind
}

func NewLink(from, to Node) NodeLink {
	return NodeLink{
		from: from,
		to:   to,
		kind: classicLink,
	}
}

func NewBranchLink(from, to Node, branch string) NodeLink {
	return NodeLink{
		from:   from,
		to:     to,
		branch: stringPointer(branch),
		kind:   classicLink,
	}
}

func NewForkLink(from, to Node) NodeLink {
	return NodeLink{
		from: from,
		to:   to,
		kind: forkLink,
	}
}

func NewBranchForkLink(from, to Node, branch string) NodeLink {
	return NodeLink{
		from:   from,
		to:     to,
		branch: stringPointer(branch),
		kind:   forkLink,
	}
}

func NewJoinLink(from, to Node) NodeLink {
	return NodeLink{
		from: from,
		to:   to,
		kind: joinLink,
	}
}

func NewBranchJoinLink(from, to Node, branch string) NodeLink {
	return NodeLink{
		from:   from,
		to:     to,
		branch: stringPointer(branch),
		kind:   joinLink,
	}
}

func NewMergeLink(from, to Node) NodeLink {
	return NodeLink{
		from: from,
		to:   to,
		kind: mergeLink,
	}
}

func NewBranchMergeLink(from, to Node, branch string) NodeLink {
	return NodeLink{
		from:   from,
		to:     to,
		branch: stringPointer(branch),
		kind:   mergeLink,
	}
}

func (n NodeLink) String() string {
	branch := ""
	if n.branch != nil {
		branch = fmt.Sprintf(", branch: %v", *n.branch)
	}
	return fmt.Sprintf("{kind: %v, from: %v, to: %v%v}", n.kind, n.from, n.to, branch)
}

var equalOptionForNodeLink = cmp.Comparer(func(x, y NodeLink) bool {
	return cmp.Equal(x.from, y.from, equalOptionForNode) && cmp.Equal(x.to, y.to, equalOptionForNode) && cmp.Equal(x.branch, y.branch)
})
