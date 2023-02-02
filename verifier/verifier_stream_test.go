package verifier_test

import (
	"bufio"
	"bytes"
	"github.com/VerifyTests/Verify.Go/verifier"
	"os"
	"testing"
)

func TestVerifyInlineBytes(t *testing.T) {
	data := []byte("Verified Text Data Stored As Binary")
	reader := bytes.NewReader(data)

	verifier.NewVerifier(t,
		verifier.UseDirectory("../_testdata"),
	).Verify(reader)
}

func TestVerifyBytes(t *testing.T) {
	f, err := os.Open("../samples/sample.jpeg")
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(f)

	verifier.NewVerifier(t,
		verifier.UseDirectory("../_testdata"),
		verifier.UseExtension("jpg"),
	).Verify(reader)
}

func TestVerifyReader(t *testing.T) {
	f, err := os.Open("../samples/sample.jpeg")
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(f)

	verifier.NewVerifier(t,
		verifier.UseDirectory("../_testdata"),
		verifier.UseExtension("jpg"),
	).Verify(reader)
}
