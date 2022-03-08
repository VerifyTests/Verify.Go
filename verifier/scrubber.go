package verifier

import (
	"bufio"
	"fmt"
	"github.com/VerifyTests/Verify.Go/utils"
	"github.com/google/uuid"
	"os"
	"regexp"
	"strings"
	"time"
)

var currentDirectoryReplacements = make(map[string]struct{})
var tempDirectoryReplacements = make(map[string]struct{})
var dirSeparator = string(os.PathSeparator)
var guidPattern = "[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"
var reg = regexp.MustCompile(guidPattern)

type dataScrubber struct {
	counter *countHolder
}

func newDataScrubber(counter *countHolder) *dataScrubber {
	return &dataScrubber{
		counter: counter,
	}
}

func init() {
	if home, err := os.UserHomeDir(); err == nil {
		addDirectory(home)
	}
	if config, err := os.UserConfigDir(); err == nil {
		addDirectory(config)
	}
	if cache, err := os.UserCacheDir(); err == nil {
		addDirectory(cache)
	}
	if exe, err := os.Executable(); err == nil {
		addDirectory(exe)
	}
	if cur, err := os.Getwd(); err == nil {
		addDirectory(cur)
	}

	addTempDirectory(os.TempDir())
}

func addTempDirectory(value string) {
	add(tempDirectoryReplacements, value)
}

func addDirectory(value string) {
	if strings.HasSuffix(value, dirSeparator) {
		add(currentDirectoryReplacements, value[0:len(value)-1])
	} else {
		add(currentDirectoryReplacements, value)
	}
}

func add(m map[string]struct{}, value string) {
	if _, exists := m[value]; !exists {
		m[value] = struct{}{}
	}
}

// Apply applies all the registered scrubbers to the target
func (s *dataScrubber) Apply(extension string, target *strings.Builder, settings *verifySettings) {
	stringData := target.String()
	target.Reset()

	for replace := range currentDirectoryReplacements {
		stringData = strings.ReplaceAll(stringData, replace, "{CurrentDirectory}")
	}

	for replace := range tempDirectoryReplacements {
		stringData = strings.ReplaceAll(stringData, replace, "{TempPath}")
	}

	for _, scrubber := range settings.instanceScrubbers {
		stringData = scrubber(stringData)
	}

	extensionScrubbers := settings.extensionMappedInstanceScrubbers[extension]
	for _, scrubber := range extensionScrubbers {
		stringData = scrubber(stringData)
	}

	stringData = fixNewlines(stringData)
	target.WriteString(stringData)
}

// ScrubTime scrubs time.Time values
func (s *dataScrubber) ScrubTime(time time.Time) string {
	if time.IsZero() {
		return "Time_Zero"
	}

	next := s.counter.GetNextTime(time)
	return fmt.Sprintf("Time_%d", next)
}

// ScrubGUID scrubs the UUID values
func (s *dataScrubber) ScrubGUID(guid uuid.UUID) string {
	if guid == uuid.Nil {
		return "Guid_Zero"
	}

	next := s.counter.GetNextUUID(guid)
	return fmt.Sprintf("Guid_%d", next)
}

// ScrubMachineName scrubs current hostname from the target
func (s *dataScrubber) ScrubMachineName(target string) string {
	if host, err := os.Hostname(); err == nil {
		if strings.Contains(target, host) {
			result := strings.ReplaceAll(target, host, "TheMachineName")
			return result
		}
	}
	return target
}

// ScrubStackTrace scrubs the stacktrace
func (s *dataScrubber) ScrubStackTrace(stacktrace string, removeParams bool) string {
	if len(stacktrace) == 0 {
		return stacktrace
	}

	//TODO: parse stacktrace

	return stacktrace
}

func (s *dataScrubber) removeLinesContaining(input string, ignoreCase bool, stringToMatch ...string) string {
	utils.Guard.AgainstNullOrEmptySlice(stringToMatch)

	return s.filterLines(input, func(line string) bool {
		return s.lineContains(line, ignoreCase, stringToMatch)
	})
}

func (s *dataScrubber) replaceLines(input string, lineReplacer func(string) string) string {
	result := strings.Builder{}
	lines := s.stringToLines(input)
	for i, line := range lines {
		if value := lineReplacer(line); len(value) > 0 {
			result.WriteString(value)
			if i != len(lines)-1 || strings.HasSuffix(input, "\n") {
				result.WriteRune('\n')
			}
		}
	}
	return result.String()
}

func (s *dataScrubber) filterLines(input string, removeLineFunc func(string) bool) string {
	result := strings.Builder{}
	lines := s.stringToLines(input)
	for i, line := range lines {
		if ok := removeLineFunc(line); !ok {
			result.WriteString(line)
			if i != len(lines)-1 || strings.HasSuffix(input, "\n") {
				result.WriteRune('\n')
			}
		}
	}
	return result.String()
}

func (s *dataScrubber) lineContains(line string, ignoreCase bool, stringToMatch []string) bool {
	for _, match := range stringToMatch {
		if ignoreCase {
			if strings.Contains(strings.ToLower(line), strings.ToLower(match)) {
				return true
			}
		} else if strings.Contains(line, match) {
			return true
		}
	}
	return false
}

func (s *dataScrubber) stringToLines(input string) (lines []string) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func (s *dataScrubber) replaceGuids(input string) string {
	if result, ok := s.tryReplaceGuids(input); ok {
		return result
	}
	return input
}

func (s *dataScrubber) tryReplaceGuids(value string) (string, bool) {

	if id, err := uuid.Parse(value); err == nil {
		return s.ScrubGUID(id), true
	}

	replaced := value
	guids := reg.FindAllString(value, -1)
	if len(guids) > 0 {
		for _, stringGuid := range guids {
			guid, err := uuid.Parse(stringGuid)
			if err == nil {
				convertedGuid := s.ScrubGUID(guid)
				replaced = strings.ReplaceAll(replaced, stringGuid, convertedGuid)
			}
		}
		return replaced, true
	}

	return "", false
}
