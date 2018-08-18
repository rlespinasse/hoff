package flow

import "github.com/google/go-cmp/cmp"

func ptrOfString(value string) *string {
	return &value
}

var equalOptionForError = cmp.Comparer(func(x, y error) bool {
	return ((x == nil || y == nil) && x == nil && y == nil) || x.Error() == y.Error()
})
