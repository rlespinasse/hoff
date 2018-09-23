package namingishard

import "github.com/google/go-cmp/cmp"

var (
	truePointer  = newBool(true)
	falsePointer = newBool(false)
)

func newBool(value bool) *bool {
	return &value
}

func boolPointer(value bool) *bool {
	if value {
		return truePointer
	}
	return falsePointer
}

var equalOptionForError = cmp.Comparer(func(x, y error) bool {
	return ((x == nil || y == nil) && x == nil && y == nil) || x.Error() == y.Error()
})
