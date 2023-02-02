package verifier

import (
	"encoding"
	"fmt"
	"github.com/VerifyTests/Verify.Go/utils"
	"github.com/google/uuid"
	"github.com/modern-go/reflect2"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	textExtension  = "txt"
	jsonExtension  = "json"
	emptyExtension = ""
)

var emptyTargets = make([]Target, 0)
var prefixList = make([]string, 0)
var locker = &sync.Mutex{}

type innerVerifier struct {
	testing             testingT
	settings            *verifySettings
	scrubber            *dataScrubber
	counter             *countHolder
	getFileNames        getFileNamesFunc
	getIndexedFileNames getIndexedFileNamesFunc
	outputDirectory     string
	verifiedFiles       []string
	receivedFiles       []string
}

func newInnerVerifier(t testingT, settings *verifySettings) *innerVerifier {
	uniqueness := newNamer(settings).getUniqueness()
	fileName, directory := defaultFileConvention(t, settings, uniqueness)
	sourceFileDirectory := filepath.Dir(fileName)

	if len(directory) == 0 {
		directory = sourceFileDirectory
	} else {
		directory = path.Join(directory, sourceFileDirectory)
		if !utils.File.Exists(directory) {
			err := utils.File.CreateDirectory(directory)
			if err != nil {
				log.Fatalf("Failed to delete %s", directory)
			}
		}
	}

	filePathPrefix := path.Join(directory, fileName)
	validatePrefix(filePathPrefix)

	pattern := fmt.Sprintf("%s.*.*", fileName)
	files, _ := fileMatch(directory, pattern)

	verifier := &innerVerifier{
		scrubber:            settings.scrubber,
		testing:             t,
		outputDirectory:     directory,
		settings:            settings,
		verifiedFiles:       findMatchingFiles(files, fileName, ".verified"),
		receivedFiles:       findMatchingFiles(files, fileName, ".received"),
		getFileNames:        getFileNamePair(filePathPrefix),
		getIndexedFileNames: getIndexFileNamePair(filePathPrefix),
	}

	for _, f := range verifier.receivedFiles {
		utils.File.Delete(f)
	}

	settings.runBeforeVerify()

	return verifier
}

func (v *innerVerifier) verifyInner(target interface{}, cleanup CleanupFunc, targets []Target) {
	if builder, extension, found := v.tryGetTargetBuilder(target); found {
		v.scrubber.Apply(extension, builder, v.settings)

		received := builder.String()
		stringTarget := newStringTarget(extension, received)
		targets = append([]Target{*stringTarget}, targets...)
	}

	targets = append(targets, v.settings.getFileAppenders()...)

	engine := newEngine(v.testing, v.outputDirectory, v.settings, v.verifiedFiles, v.getFileNames, v.getIndexedFileNames)

	engine.handleResults(targets)

	if cleanup != nil {
		cleanup()
	}

	engine.throwIfRequired()
}

// verifyStream verifies a target of type []byte
func (v *innerVerifier) verifyStream(target []byte) {
	panic("not implemented")
}

func (v *innerVerifier) tryGetTargetBuilder(target interface{}) (builder *strings.Builder, extension string, found bool) {
	appenders := v.settings.getJSONAppenders()
	hasAppends := len(appenders) > 0

	if target == nil {
		if !hasAppends {
			builder = nil
			extension = emptyExtension
			found = false
			return
		}

		extension = textExtension
		if v.settings.strictJSON {
			extension = jsonExtension
		}

		builder = asJSON(target, appenders, v.settings)
		found = true
		return
	}

	if !hasAppends {
		if stringTarget, ok := target.(string); ok {
			b := strings.Builder{}
			b.WriteString(fixNewlines(stringTarget))
			extension = v.settings.extensionOrTxt()
			found = true
			builder = &b
			return
		}
	}

	extension = textExtension
	if v.settings.strictJSON {
		extension = jsonExtension
	}

	builder = asJSON(target, appenders, v.settings)
	found = true
	return
}

func (v *innerVerifier) tryGetToString(target interface{}) (asStringResult, bool) {
	//TODO: implement SimpleName like conversion
	result, converted := v.tryGetPrimitiveToString(target)
	if converted {
		return result, true
	}

	typ := reflect2.TypeOf(target)

	if typ.AssignableTo(stringBuilderType) {
		return stringBuilderToString(target.(strings.Builder))
	}

	if typ.Implements(stringerType) {
		return stringerToString(target.(fmt.Stringer)), true
	}

	if typ.Implements(textMarshalerType) {
		return textMarshallerToString(target.(encoding.TextMarshaler))
	}

	return asStringResult{}, false
}

func (v *innerVerifier) tryGetPrimitiveToString(target interface{}) (asStringResult, bool) {
	switch v := target.(type) {
	case string:
		return stringToString(v), true
	case int:
		return intToString(v), true
	case int8:
		return int8ToString(v), true
	case int16:
		return int16ToString(v), true
	case int32:
		return int32ToString(v), true
	case int64:
		return int64ToString(v), true
	case uint:
		return uIntToString(v), true
	case uint8:
		return uInt8ToString(v), true
	case uint16:
		return uInt16ToString(v), true
	case uint32:
		return uInt32ToString(v), true
	case uint64:
		return uInt64ToString(v), true
	case bool:
		return boolToString(v), true
	case float64:
		return float64ToString(v), true
	case float32:
		return float32ToString(v), true
	case time.Time:
		return timeToString(v), true
	case uuid.UUID:
		return uUIDToString(v), true
	}

	return asStringResult{}, false
}

