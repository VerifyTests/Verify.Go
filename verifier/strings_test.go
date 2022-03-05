package verifier

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveItemFromSlice(t *testing.T) {
	items := []string{"first", "second", "third", "fourth"}
	items = removeStringItem(items, "second")
	assert.Len(t, items, 3)
}
