package verifier_test

import (
	"github.com/VerifyTests/Verify.Go/verifier"
	"github.com/google/uuid"
	"os"
	"strings"
	"testing"
)

type Test struct {
	content string
}

func TestCreateVerifier(t *testing.T) {
	arg := "This is a test arg"

	verifier.NewVerifier(t,
		verifier.UseDirectory("../_testdata"),
	).Verify(arg)
}

func TestNilTargets(t *testing.T) {
	getNilTarget := func() *Test {
		return nil
	}

	arg := getNilTarget()

	verifier.NewVerifier(t,
		verifier.UseDirectory("../_testdata"),
	).Verify(arg)
}

func TestCallbacks_IgnoreVerify(t *testing.T) {
	afterCalled := false
	firstCalled := false
	beforeCalled := false

	//remove to received new
	_ = os.Remove("../_testdata/verifier_core_test.Callbacks_IgnoreVerify.verified.txt")

	verifier.NewVerifier(t,
		verifier.UseDirectory("../_testdata"),
		verifier.AutoVerify(),
		verifier.OnFirstVerify(func(file verifier.FilePair) {
			firstCalled = true
		}),
		verifier.OnBeforeVerify(func() {
			beforeCalled = true
		}),
		verifier.OnAfterVerify(func() {
			afterCalled = true
		}),
	).Verify(uuid.New())

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
	verifier.NewVerifier(t,
		verifier.UseDirectory("../_testdata"),
		verifier.UseExtension("js"),
		verifier.AddScrubberForExtension("js", func(target string) string {
			return strings.Replace(target, "Hello World!", "{redacted}", 1)
		}),
	).Verify("alert(\"Hello World!\");")
}
