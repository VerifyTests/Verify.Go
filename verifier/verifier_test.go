package verifier_test

import (
	"github.com/heskandari/Verify.Go/verifier"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

type Test struct {
	content string
}

func TestCreateVerifier(t *testing.T) {
	testSettings := verifier.NewSettings()
	testSettings.UseDirectory("../_testdata")

	arg := "This is a test arg"
	verifier.VerifyWithSetting(t, testSettings, arg)
}

func TestNilTargets(t *testing.T) {
	testSettings := verifier.NewSettings()
	testSettings.UseDirectory("../_testdata")

	getNilTarget := func() *Test {
		return nil
	}

	arg := getNilTarget()

	verifier.VerifyWithSetting(t, testSettings, arg)
}

func TestCallbacks(t *testing.T) {
	settings := verifier.NewSettings()

	afterCalled := false
	firstCalled := false
	beforeCalled := false

	//rename to received to cause new
	_ = os.Rename("../_testdata/verifier_test.TestCallbacks.verified.txt",
		"../_testdata/verifier_test.TestCallbacks.received.txt")

	settings.UseDirectory("../_testdata")
	settings.AutoVerify()
	settings.OnFirstVerify(func(file verifier.FilePair) {
		firstCalled = true
	})
	settings.OnBeforeVerify(func() {
		beforeCalled = true
	})
	settings.OnAfterVerify(func() {
		afterCalled = true
	})

	verifier.VerifyWithSetting(t, settings, Test{})

	assert.True(t, firstCalled)
	assert.True(t, beforeCalled)
	assert.True(t, afterCalled)
}

func TestRegisteringScrubberWithExtension(t *testing.T) {
	settings := verifier.NewSettings()

	settings.UseDirectory("../_testdata")
	settings.UseExtension("js")
	settings.AddScrubberForExtension("js", func(target string) string {
		return strings.Replace(target, "Hello World!", "{msg}", 1)
	})

	verifier.VerifyWithSetting(t, settings, "alert(\"Hello World!\");")
}
