package verifier

import (
	"testing"
)

func TestRemoveItemFromSlice(t *testing.T) {
	items := []string{"first", "second", "third", "fourth"}
	items = removeStringItem(items, "second")

	if len(items) != 3 {
		t.Fatalf("Should contain 3 items")
	}
}
