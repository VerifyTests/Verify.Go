package verifier

import (
	"fmt"
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

	if strings.Contains(scrubbed, "test") {
		t.Fatalf("Should have the value")
	}
	if strings.Contains(scrubbed, "multiline") {
		t.Fatalf("Should have the value")
	}
}

func TestScrubber_ScrubGuid(t *testing.T) {

	guid := "50a7501e-8af3-11ec-b80a-8c859053a807"

	scrubber := newDataScrubber(startCounter())
	scrubbed := scrubber.replaceGuids(guid)

	if !strings.Contains(scrubbed, "Guid_1") {
		t.Fatalf("Should have the value")
	}
	if strings.Contains(scrubbed, guid) {
		t.Fatalf("Should not contain the guid value")
	}
}

func TestScrubber_ScrubInlineGuid(t *testing.T) {

	guid := "50a7501e-8af3-11ec-b80a-8c859053a807"
	builder := strings.Builder{}
	builder.WriteString("this is\n")
	builder.WriteString(fmt.Sprintf("multiline line with uuid of %s\n", guid))
	builder.WriteString("as a test\n")

	scrubber := newDataScrubber(startCounter())
	scrubbed := scrubber.replaceGuids(builder.String())

	if strings.Contains(scrubbed, guid) {
		t.Fatalf("Should not contain the guid value")
	}
	if !strings.Contains(scrubbed, "multiline line with uuid of Guid_1") {
		t.Fatalf("Should contain the line with scrubbed guid")
	}
}

func TestScrubber_FilterLinesWithoutLastLineEnd(t *testing.T) {

	builder := strings.Builder{}
	builder.WriteString("this is\n")
	builder.WriteString("my multiline\n")
	builder.WriteString("test\n")
	builder.WriteString("string")

	scrubber := newDataScrubber(startCounter())
	scrubbed := scrubber.replaceLines(builder.String(), testLineRemover)

	if strings.Contains(scrubbed, "test") {
		t.Fatalf("Should contain the value")
	}
	if !strings.Contains(scrubbed, "replacement") {
		t.Fatalf("Should contain the value")
	}
	if strings.HasSuffix(scrubbed, "\n") {
		t.Fatalf("Should not contain the ending new line")
	}
}

func TestScrubber_ScrubMachineName(t *testing.T) {

	host, _ := os.Hostname()
	target := fmt.Sprintf("machine name: %s\n", host)

	scrubber := newDataScrubber(startCounter())
	scrubbed := scrubber.ScrubMachineName(target)

	if strings.Contains(scrubbed, host) {
		t.Fatalf("Should not contain the initial value")
	}
}

func TestScrubber_CurrentDirectory(t *testing.T) {
	wd, _ := os.Getwd()
	exe, _ := os.Executable()
	tmp := os.TempDir()
	cache, _ := os.UserCacheDir()

	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("%s\n%s\n%s\n%s", exe, wd, tmp, cache))

	scrubber := newDataScrubber(startCounter())

	scrubber.Apply("txt", &builder, newSettings(t))

	output := builder.String()

	if !strings.Contains(output, "{ExeDir}\n{CurrentDirectory}\n{TempDir}\n{CacheDir}") {
		t.Fatalf("Should scrub well known directories")
	}
}

func TestScrubber_FilterLinesWithLineEnd(t *testing.T) {

	builder := strings.Builder{}
	builder.WriteString("this is\n")
	builder.WriteString("a test\n")

	scrubber := newDataScrubber(startCounter())
	scrubbed := scrubber.replaceLines(builder.String(), testLineRemover)

	if !strings.Contains(scrubbed, "this is\n") {
		t.Fatalf("should contain the value")
	}
	if !strings.Contains(scrubbed, "a replacement\n") {
		t.Fatalf("should contain the value")
	}
	if !strings.HasSuffix(scrubbed, "\n") {
		t.Fatalf("should contain the last new line")
	}
}

var testLineRemover = func(value string) string {
	if strings.Contains(value, "test") {
		return strings.ReplaceAll(value, "test", "replacement")
	}
	return value
}
