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
}

// VerifySettings provides customization for the Verify process.
type VerifySettings interface {
	AutoVerify()
	DisableDiff()
	UseStrictJSON()
	DontScrubGuids()
	DontScrubTimes()
	UniqueForArchitecture()
	UniqueForOperatingSystem()
	UniqueForRuntime()
	UseDirectory(directory string)
	UseExtension(extension string)
	ScrubMachineName()
	AddScrubber(fun InstanceScrubber)
	AddScrubberForExtension(extension string, fun InstanceScrubber)
	ScrubLinesContainingAnyCase(stringToMatch ...string)
	ScrubLinesContaining(stringToMatch ...string)
	ScrubInlineGuids()
	ScrubLines(fun RemoveLineFunc)
	ScrubLinesWithReplace(fun ReplaceLineFunc)
	ScrubEmptyLines()
	UseStreamComparer(fun StreamComparerFunc)
	UseStringComparer(fun StringComparerFunc)
	OnAfterVerify(fun AfterVerifyFunc)
	OnBeforeVerify(fun BeforeVerifyFunc)
	OnFirstVerify(fun FirstVerifyFunc)
	OnVerifyMismatch(fun VerifyMismatchFunc)
	OnVerifyDelete(fun VerifyDeleteFunc)
	OmitContentFromError()
	TestCase(name string)
}

// OnVerifyDelete callback that is executed before a file is deleted
func (v *verifySettings) OnVerifyDelete(fun VerifyDeleteFunc) {
	v.onVerifyDelete = fun
}

// OnAfterVerify callback that is executed after the verify process
func (v *verifySettings) OnAfterVerify(fun AfterVerifyFunc) {
	v.onAfterVerify = fun
}

// OnBeforeVerify callback that is executed before the verify process
func (v *verifySettings) OnBeforeVerify(fun BeforeVerifyFunc) {
	v.onBeforeVerify = fun
}

// OnFirstVerify callback that is executed on the first verify
func (v *verifySettings) OnFirstVerify(fun FirstVerifyFunc) {
	v.onFirstVerify = fun
}

// OnVerifyMismatch callback that is executed when a mismatch happens
func (v *verifySettings) OnVerifyMismatch(fun VerifyMismatchFunc) {
	v.onVerifyMismatch = fun
}

// AutoVerify automatically accepts the received files
func (v *verifySettings) AutoVerify() {
	v.autoVerify = true
}

// UniqueForArchitecture create file names based on the runtime architecture
func (v *verifySettings) UniqueForArchitecture() {
	v.uniqueForArchitecture = true
}

// UniqueForOperatingSystem create file names based on the runtime operating system
func (v *verifySettings) UniqueForOperatingSystem() {
	v.uniqueForOperatingSystem = true
}

// UniqueForRuntime create file names based on the Go runtime versions
func (v *verifySettings) UniqueForRuntime() {
	v.uniqueForRuntime = true
}

// OmitContentFromError show the content differences when a mismatch occurs during verification
func (v *verifySettings) OmitContentFromError() {
	v.omitContentFromError = true
}

// DisableDiff enables the diff tools
func (v *verifySettings) DisableDiff() {
	v.diffDisabled = false
}

// UseStrictJSON use .json extension for the outputted files
func (v *verifySettings) UseStrictJSON() {
	v.strictJSON = true
}

// DontScrubGuids do not auto-scrub UUID values
func (v *verifySettings) DontScrubGuids() {
	v.scrubGuids = false
}

// DontScrubTimes do not auto-scrub time.Time values
func (v *verifySettings) DontScrubTimes() {
	v.scrubTimes = false
}

// UseExtension specify an extension to use for the outputted files
func (v *verifySettings) UseExtension(extension string) {
	utils.Guard.AgainstBadExtension(extension)
	v.extension = extension
}

func (v *verifySettings) getJSONAppenders() []toAppend {
	result := make([]toAppend, 0)
	for _, appender := range v.jsonAppender {
		data := appender()
		if data != nil {
			result = append(result, *data)
		}
	}
	return result
}

func (v *verifySettings) getFileAppenders() []Target {
	result := make([]Target, 0)
	for _, appender := range v.fileAppender {
		stream := appender()
		if stream != nil {
			result = append(result, *stream)
		}
	}
	return result
}

// UseStreamComparer use the specified function for stream comparison
func (v *verifySettings) UseStreamComparer(fun StreamComparerFunc) {
	v.streamComparer = fun
}

// UseStringComparer use the specified function for string comparison
func (v *verifySettings) UseStringComparer(fun StringComparerFunc) {
	v.stringComparer = fun
}

