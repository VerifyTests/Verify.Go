package verifier

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorReporting(t *testing.T) {
	notEqual := NotEqualFiles{
		NotEqualFilePair{
			Message: "File is not the same",
			File: FilePair{
				Name:         "NotEqualFileName",
				ReceivedPath: "../_testdata/test.verified.txt",
			},
		},
	}

	equalFiles := EqualFiles{
		FilePair{
			Name:         "EqualFileName",
			VerifiedName: "VerifiedEqual",
			ReceivedName: "ReceivedEqual",
			IsText:       true,
			Extension:    textExtension,
		},
	}

	newFiles := NewFiles{
		FilePair{
			Name:         "NewFileName",
			ReceivedName: "ReceivedNew",
			ReceivedPath: "../_testdata/test.verified.txt",
			IsText:       true,
			Extension:    textExtension,
		},
	}
	deletedFiles := []string{"Deleted1.txt"}

	builder := failingMessageBuilder{
		settings:      newSettings(),
		testCase:      "ManualTestCase",
		testName:      t.Name(),
		directory:     "../_testdata",
		notEqualFiles: notEqual,
		equalFiles:    equalFiles,
		newFiles:      newFiles,
		delete:        deletedFiles,
	}

	msg := builder.build()

	assert.NotEmpty(t, msg)
	assert.Contains(t, msg, "Deleted")
	assert.Contains(t, msg, "New:")
	assert.Contains(t, msg, "NotEqual:")
	assert.Contains(t, msg, "Delete:")
	assert.Contains(t, msg, "FileContent:")
	assert.Contains(t, msg, "ManualTestCase")
	assert.Contains(t, msg, t.Name())
}
