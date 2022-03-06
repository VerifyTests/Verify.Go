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

type EqualityResult struct {
	Equality FileEquality
	Message  string
}

type FileEquality int

const (
	FileEqual FileEquality = iota
	FileNotEqual
	FileNew
)

type CompareResult struct {
	IsEqual bool
	Message string
}

type InstanceScrubber func(target string) string

type StringComparerFunc func(received, verified string) CompareResult

type StreamComparerFunc func(received, verified []byte) CompareResult

type RemoveLineFunc func(string) bool

type ReplaceLineFunc func(string) string

type CleanupFunc func()

type FileConventionFunc func(uniqueness string) (fileNamePrefix string, directory string)

type FileAppenderFunc func() *Target

type JSONAppenderFunc func() *toAppend

type FirstVerifyFunc func(file FilePair)

type BeforeVerifyFunc func()

type AfterVerifyFunc func()

type VerifyMismatchFunc func(file FilePair, message string)

type VerifyDeleteFunc func(file string)

type FilePair struct {
	Extension    string
	ReceivedPath string
	VerifiedPath string
	ReceivedName string
	VerifiedName string
	Name         string
	IsText       bool
}

type NewFiles []FilePair
type EqualFiles []FilePair
type NotEqualFiles []NotEqualFilePair
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
