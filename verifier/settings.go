package verifier

import (
	"github.com/VerifyTests/Verify.Go/diff"
	"github.com/VerifyTests/Verify.Go/utils"
	"strings"
)

type verifySettings struct {
	directory                        string
	autoVerify                       bool
	diffDisabled                     bool
	strictJSON                       bool
	scrubGuids                       bool
	scrubTimes                       bool
	omitContentFromError             bool
	uniqueForArchitecture            bool
	uniqueForOperatingSystem         bool
	uniqueForRuntime                 bool
	instanceScrubbers                []InstanceScrubber
	fileAppender                     []FileAppenderFunc
	jsonAppender                     []JSONAppenderFunc
	extensionMappedInstanceScrubbers map[string][]InstanceScrubber
	testCase                         string
	extension                        string
	defaultExtension                 string
	ciDetected                       diff.CIDetected
	scrubber                         *dataScrubber
	counter                          *countHolder
	defaultStringComparer            StringComparerFunc
	stringComparers                  map[string]StringComparerFunc
	streamComparers                  map[string]StreamComparerFunc
	streamComparer                   StreamComparerFunc
	stringComparer                   StringComparerFunc
	onBeforeVerify                   BeforeVerifyFunc
	onAfterVerify                    AfterVerifyFunc
	onFirstVerify                    FirstVerifyFunc
	onVerifyMismatch                 VerifyMismatchFunc
	onVerifyDelete                   VerifyDeleteFunc
	t                                testingT
}

type VerifyConfigure = func(settings *verifySettings)

// OnVerifyDelete callback that is executed before a file is deleted
func OnVerifyDelete(fun VerifyDeleteFunc) VerifyConfigure {
	return func(s *verifySettings) {
		s.onVerifyDelete = fun
	}
}

// OnAfterVerify callback that is executed after the verify process
func OnAfterVerify(fun AfterVerifyFunc) VerifyConfigure {
	return func(s *verifySettings) {
		s.onAfterVerify = fun
	}
}

// OnBeforeVerify callback that is executed before the verify process
func OnBeforeVerify(fun BeforeVerifyFunc) VerifyConfigure {
	return func(s *verifySettings) {
		s.onBeforeVerify = fun
	}
}

// OnFirstVerify callback that is executed on the first verify
func OnFirstVerify(fun FirstVerifyFunc) VerifyConfigure {
	return func(s *verifySettings) {
		s.onFirstVerify = fun
	}
}

// OnVerifyMismatch callback that is executed when a mismatch happens
func OnVerifyMismatch(fun VerifyMismatchFunc) VerifyConfigure {
	return func(s *verifySettings) {
		s.onVerifyMismatch = fun
	}
}

// AutoVerify automatically accepts the received files
func AutoVerify() VerifyConfigure {
	return func(s *verifySettings) {
		s.autoVerify = true
	}
}

// UniqueForArchitecture create file names based on the runtime architecture
func UniqueForArchitecture() VerifyConfigure {
	return func(s *verifySettings) {
		s.uniqueForArchitecture = true
	}
}

// UniqueForOperatingSystem create file names based on the runtime operating system
func UniqueForOperatingSystem() VerifyConfigure {
	return func(s *verifySettings) {
		s.uniqueForOperatingSystem = true
	}
}

// UniqueForRuntime create file names based on the Go runtime versions
func UniqueForRuntime() VerifyConfigure {
	return func(s *verifySettings) {
		s.uniqueForRuntime = true
	}
}

// OmitContentFromError show the content differences when a mismatch occurs during verification
func OmitContentFromError() VerifyConfigure {
	return func(s *verifySettings) {
		s.omitContentFromError = true
	}
}

// DisableDiff enables the diff tools
func DisableDiff() VerifyConfigure {
	return func(s *verifySettings) {
		s.diffDisabled = true
	}
}

