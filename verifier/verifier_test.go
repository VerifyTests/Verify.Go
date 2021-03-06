package verifier_test

import (
	"github.com/VerifyTests/Verify.Go/verifier"
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

	if !firstCalled {
		t.Fatalf("OnFirstVerify was not called")
	}

	if !beforeCalled {
		t.Fatalf("OnBeforeVerify was not called")
	}

	if !afterCalled {
		t.Fatalf("OnAfterVerify was not called")
	}
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