func defaultFileConvention(t testingT, settings *verifySettings, uniqueness string) (string, string) {
	testName, sourceFile, _ := testCallerInfo()
	testNameParts := strings.Split(t.Name(), "/")
	name := getTestCaseName(testNameParts, testName, settings.testCase)
	sourceFileWithoutExt := utils.File.GetFileNameWithoutExtension(sourceFile)

	if len(settings.directory) == 0 {
		settings.directory = filepath.Dir(sourceFile)
	}

	if len(uniqueness) > 0 {
		return fmt.Sprintf("%s.%s.%s", sourceFileWithoutExt, name, uniqueness), settings.directory
	}

	return fmt.Sprintf("%s.%s", sourceFileWithoutExt, name), settings.directory
}

func getFileNamePair(filePathPrefix string) getFileNamesFunc {
	return func(extension string) FilePair {
		return newFilePair(extension, filePathPrefix)
	}
}

func getIndexFileNamePair(filePathPrefix string) getIndexedFileNamesFunc {
	return func(extension string, index int) FilePair {
		return newFilePair(extension, fmt.Sprintf("%s.%02d", filePathPrefix, index))
	}
}

func getTestCaseName(parts []string, testName, caseName string) string {
	test := removeTestPrefix(testName)

	//test case is not specified, but can be
	//determined, based on the table test-case name
	if len(parts) == 2 && len(caseName) == 0 {
		caseName = parts[1]
	}

	if len(test) == 0 && len(caseName) > 0 {
		return caseName
	}

	if len(caseName) == 0 && len(test) > 0 {
		return test
	}

	if len(test) > 0 && len(caseName) > 0 {
		return fmt.Sprintf("%s.%s", test, caseName)
	}

	panic("Test name can't be determined. Provide the name via settings.TestCase method.")
}

func removeTestPrefix(testName string) string {
	if strings.HasPrefix(testName, "Test") { // remove 'Test' convention prefix from test name
		var testName = testName[4:]
		return testName
	}
	return testName
}

func testCallerInfo() (functionName string, filePath string, line int) {

	var pc uintptr
	var ok bool
	var file string
	var name string

	for i := 0; ; i++ {
		pc, file, line, ok = runtime.Caller(i)
		if !ok {
			// The breaks below failed to terminate the loop, and we ran off the
			// end of the call stack.
			break
		}

		// This is a huge edge case, but it will panic if this is the case
		if file == "<autogenerated>" {
			break
		}

		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}
		name = f.Name()

		// testing.tRunner is the standard library function that calls
		// tests. Subtests are called directly by tRunner, without going through
		// the Test/Benchmark/Example function that contains the t.Run calls, so
		// with subtests we should break when we hit tRunner, without adding it
		// to the list of callers.
		if name == "testing.tRunner" {
			break
		}

		segments := strings.Split(name, ".")

		//test names are usually the last part of the
		//fully qualified test name and start with a set of
		//well-known prefixes
		if len(segments)-1 >= 0 {
			potentialTestName := segments[len(segments)-1]
			if isTest(potentialTestName, "Test") ||
				isTest(potentialTestName, "Benchmark") ||
				isTest(potentialTestName, "Example") {
				functionName = potentialTestName
				filePath = file
				return
			}
		}

		//Table test name have an inner function where
		//the actual test name is not the last part of the
		//full method name, but prior to last
		if len(segments)-2 >= 0 {
			tableTestName := segments[len(segments)-2]
			if isTest(tableTestName, "Test") ||
				isTest(tableTestName, "Benchmark") ||
				isTest(tableTestName, "Example") {
				functionName = tableTestName
				filePath = file
				return
			}
		}
	}
	return
}

func isTest(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) { // "Test" is ok
		return true
	}
	r, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(r)
}

func fileMatch(root, pattern string) ([]string, error) {
	matches := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			//NOTE: Replace backslash with slash to get paths compatible with the builtin path.go functions
			matches = append(matches, strings.ReplaceAll(path, "\\", "/"))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return matches, nil
}

func findMatchingFiles(files []string, fileNamePrefix string, suffix string) []string {
	matches := make([]string, 0)
	for _, f := range files {
		name := utils.File.GetFileNameWithoutExtension(f)
		if !strings.HasPrefix(name, fileNamePrefix) {
			continue
		}

		if !strings.HasSuffix(name, suffix) {
			continue
		}

		prefixRemoved := name[len(fileNamePrefix):]
		if prefixRemoved == suffix {
			matches = append(matches, f)
			continue
		}

		numberPart := prefixRemoved[1 : len(prefixRemoved)-len(suffix)-1]
		if _, err := strconv.Atoi(numberPart); err == nil {
			matches = append(matches, f)
		}
	}
	return matches
}

func validatePrefix(prefix string) {
	locker.Lock()
	defer locker.Unlock()

	for _, f := range prefixList {
		if f == prefix {
			panic(fmt.Sprintf("The prefix has already been used: %s.\n"+
				"This is mostly caused by a conflicting combination of "+
				"`VerifySettings.UseDirectory()`, `VerifySettings.TestCase(), and `VerifySettings.TestName()`.\n"+
				"If that's not the case, and having multiple identical prefixes is acceptable, then call `VerifierSettings.DisableRequireUniquePrefix()` "+
				"to disable this uniqueness validation.", prefix))
		}
	}
	prefixList = append(prefixList, prefix)
}
