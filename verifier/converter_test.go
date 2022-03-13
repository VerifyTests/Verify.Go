package verifier

import (
	"fmt"
	"github.com/google/uuid"
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

	if getString(inner, 1) != "1" {
		t.Fatalf("should correctly convert to string")
	}
	if getString(inner, float32(1.1234)) != "1.1234" {
		t.Fatalf("should correctly convert to string")
	}
	if getString(inner, float64(2.45678)) != "2.45678" {
		t.Fatalf("should correctly convert to string")
	}
	if getString(inner, time.Date(2022, 02, 01, 1, 2, 3, 4, time.UTC)) != "2022-02-01T01:02:03Z" {
		t.Fatalf("should correctly convert to string")
	}
	if getString(inner, uuid.MustParse("9b462506-8147-4534-911a-7cdfe47ed989")) != "9b462506-8147-4534-911a-7cdfe47ed989" {
		t.Fatalf("should correctly convert to string")
	}
	if getString(inner, testMarshaller{data: "marshaller"}) != "marshaller" {
		t.Fatalf("should correctly convert to string")
	}
	if getString(inner, testStringer{data: "stringer"}) != "stringer" {
		t.Fatalf("should correctly convert to string")
	}
	if getString(inner, int8(1)) != "1" {
		t.Fatalf("should correctly convert to string")
	}
	if getString(inner, int16(2)) != "2" {
		t.Fatalf("should correctly convert to string")
	}
	if getString(inner, int32(3)) != "3" {
		t.Fatalf("should correctly convert to string")
	}
	if getString(inner, int64(4)) != "4" {
		t.Fatalf("should correctly convert to string")
	}
	if getString(inner, uint(5)) != "5" {
		t.Fatalf("should correctly convert to string")
	}
	if getString(inner, uint8(6)) != "6" {
		t.Fatalf("should correctly convert to string")
	}
	if getString(inner, uint16(7)) != "7" {
		t.Fatalf("should correctly convert to string")
	}
	if getString(inner, uint32(8)) != "8" {
		t.Fatalf("should correctly convert to string")
	}
	if getString(inner, uint64(9)) != "9" {
		t.Fatalf("should correctly convert to string")
	}
}

func getString(inner *innerVerifier, target interface{}) string {
	r, ok := inner.tryGetToString(target)
	if !ok {
		panic(fmt.Sprintf("could not convert to string from %v", target))
	}
	return r.Value
}
