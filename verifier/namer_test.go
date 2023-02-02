package verifier

import (
	"testing"
)

func TestNamerUniqueness(t *testing.T) {

	testCases := map[string]VerifyConfigure{
		"runtime": UniqueForRuntime(),
		"arch":    UniqueForArchitecture(),
		"os":      UniqueForOperatingSystem(),
		"os/arch": func(s *verifySettings) {
			s.uniqueForArchitecture = true
			s.uniqueForOperatingSystem = true
		},
		"os/runtime": func(s *verifySettings) {
			s.uniqueForOperatingSystem = true
			s.uniqueForRuntime = true
		},
		"os/runtime/arch": func(s *verifySettings) {
			s.uniqueForArchitecture = true
			s.uniqueForOperatingSystem = true
			s.uniqueForRuntime = true
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {

			setting := newSettings(t)
			test(setting)
			namer := newNamer(setting)
			result := namer.getUniqueness()
			t.Logf(result)
			if len(result) == 0 {
				t.Fatalf("Should have created a name")
			}
		})
	}
}
