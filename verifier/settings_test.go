package verifier

import (
	"testing"
)

func TestVerifySettings_AddScrubber(t *testing.T) {

	testScrubber := func(input string) string {
		return input
	}

	s := newSettings()
	s.AddScrubber(testScrubber)

	if len(s.instanceScrubbers) != 1 {
		t.Fatalf("InstanceScrubbers should have 1 item")
	}
}

func TestVerifySettings_AddScrubberForExtension(t *testing.T) {

	firstTextScrubber := func(builder string) string { return "" }
	secondTextScrubber := func(builder string) string { return "" }
	jsonScrubber := func(builder string) string { return "" }

	s := newSettings()

	if len(s.extensionMappedInstanceScrubbers) != 0 {
		t.Fatalf("mapped instance scrubbers should be empty")
	}

	s.AddScrubberForExtension("txt", firstTextScrubber)
	if len(s.extensionMappedInstanceScrubbers["txt"]) != 1 {
		t.Fatalf("extension scrubber for 'txt' should have 1 instance")
	}

	s.AddScrubberForExtension("txt", secondTextScrubber)
	if len(s.extensionMappedInstanceScrubbers["txt"]) != 2 {
		t.Fatalf("extension scrubber for 'txt' should have 2 instances")
	}

	s.AddScrubberForExtension("json", jsonScrubber)
	if len(s.extensionMappedInstanceScrubbers["txt"]) != 2 {
		t.Fatalf("extension scrubber for 'txt' should have 2 instances")
	}
	if len(s.extensionMappedInstanceScrubbers["json"]) != 1 {
		t.Fatalf("extension scrubber for 'json' should have 1 instance")
	}
}
