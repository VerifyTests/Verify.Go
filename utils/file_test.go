package utils

import (
	"testing"
)

func TestGetFileExtension(t *testing.T) {
	if File.GetFileExtension("textfile.txt") != "txt" {
		t.Fatalf("should get the right file extension")
	}
	if File.GetFileExtension("textfile.json") != "json" {
		t.Fatalf("should get the right file extension")
	}
	if File.GetFileExtension("txt") != "txt" {
		t.Fatalf("should get the right file extension")
	}
}

func TestFileReading(t *testing.T) {
	content := File.ReadFile("../_testdata/verifier_test.NilTargets.verified.txt")
	if len(content) == 0 {
		t.Fatalf("should read file content")
	}
}
