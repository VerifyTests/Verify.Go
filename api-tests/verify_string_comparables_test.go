package api_tests

import (
	"github.com/VerifyTests/Verify.Go/verifier"
	"testing"
)

type stringerStruct struct {
}

func (s stringerStruct) String() string {
	return "stringerStruct"
}

type textMarshallerStruct struct {
}

func (t textMarshallerStruct) MarshalText() (text []byte, err error) {
	return []byte("textMarshallerStruct"), nil
}

func TestStringConversion(t *testing.T) {

	testCases := map[string]interface{}{
		"nil":        nil,
		"stringer":   stringerStruct{},
		"marshaller": textMarshallerStruct{},
		"int":        int(1),
		"int8":       int8(2),
		"int32":      int32(3),
		"int64":      int64(4),
		"uint":       uint(5),
		"uint8":      uint8(6),
		"uint16":     uint16(7),
		"uint32":     uint32(8),
		"uint64":     uint64(9),
		"float32":    float32(10.0),
		"float64":    float64(11.0),
	}

	settings := verifier.NewSettings()
	settings.UseDirectory("../_testdata")

	for k, v := range testCases {

		settings.TestCase(k)
		verifier.VerifyWithSetting(t, settings, v)
	}
}