// UseStrictJSON use .json extension for the outputted files
func UseStrictJSON() VerifyConfigure {
	return func(s *verifySettings) {
		s.strictJSON = true
	}
}

// DontScrubGuids do not auto-scrub UUID values
func DontScrubGuids() VerifyConfigure {
	return func(s *verifySettings) {
		s.scrubGuids = false
	}
}

// DontScrubTimes do not auto-scrub time.Time values
func DontScrubTimes() VerifyConfigure {
	return func(s *verifySettings) {
		s.scrubTimes = false
	}
}

// UseExtension specify an extension to use for the outputted files
func UseExtension(extension string) VerifyConfigure {
	return func(s *verifySettings) {
		utils.Guard.AgainstBadExtension(extension)
		s.extension = extension
	}
}

// UseStreamComparer use the specified function for stream comparison
func UseStreamComparer(fun StreamComparerFunc) VerifyConfigure {
	return func(s *verifySettings) {
		s.streamComparer = fun
	}
}

// UseStringComparer use the specified function for string comparison
func UseStringComparer(fun StringComparerFunc) VerifyConfigure {
	return func(s *verifySettings) {
		s.stringComparer = fun
	}
}

// UseDirectory place the output files in the specified directory
func UseDirectory(directory string) VerifyConfigure {
	return func(s *verifySettings) {
		s.directory = directory
	}
}

// AddScrubber add a function to the front of the scrubber collections.
func AddScrubber(fun InstanceScrubber) VerifyConfigure {
	return func(s *verifySettings) {
		s.instanceScrubbers = append(s.instanceScrubbers, fun)
	}
}

// ScrubMachineName scrubs the machine name from the target data
func ScrubMachineName() VerifyConfigure {
	return func(s *verifySettings) {
		s.instanceScrubbers = append(s.instanceScrubbers, s.scrubber.ScrubMachineName)
	}
}

// AddScrubberForExtension adds a function for a specified extension to the front of the scrubber collection.
func AddScrubberForExtension(extension string, fun InstanceScrubber) VerifyConfigure {
	return func(s *verifySettings) {
		current, found := s.extensionMappedInstanceScrubbers[extension]
		if !found {
			list := make([]InstanceScrubber, 0)
			list = append(list, fun)
			s.extensionMappedInstanceScrubbers[extension] = list
		} else {
			s.extensionMappedInstanceScrubbers[extension] = append([]InstanceScrubber{fun}, current...)
		}
	}
}

// ScrubLinesContainingAnyCase scrubs strings that match the data in the target
func ScrubLinesContainingAnyCase(stringToMatch ...string) VerifyConfigure {
	return func(s *verifySettings) {
		removeLines := func(target string) string {
			return s.scrubber.removeLinesContaining(target, true, stringToMatch...)
		}
		s.instanceScrubbers = append([]InstanceScrubber{removeLines}, s.instanceScrubbers...)
	}
}

// ScrubLinesContaining scrubs the line containing specified strings
func ScrubLinesContaining(stringToMatch ...string) VerifyConfigure {
	return func(s *verifySettings) {
		removeLines := func(target string) string {
			return s.scrubber.removeLinesContaining(target, false, stringToMatch...)
		}
		s.instanceScrubbers = append([]InstanceScrubber{removeLines}, s.instanceScrubbers...)
	}
}

// ScrubInlineGuids scrubs inline UUID values with string types
func ScrubInlineGuids() VerifyConfigure {
	return func(s *verifySettings) {
		s.instanceScrubbers = append([]InstanceScrubber{s.scrubber.replaceGuids}, s.instanceScrubbers...)
	}
}

// ScrubInlineTime scrubs inline Time values with string types
func ScrubInlineTime(format string) VerifyConfigure {
	return func(s *verifySettings) {
		s.instanceScrubbers = append([]InstanceScrubber{
			func(target string) string {
				return s.scrubber.replaceTime(format, target)
			},
		}, s.instanceScrubbers...)
	}
}

