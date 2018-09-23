package namingishard

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

type nodeLink struct {
	from   Node
	to     Node
	branch *bool
}

func newLink(from, to Node) nodeLink {
	return nodeLink{
		from: from,
		to:   to,
	}
}

func newBranchLink(from, to Node, branch bool) nodeLink {
	return nodeLink{
		from:   from,
		to:     to,
		branch: boolPointer(branch),
	}
}

func (n nodeLink) String() string {
	branch := ""
	if n.branch != nil {
		branch = fmt.Sprintf(" branch:%v", *n.branch)
	}
	return fmt.Sprintf("{from:'%v' to:'%v'%v}", n.from, n.to, branch)
}

var equalOptionFornodeLink = cmp.Comparer(func(x, y nodeLink) bool {
	return cmp.Equal(x.from, y.from, equalOptionForNode) && cmp.Equal(x.to, y.to, equalOptionForNode) && cmp.Equal(x.branch, y.branch)
})
