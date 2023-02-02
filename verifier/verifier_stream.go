package verifier

import (
	"bytes"
	"io"
)

// verifyStream verifies a target of type []byte
func (v *innerVerifier) verifyStream(data []byte, extension string) {

	if len(data) == 0 {
		panic("empty data is not allowed")
	}

	//TODO: Check for ExtensionConverter

	targets := []Target{*newStreamTarget(extension, data)}

	var cleanup CleanupFunc

	v.verifyInner(nil, cleanup, targets)
}

func (v *innerVerifier) verifyReader(data io.Reader) {
	buff := new(bytes.Buffer)
	_, err := buff.ReadFrom(data)
	if err != nil {
		panic("failed to read from the io.Reader: " + err.Error())
	}

	//default bin extension for reader data,
	//or use user provided
	ext := "bin"
	if len(v.settings.extension) > 0 {
		ext = v.settings.extension
	}

	v.verifyStream(buff.Bytes(), ext)
}
