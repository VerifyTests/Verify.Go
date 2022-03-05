package verifier

import (
	"github.com/VerifyTests/Verify.Go/diff"
	"strings"
)

type verifySettings struct {
	directory                        string
	autoVerify                       bool
	diffEnabled                      bool
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

type VerifySettings interface {
	AutoVerify()
	EnableDiff()
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

func (v *verifySettings) OnVerifyDelete(fun VerifyDeleteFunc) {
	v.onVerifyDelete = fun
}

func (v *verifySettings) OnAfterVerify(fun AfterVerifyFunc) {
	v.onAfterVerify = fun
}

func (v *verifySettings) OnBeforeVerify(fun BeforeVerifyFunc) {
	v.onBeforeVerify = fun
}

func (v *verifySettings) OnFirstVerify(fun FirstVerifyFunc) {
	v.onFirstVerify = fun
}

func (v *verifySettings) OnVerifyMismatch(fun VerifyMismatchFunc) {
	v.onVerifyMismatch = fun
}

func (v *verifySettings) AutoVerify() {
	v.autoVerify = true
}

func (v *verifySettings) UniqueForArchitecture() {
	v.uniqueForArchitecture = true
}

func (v *verifySettings) UniqueForOperatingSystem() {
	v.uniqueForOperatingSystem = true
}

func (v *verifySettings) UniqueForRuntime() {
	v.uniqueForRuntime = true
}

func (v *verifySettings) OmitContentFromError() {
	v.omitContentFromError = true
}

func (v *verifySettings) EnableDiff() {
	v.diffEnabled = true
}

func (v *verifySettings) UseStrictJSON() {
	v.strictJSON = true
}

func (v *verifySettings) DontScrubGuids() {
	v.scrubGuids = false
}

func (v *verifySettings) DontScrubTimes() {
	v.scrubTimes = false
}

func (v *verifySettings) UseExtension(extension string) {
	guard.AgainstBadExtension(extension)
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

func (v *verifySettings) GetFileAppenders() []Target {
	result := make([]Target, 0)
	for _, appender := range v.fileAppender {
		stream := appender()
		if stream != nil {
			result = append(result, *stream)
		}
	}
	return result
}

func (v *verifySettings) UseStreamComparer(fun StreamComparerFunc) {
	v.streamComparer = fun
}

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

func (v *verifySettings) UseDirectory(directory string) {
	v.directory = directory
}

func (v *verifySettings) extensionOrTxt() string {
	if len(v.extension) == 0 {
		return textExtension
	}
	return v.extension
}

func (v *verifySettings) ScrubMachineName() {
	v.AddScrubber(v.scrubber.ScrubMachineName)
}

func (v *verifySettings) AddScrubber(fun InstanceScrubber) {
	v.instanceScrubbers = append(v.instanceScrubbers, fun)
}

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

func (v *verifySettings) ScrubLinesContainingAnyCase(stringToMatch ...string) {
	removeLines := func(target string) string {
		return v.scrubber.removeLinesContaining(target, true, stringToMatch...)
	}
	v.instanceScrubbers = append([]InstanceScrubber{removeLines}, v.instanceScrubbers...)
}

func (v *verifySettings) ScrubLinesContaining(stringToMatch ...string) {
	removeLines := func(target string) string {
		return v.scrubber.removeLinesContaining(target, false, stringToMatch...)
	}
	v.instanceScrubbers = append([]InstanceScrubber{removeLines}, v.instanceScrubbers...)
}

func (v *verifySettings) ScrubInlineGuids() {
	v.instanceScrubbers = append([]InstanceScrubber{v.scrubber.replaceGuids}, v.instanceScrubbers...)
}

func (v *verifySettings) ScrubLines(fun RemoveLineFunc) {
	filterLines := func(input string) string {
		return v.scrubber.filterLines(input, fun)
	}
	v.instanceScrubbers = append([]InstanceScrubber{filterLines}, v.instanceScrubbers...)
}

func (v *verifySettings) ScrubLinesWithReplace(fun ReplaceLineFunc) {
	filterLines := func(input string) string {
		return v.scrubber.replaceLines(input, fun)
	}
	v.instanceScrubbers = append([]InstanceScrubber{filterLines}, v.instanceScrubbers...)
}

func (v *verifySettings) ScrubEmptyLines() {
	isNullOrWhitespace := func(line string) bool {
		return len(line) == 0 || strings.TrimSpace(line) == ""
	}
	filterLines := func(input string) string {
		return v.scrubber.filterLines(input, isNullOrWhitespace)
	}
	v.instanceScrubbers = append([]InstanceScrubber{filterLines}, v.instanceScrubbers...)
}

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
		diffEnabled:                      false,
		autoVerify:                       false,
	}
}
