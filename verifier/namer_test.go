package verifier

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNamerUniqueness(t *testing.T) {

	testCases := map[string]func(VerifySettings){
		"runtime": func(s VerifySettings) {
			s.UniqueForRuntime()
		},
		"arch": func(s VerifySettings) {
			s.UniqueForArchitecture()
		},
		"os": func(s VerifySettings) {
			s.UniqueForOperatingSystem()
		},
		"os/arch": func(s VerifySettings) {
			s.UniqueForOperatingSystem()
			s.UniqueForArchitecture()
		},
		"os/runtime": func(s VerifySettings) {
			s.UniqueForOperatingSystem()
			s.UniqueForRuntime()
		},
		"os/runtime/arch": func(s VerifySettings) {
			s.UniqueForOperatingSystem()
			s.UniqueForRuntime()
			s.UniqueForArchitecture()
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {

			setting := newSettings()
			test(setting)
			namer := newNamer(setting)
			result := namer.getUniqueness()
			t.Logf(result)
			assert.NotEmpty(t, result)
		})
	}
}