// ScrubLines scrub target lines with the provided function
func ScrubLines(fun RemoveLineFunc) VerifyConfigure {
	return func(s *verifySettings) {
		filterLines := func(input string) string {
			return s.scrubber.filterLines(input, fun)
		}
		s.instanceScrubbers = append([]InstanceScrubber{filterLines}, s.instanceScrubbers...)
	}
}

// ScrubLinesWithReplace scrubs target lines and replace with the value provided by the function
func ScrubLinesWithReplace(fun ReplaceLineFunc) VerifyConfigure {
	return func(s *verifySettings) {
		filterLines := func(input string) string {
			return s.scrubber.replaceLines(input, fun)
		}
		s.instanceScrubbers = append([]InstanceScrubber{filterLines}, s.instanceScrubbers...)
	}
}

// ScrubEmptyLines scrubs all the empty lines from the target
func ScrubEmptyLines() VerifyConfigure {
	return func(s *verifySettings) {
		isNullOrWhitespace := func(line string) bool {
			return len(line) == 0 || strings.TrimSpace(line) == ""
		}
		filterLines := func(input string) string {
			return s.scrubber.filterLines(input, isNullOrWhitespace)
		}
		s.instanceScrubbers = append([]InstanceScrubber{filterLines}, s.instanceScrubbers...)
	}
}

// TestCase specify a case name for the test.
func TestCase(name string) VerifyConfigure {
	return func(s *verifySettings) {
		s.testCase = name
	}
}

func (s *verifySettings) tryGetStringComparer(extension string) (StringComparerFunc, bool) {
	comp, ok := s.stringComparers[extension]
	if ok {
		return comp, true
	}

	if s.defaultStringComparer != nil {
		return s.defaultStringComparer, true
	}

	return nil, false
}

func (s *verifySettings) extensionOrTxt() string {
	if len(s.extension) == 0 {
		return textExtension
	}
	return s.extension
}

func (s *verifySettings) runOnFirstVerify(file FilePair) {
	if s.onFirstVerify != nil {
		s.onFirstVerify(file)
	}
}

func (s *verifySettings) runAfterVerify() {
	if s.onAfterVerify != nil {
		s.onAfterVerify()
	}
}

func (s *verifySettings) runBeforeVerify() {
	if s.onBeforeVerify != nil {
		s.onBeforeVerify()
	}
}

func (s *verifySettings) runOnVerifyMismatch(file FilePair, message string) {
	if s.onVerifyMismatch != nil {
		s.onVerifyMismatch(file, message)
	}
}

func (s *verifySettings) runOnVerifyDelete(file string) {
	if s.onVerifyMismatch != nil {
		s.onVerifyDelete(file)
	}
}

func (s *verifySettings) getJSONAppenders() []toAppend {
	result := make([]toAppend, 0)
	for _, appender := range s.jsonAppender {
		data := appender()
		if data != nil {
			result = append(result, *data)
		}
	}
	return result
}

func (s *verifySettings) getFileAppenders() []Target {
	result := make([]Target, 0)
	for _, appender := range s.fileAppender {
		stream := appender()
		if stream != nil {
			result = append(result, *stream)
		}
	}
	return result
}

func newSettings(t testingT) *verifySettings {
	return &verifySettings{
		scrubGuids:                       true,
		scrubTimes:                       true,
		extensionMappedInstanceScrubbers: make(map[string][]InstanceScrubber),
		instanceScrubbers:                make([]InstanceScrubber, 0),
		fileAppender:                     make([]FileAppenderFunc, 0),
		jsonAppender:                     make([]JSONAppenderFunc, 0),
		streamComparers:                  make(map[string]StreamComparerFunc),
		stringComparers:                  make(map[string]StringComparerFunc),
		scrubber:                         newDataScrubber(startCounter()),
		ciDetected:                       diff.CheckCI(),
		diffDisabled:                     diff.CheckDisabled(),
		autoVerify:                       false,
		t:                                t,
	}
}
