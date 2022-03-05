package verifier

import (
	"encoding"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

var stringToString = func(value string) asStringResult {
	return newStringResult(value)
}

var intToString = func(value int) asStringResult {
	return newStringResult(fmt.Sprintf("%d", value))
}

var int8ToString = func(value int8) asStringResult {
	return newStringResult(fmt.Sprintf("%d", value))
}

var int16ToString = func(value int16) asStringResult {
	return newStringResult(fmt.Sprintf("%d", value))
}

var int32ToString = func(value int32) asStringResult {
	return newStringResult(fmt.Sprintf("%d", value))
}

var int64ToString = func(value int64) asStringResult {
	return newStringResult(fmt.Sprintf("%d", value))
}

var uIntToString = func(value uint) asStringResult {
	return newStringResult(fmt.Sprintf("%d", value))
}

var uInt8ToString = func(value uint8) asStringResult {
	return newStringResult(fmt.Sprintf("%d", value))
}

var uInt16ToString = func(value uint16) asStringResult {
	return newStringResult(fmt.Sprintf("%d", value))
}

var uInt32ToString = func(value uint32) asStringResult {
	return newStringResult(fmt.Sprintf("%d", value))
}

var uInt64ToString = func(value uint64) asStringResult {
	return newStringResult(fmt.Sprintf("%d", value))
}

var boolToString = func(value bool) asStringResult {
	if value {
		return newStringResult("True")
	}
	return newStringResult("False")
}

var float64ToString = func(value float64) asStringResult {
	return newStringResult(strconv.FormatFloat(value, 'f', -1, 64))
}

var float32ToString = func(value float32) asStringResult {
	return newStringResult(strconv.FormatFloat(float64(value), 'f', -1, 32))
}

var timeToString = func(value time.Time) asStringResult {
	return newStringResult(value.Format(time.RFC3339))
}

var uUIDToString = func(value uuid.UUID) asStringResult {
	return newStringResult(value.String())
}

var stringerToString = func(stringer fmt.Stringer) asStringResult {
	return newStringResult(stringer.String())
}

var textMarshallerToString = func(marshaller encoding.TextMarshaler) (asStringResult, bool) {
	b, err := marshaller.MarshalText()
	if err == nil {
		return newStringResult(string(b)), true
	}
	return asStringResult{}, false
}

var stringBuilderToString = func(builder strings.Builder) (asStringResult, bool) {
	return newStringResult(builder.String()), true
}
