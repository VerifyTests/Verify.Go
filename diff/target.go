package diff

import (
	"fmt"
	"os"
	"strings"
)

type targetPosition struct {
	TargetOnLeft bool
}

var position = newTargetPosition()

func newTargetPosition() targetPosition {
	pos := targetPosition{}
	onLeft, found := pos.ReadTargetOnLeft()
	if !found {
		onLeft = false
	}
	pos.TargetOnLeft = onLeft
	return pos
}

func (t *targetPosition) ReadTargetOnLeft() (result bool, found bool) {
	value, ok := os.LookupEnv("DiffEngine_TargetOnLeft")
	if !ok {
		return false, false
	}

	if strings.ToLower(value) == "true" {
		return true, true
	}

	if strings.ToLower(value) == "false" {
		return true, false
	}

	panic(fmt.Sprintf("Unable to parse Position from `DiffEngine_TargetOnLeft`. Must be `true` or `false`. Environment variable: %s", value))
}

func (t *targetPosition) SetTargetOnLeft(value bool) {
	if t.TargetOnLeft == value {
		return
	}

	t.TargetOnLeft = value
	var envValue string

	if value {
		envValue = "true"
	}

	_ = os.Setenv("DiffEngine_TargetOnLeft", envValue)
}
