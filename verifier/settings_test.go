package verifier

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifySettings_AddScrubber(t *testing.T) {

	testScrubber := func(input string) string {
		return input
	}

	s := newSettings()
	s.AddScrubber(testScrubber)

	assert.NotEmpty(t, s.instanceScrubbers)
	assert.Equal(t, 1, len(s.instanceScrubbers))
}

func TestVerifySettings_AddScrubberForExtension(t *testing.T) {

	firstTextScrubber := func(builder string) string { return "" }
	secondTextScrubber := func(builder string) string { return "" }
	jsonScrubber := func(builder string) string { return "" }

	s := newSettings()

	assert.NotNil(t, s.extensionMappedInstanceScrubbers)
	assert.Len(t, s.extensionMappedInstanceScrubbers, 0)

	s.AddScrubberForExtension("txt", firstTextScrubber)
	assert.Len(t, s.extensionMappedInstanceScrubbers["txt"], 1)

	s.AddScrubberForExtension("txt", secondTextScrubber)
	assert.Len(t, s.extensionMappedInstanceScrubbers["txt"], 2)

	s.AddScrubberForExtension("json", jsonScrubber)
	assert.Len(t, s.extensionMappedInstanceScrubbers["txt"], 2)
	assert.Len(t, s.extensionMappedInstanceScrubbers["json"], 1)
}
