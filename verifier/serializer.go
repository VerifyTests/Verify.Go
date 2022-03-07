package verifier

import (
	"encoding"
	"fmt"
	"github.com/google/uuid"
	"github.com/heskandari/jsoner"
	"github.com/modern-go/reflect2"
	"strings"
	"sync"
	"time"
	"unsafe"
)

var (
	stringType        = reflect2.TypeOf("")
	textMarshalerType = reflect2.TypeOfPtr((*encoding.TextMarshaler)(nil)).Elem()
	stringBuilderType = reflect2.TypeOfPtr((*strings.Builder)(nil)).Elem()
	stringerType      = reflect2.TypeOfPtr((*fmt.Stringer)(nil)).Elem()
	timeType          = reflect2.TypeOfPtr((*time.Time)(nil)).Elem()
	uuidType          = reflect2.TypeOfPtr((*uuid.UUID)(nil)).Elem()
)

type serializer struct {
	scrubber *dataScrubber
	json     jsoner.API
	lock     *sync.Mutex
	settings *verifySettings
}

func newSerializer(settings *verifySettings, scrubber *dataScrubber) *serializer {
	serializer := serializer{
		settings: settings,
		scrubber: scrubber,
		lock:     &sync.Mutex{},
		json:     createMarshaller(),
	}

	serializer.registerExtensions()

	return &serializer
}

func (s *serializer) registerExtensions() {
	s.json.RegisterTypeEncoder(stringType.String(), newStringEncoder(s))
	s.json.RegisterTypeEncoder(stringerType.String(), newStringerEncoder(s))
	s.json.RegisterTypeEncoder(uuidType.String(), newUUIDEncoder(s))
	s.json.RegisterTypeEncoder(timeType.String(), newTimeEncoder(s))
	s.json.RegisterTypeEncoder(textMarshalerType.String(), newTextMarshallerEncoder(s))
}

func createMarshaller() jsoner.API {
	return jsoner.Config{
		EscapeHTML:    false,
		SortMapKeys:   true,
		IndentionStep: 4,
	}.Froze()
}

func (s *serializer) Serialize(v interface{}) string {
	s.lock.Lock()
	defer s.lock.Unlock()

	switch val := v.(type) {
	case fmt.Stringer:
		return s.convertStringer(val)
	case strings.Builder:
		return s.convertString(val.String())
	}
	return s.toJSON(v)
}

func (s *serializer) convertUUID(value uuid.UUID) string {
	if s.settings.scrubGuids {
		return s.scrubber.ScrubGUID(value)
	}
	return value.String()
}

func (s *serializer) convertTime(value time.Time) string {
	if s.settings.scrubTimes {
		return s.scrubber.ScrubTime(value)
	}
	converted := value.Format(time.RFC3339)
	return converted
}

func (s *serializer) toJSON(v interface{}) string {
	js, err := s.json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(fmt.Sprintf("failed to serialize to json: %s", err.Error()))
	}

	r := string(js)
	return r
}

type encoderUUID struct{ serializer *serializer }
type encoderTextMarshaller struct{ serializer *serializer }
type encoderTime struct{ serializer *serializer }
type encoderStringer struct{ serializer *serializer }
type encoderString struct{ serializer *serializer }

func newStringerEncoder(s *serializer) jsoner.ValEncoder {
	return &encoderStringer{
		serializer: s,
	}
}

func newStringEncoder(s *serializer) jsoner.ValEncoder {
	return &encoderString{
		serializer: s,
	}
}

func newTextMarshallerEncoder(s *serializer) jsoner.ValEncoder {
	return &encoderTextMarshaller{
		serializer: s,
	}
}

func newTimeEncoder(s *serializer) jsoner.ValEncoder {
	return &encoderTime{
		serializer: s,
	}
}

func newUUIDEncoder(s *serializer) jsoner.ValEncoder {
	return &encoderUUID{
		serializer: s,
	}
}

func (t encoderUUID) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (t encoderUUID) Encode(ptr unsafe.Pointer, stream *jsoner.Stream) {
	val := (*uuid.UUID)(ptr)
	stream.WriteString(t.serializer.convertUUID(*val))
}

func (t encoderTextMarshaller) IsEmpty(ptr unsafe.Pointer) bool {
	val := (*encoding.TextMarshaler)(ptr)
	return val == nil
}

func (t encoderTextMarshaller) Encode(ptr unsafe.Pointer, stream *jsoner.Stream) {
	val := *(*encoding.TextMarshaler)(ptr)
	str, err := val.MarshalText()
	if err != nil {
		panic("Failed to marshal text")
	}
	stream.WriteString(t.serializer.convertString(string(str)))
}

func (t encoderString) IsEmpty(ptr unsafe.Pointer) bool {
	val := (*string)(ptr)
	return val == nil
}

func (t encoderString) Encode(ptr unsafe.Pointer, stream *jsoner.Stream) {
	val := *(*string)(ptr)
	stream.WriteString(t.serializer.convertString(val))
}

func (t encoderStringer) IsEmpty(ptr unsafe.Pointer) bool {
	val := (*fmt.Stringer)(ptr)
	return val == nil
}

func (t encoderStringer) Encode(ptr unsafe.Pointer, stream *jsoner.Stream) {
	val := *(*fmt.Stringer)(ptr)
	stream.WriteString(t.serializer.convertString(val.String()))
}

func (t encoderTime) IsEmpty(ptr unsafe.Pointer) bool {
	val := (*time.Time)(ptr)
	return val == nil
}

func (t encoderTime) Encode(ptr unsafe.Pointer, stream *jsoner.Stream) {
	val := *(*time.Time)(ptr)
	stream.WriteString(t.serializer.convertTime(val))
}

func asJSON(input interface{}, appenders []toAppend, settings *verifySettings) *strings.Builder {
	if len(appenders) > 0 {
		dictionary := make(map[string]interface{})
		if input == nil {
			dictionary["target"] = "nil"
		} else {
			dictionary["target"] = input
		}

		input = dictionary
		for _, appender := range appenders {
			dictionary[appender.Name] = appender.Data
		}
	}

	serializer := newSerializer(settings, settings.scrubber)
	serialized := serializer.Serialize(input)

	builder := strings.Builder{}
	builder.WriteString(serialized)

	return &builder
}

func (s *serializer) convertString(value string) string {
	converted, valid := s.tryParseConvertGUID(value)
	if valid {
		return converted
	}

	converted, valid = s.tryParseDateTime(value)
	if valid {
		return converted
	}

	return fixNewlines(value)
}

func (s *serializer) convertStringer(str fmt.Stringer) string {
	stringValue := (str).String()
	return s.convertString(stringValue)
}

func (s *serializer) tryParseConvertGUID(value string) (string, bool) {
	if s.settings.scrubGuids {
		parsed, err := uuid.Parse(value)
		if err != nil {
			return value, false
		}
		return s.scrubber.ScrubGUID(parsed), true
	}
	return value, false
}

func (s *serializer) tryParseDateTime(value string) (string, bool) {
	if s.settings.scrubTimes {
		parsed, err := time.Parse(time.RFC3339, value)
		if err == nil {
			return s.scrubber.ScrubTime(parsed), true
		}

		parsed, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", value) // Stringer date format
		if err == nil {
			return s.scrubber.ScrubTime(parsed), true
		}
	}
	return value, false
}
