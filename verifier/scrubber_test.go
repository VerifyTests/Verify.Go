package verifier

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestScrubber_RemoveLinesContaining(t *testing.T) {

	builder := strings.Builder{}
	builder.WriteString("this is\n")
	builder.WriteString("multiline\n")
	builder.WriteString("a test\n")

	scrubber := newDataScrubber(startCounter())
	scrubbed := scrubber.removeLinesContaining(builder.String(), true, "test", "multiline")

	assert.NotContains(t, scrubbed, "test")
	assert.NotContains(t, scrubbed, "multiline")
}

func TestScrubber_ScrubGuid(t *testing.T) {

	guid := "50a7501e-8af3-11ec-b80a-8c859053a807"

	scrubber := newDataScrubber(startCounter())
	scrubbed := scrubber.replaceGuids(guid)

	assert.Equal(t, scrubbed, "Guid_1")
	assert.NotContains(t, scrubbed, guid)
}

func TestScrubber_ScrubInlineGuid(t *testing.T) {

	guid := "50a7501e-8af3-11ec-b80a-8c859053a807"
	builder := strings.Builder{}
	builder.WriteString("this is\n")
	builder.WriteString(fmt.Sprintf("multiline line with uuid of %s\n", guid))
	builder.WriteString("as a test\n")

	scrubber := newDataScrubber(startCounter())
	scrubbed := scrubber.replaceGuids(builder.String())

	assert.NotContains(t, scrubbed, guid)
	assert.Contains(t, scrubbed, "multiline line with uuid of Guid_1")
}

func TestScrubber_FilterLinesWithoutLastLineEnd(t *testing.T) {

	builder := strings.Builder{}
	builder.WriteString("this is\n")
	builder.WriteString("my multiline\n")
	builder.WriteString("test\n")
	builder.WriteString("string")

	scrubber := newDataScrubber(startCounter())
	scrubbed := scrubber.replaceLines(builder.String(), testLineRemover)

	assert.NotContains(t, scrubbed, "test")
	assert.Contains(t, scrubbed, "replacement")
	assert.False(t, strings.HasSuffix(scrubbed, "\n"))
}

func TestScrubber_ScrubMachineName(t *testing.T) {

	host, _ := os.Hostname()
	target := fmt.Sprintf("machine name: %s\n", host)

	scrubber := newDataScrubber(startCounter())
	scrubbed := scrubber.ScrubMachineName(target)

	assert.NotContains(t, scrubbed, host)
}

func TestScrubber_CurrentDirectory(t *testing.T) {
	wd, _ := os.Getwd()
	builder := strings.Builder{}
	builder.WriteString(wd)

	scrubber := newDataScrubber(startCounter())

	scrubber.Apply("txt", &builder, newSettings())

	output := builder.String()
	assert.Equal(t, "{CurrentDirectory}", output)
}

func TestScrubber_FilterLinesWithLineEnd(t *testing.T) {

	builder := strings.Builder{}
	builder.WriteString("this is\n")
	builder.WriteString("a test\n")

	scrubber := newDataScrubber(startCounter())
	scrubbed := scrubber.replaceLines(builder.String(), testLineRemover)

	assert.Contains(t, scrubbed, "this is\n")
	assert.Contains(t, scrubbed, "a replacement\n")
	assert.True(t, strings.HasSuffix(scrubbed, "\n"))
}

var testLineRemover = func(value string) string {
	if strings.Contains(value, "test") {
		return strings.ReplaceAll(value, "test", "replacement")
	}
	return value
}