func (v *verifySettings) tryGetStringComparer(extension string) (StringComparerFunc, bool) {
	comp, ok := v.stringComparers[extension]
	if ok {
		return comp, true
	}

	if v.defaultStringComparer != nil {
		return v.defaultStringComparer, true
	}

	return nil, false
}

// UseDirectory place the output files in the specified directory
func (v *verifySettings) UseDirectory(directory string) {
	v.directory = directory
}

func (v *verifySettings) extensionOrTxt() string {
	if len(v.extension) == 0 {
		return textExtension
	}
	return v.extension
}

// ScrubMachineName scrubs the machine name from the target data
func (v *verifySettings) ScrubMachineName() {
	v.AddScrubber(v.scrubber.ScrubMachineName)
}

// AddScrubber add a function to the front of the scrubber collections.
func (v *verifySettings) AddScrubber(fun InstanceScrubber) {
	v.instanceScrubbers = append(v.instanceScrubbers, fun)
}

// AddScrubberForExtension adds a function for a specified extension to the front of the scrubber collection.
func (v *verifySettings) AddScrubberForExtension(extension string, fun InstanceScrubber) {
	current, found := v.extensionMappedInstanceScrubbers[extension]
	if !found {
		list := make([]InstanceScrubber, 0)
		list = append(list, fun)
		v.extensionMappedInstanceScrubbers[extension] = list
	} else {
		v.extensionMappedInstanceScrubbers[extension] = append([]InstanceScrubber{fun}, current...)
	}
}

// ScrubLinesContainingAnyCase scrubs strings that match the data in the target
func (v *verifySettings) ScrubLinesContainingAnyCase(stringToMatch ...string) {
	removeLines := func(target string) string {
		return v.scrubber.removeLinesContaining(target, true, stringToMatch...)
	}
	v.instanceScrubbers = append([]InstanceScrubber{removeLines}, v.instanceScrubbers...)
}

// ScrubLinesContaining scrubs the line containing specified strings
func (v *verifySettings) ScrubLinesContaining(stringToMatch ...string) {
	removeLines := func(target string) string {
		return v.scrubber.removeLinesContaining(target, false, stringToMatch...)
	}
	v.instanceScrubbers = append([]InstanceScrubber{removeLines}, v.instanceScrubbers...)
}

// ScrubInlineGuids scrubs inline UUID values with string types
func (v *verifySettings) ScrubInlineGuids() {
	v.instanceScrubbers = append([]InstanceScrubber{v.scrubber.replaceGuids}, v.instanceScrubbers...)
}

// ScrubLines scrub target lines with the provided function
func (v *verifySettings) ScrubLines(fun RemoveLineFunc) {
	filterLines := func(input string) string {
		return v.scrubber.filterLines(input, fun)
	}
	v.instanceScrubbers = append([]InstanceScrubber{filterLines}, v.instanceScrubbers...)
}

// ScrubLinesWithReplace scrubs target lines and replace with the value provided by the function
func (v *verifySettings) ScrubLinesWithReplace(fun ReplaceLineFunc) {
	filterLines := func(input string) string {
		return v.scrubber.replaceLines(input, fun)
	}
	v.instanceScrubbers = append([]InstanceScrubber{filterLines}, v.instanceScrubbers...)
}

// ScrubEmptyLines scrubs all the empty lines from the target
func (v *verifySettings) ScrubEmptyLines() {
	isNullOrWhitespace := func(line string) bool {
		return len(line) == 0 || strings.TrimSpace(line) == ""
	}
	filterLines := func(input string) string {
		return v.scrubber.filterLines(input, isNullOrWhitespace)
	}
	v.instanceScrubbers = append([]InstanceScrubber{filterLines}, v.instanceScrubbers...)
}

// TestCase specify a case name for the test.
func (v *verifySettings) TestCase(name string) {
	v.testCase = name
}

func (v *verifySettings) runOnFirstVerify(file FilePair) {
	if v.onFirstVerify != nil {
		v.onFirstVerify(file)
	}
}

func (v *verifySettings) runAfterVerify() {
	if v.onAfterVerify != nil {
		v.onAfterVerify()
	}
}

func (v *verifySettings) runBeforeVerify() {
	if v.onBeforeVerify != nil {
		v.onBeforeVerify()
	}
}

func (v *verifySettings) runOnVerifyMismatch(file FilePair, message string) {
	if v.onVerifyMismatch != nil {
		v.onVerifyMismatch(file, message)
	}
}

func (v *verifySettings) runOnVerifyDelete(file string) {
	if v.onVerifyMismatch != nil {
		v.onVerifyDelete(file)
	}
}

// NewSettings returns a new default instance of the VerifySettings
func NewSettings() VerifySettings {
	return newSettings()
}

func newSettings() *verifySettings {
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
	}
}