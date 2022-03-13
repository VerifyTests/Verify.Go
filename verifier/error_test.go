package verifier

import (
	"strings"
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

	if len(msg) == 0 {
		t.Fatalf("Should have generated report")
	}
	if !strings.Contains(msg, "Deleted") {
		t.Fatalf("Should contain the value")
	}
	if !strings.Contains(msg, "New:") {
		t.Fatalf("Should contain the value")
	}
	if !strings.Contains(msg, "NotEqual:") {
		t.Fatalf("Should contain the value")
	}
	if !strings.Contains(msg, "Delete:") {
		t.Fatalf("Should contain the value")
	}
	if !strings.Contains(msg, "FileContent:") {
		t.Fatalf("Should contain the value")
	}
	if !strings.Contains(msg, "ManualTestCase") {
		t.Fatalf("Should contain the value")
	}
	if !strings.Contains(msg, t.Name()) {
		t.Fatalf("Should contain the value")
	}
}
