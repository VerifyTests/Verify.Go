package verifier

import (
	"fmt"
	"github.com/VerifyTests/Verify.Go/utils"
	"path"
)

type testingT interface {
	Name() string
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Skip(args ...interface{})
	Skipf(format string, args ...interface{})
	Cleanup(f func())
	TempDir() string
}

// EqualityResult specifies file content equality with an optional message
type EqualityResult struct {
	Equality FileEquality
	Message  string
}

// FileEquality file content equality
type FileEquality int

const (
	//FileEqual designates files with equal contents
	FileEqual FileEquality = iota
	//FileNotEqual designates files with different contents
	FileNotEqual
	//FileNew designates new files
	FileNew
)

// CompareResult specifies comparison results with an optional message
type CompareResult struct {
	IsEqual bool
	Message string
}

// InstanceScrubber a function that scrubs the input target and returns the scrubbed version
type InstanceScrubber func(target string) string

// StringComparerFunc a function that compares string content of received and verified files and returns the CompareResult
type StringComparerFunc func(received, verified string) CompareResult

// StreamComparerFunc a function that compares stream contents of received and verified files and returns the CompareResult
type StreamComparerFunc func(received, verified []byte) CompareResult

// RemoveLineFunc decides if a line from the target should be removed
type RemoveLineFunc func(string) bool

// ReplaceLineFunc replaces a line from the target the the output
type ReplaceLineFunc func(string) string

// CleanupFunc cleanup function
type CleanupFunc func()

// FileConventionFunc provides a unique file and directory for the test
type FileConventionFunc func(uniqueness string) (fileNamePrefix string, directory string)

// FileAppenderFunc returns a target
type FileAppenderFunc func() *Target

// JSONAppenderFunc returns an appender
type JSONAppenderFunc func() *toAppend

// FirstVerifyFunc a function to run on first verification
type FirstVerifyFunc func(file FilePair)

// BeforeVerifyFunc a function to run before verification
type BeforeVerifyFunc func()

// AfterVerifyFunc a function to run after verification
type AfterVerifyFunc func()

// VerifyMismatchFunc a function to run when a content mismatch happens
type VerifyMismatchFunc func(file FilePair, message string)

// VerifyDeleteFunc a function to run before deletion of a file
type VerifyDeleteFunc func(file string)

// FilePair information for a file being compared containing received and verified files.
type FilePair struct {
	Extension    string
	ReceivedPath string
	VerifiedPath string
	ReceivedName string
	VerifiedName string
	Name         string
	IsText       bool
}

// NewFiles a slice of FilePair that contains new files
type NewFiles []FilePair

// EqualFiles a slice of FilePair that contains equal files
type EqualFiles []FilePair

// NotEqualFiles a slice of NotEqualFilePair that contains non-equal files
type NotEqualFiles []NotEqualFilePair

// NotEqualFilePair designates a non-equal file information with an optional message
type NotEqualFilePair struct {
	File    FilePair
	Message string
}

type getFileNamesFunc func(extension string) FilePair

type getIndexedFileNamesFunc func(extension string, index int) FilePair

type asStringResult struct {
	Value     string
	Extension string
}

type toAppend struct {
	Name string
	Data interface{}
}

func newStringResult(value string) asStringResult {
	return asStringResult{
		Value:     value,
		Extension: "",
	}
}

func newFilePair(extension, prefix string) FilePair {

	received := fmt.Sprintf("%s.received.%s", prefix, extension)
	verified := fmt.Sprintf("%s.verified.%s", prefix, extension)

	return FilePair{
		Extension:    extension,
		Name:         path.Base(prefix),
		ReceivedPath: received,
		VerifiedPath: verified,
		ReceivedName: path.Base(received),
		VerifiedName: path.Base(verified),
		IsText:       utils.File.IsText(extension),
	}
}
