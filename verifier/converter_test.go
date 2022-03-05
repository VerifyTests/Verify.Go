package verifier

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testMarshaller struct {
	data string
}

type testStringer struct {
	data string
}

func (t testStringer) String() string {
	return t.data
}

func (t testMarshaller) MarshalText() (text []byte, err error) {
	return []byte(t.data), nil
}

func TestStringVerifier(t *testing.T) {

	settings := newSettings()
	inner := newInnerVerifier(t, settings)

	assert.Equal(t, "1", getString(inner, 1))
	assert.Equal(t, "1.1234", getString(inner, float32(1.1234)))
	assert.Equal(t, "2.45678", getString(inner, float64(2.45678)))
	assert.Equal(t, "2022-02-01T01:02:03Z", getString(inner, time.Date(2022, 02, 01, 1, 2, 3, 4, time.UTC)))
	assert.Equal(t, "9b462506-8147-4534-911a-7cdfe47ed989", getString(inner, uuid.MustParse("9b462506-8147-4534-911a-7cdfe47ed989")))
	assert.Equal(t, "marshaller", getString(inner, testMarshaller{data: "marshaller"}))
	assert.Equal(t, "stringer", getString(inner, testStringer{data: "stringer"}))
	assert.Equal(t, "1", getString(inner, int8(1)))
	assert.Equal(t, "2", getString(inner, int16(2)))
	assert.Equal(t, "3", getString(inner, int32(3)))
	assert.Equal(t, "4", getString(inner, int64(4)))
	assert.Equal(t, "5", getString(inner, uint(5)))
	assert.Equal(t, "6", getString(inner, uint8(6)))
	assert.Equal(t, "7", getString(inner, uint16(7)))
	assert.Equal(t, "8", getString(inner, uint32(8)))
	assert.Equal(t, "9", getString(inner, uint64(9)))
}

func getString(inner *innerVerifier, target interface{}) string {
	r, ok := inner.TryGetToString(target)
	if !ok {
		panic(fmt.Sprintf("could not convert to string from %v", target))
	}
	return r.Value
}
