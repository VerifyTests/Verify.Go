package verifier

import (
	"encoding"
	"fmt"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
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
	json     jsoniter.API
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
	s.json.RegisterTypeEncoder(stringType.String(), NewStringEncoder(s))
	s.json.RegisterTypeEncoder(stringerType.String(), NewStringerEncoder(s))
	s.json.RegisterTypeEncoder(uuidType.String(), NewUUIDEncoder(s))
	s.json.RegisterTypeEncoder(timeType.String(), NewTimeEncoder(s))
	s.json.RegisterTypeEncoder(textMarshalerType.String(), NewTextMarshallerEncoder(s))
}

func createMarshaller() jsoniter.API {
	return jsoniter.Config{
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

type UUIDEncoder struct{ serializer *serializer }
type TextMarshallerEncoder struct{ serializer *serializer }
type TimeEncoder struct{ serializer *serializer }
type StringerEncoder struct{ serializer *serializer }
type StringEncoder struct{ serializer *serializer }

func NewStringerEncoder(s *serializer) jsoniter.ValEncoder {
	return &StringerEncoder{
		serializer: s,
	}
}

func NewStringEncoder(s *serializer) jsoniter.ValEncoder {
	return &StringEncoder{
		serializer: s,
	}
}

func NewTextMarshallerEncoder(s *serializer) jsoniter.ValEncoder {
	return &TextMarshallerEncoder{
		serializer: s,
	}
}

func NewTimeEncoder(s *serializer) jsoniter.ValEncoder {
	return &TimeEncoder{
		serializer: s,
	}
}

func NewUUIDEncoder(s *serializer) jsoniter.ValEncoder {
	return &UUIDEncoder{
		serializer: s,
	}
}

func (t UUIDEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	//val := (*uuid.UUID)(ptr)
	//if val == nil {
	//	return true
	//}
	//
	//guid := *val
	//if guid == uuid.Nil {
	//	return true
	//}

	return false
}

func (t UUIDEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	val := (*uuid.UUID)(ptr)
	stream.WriteString(t.serializer.convertUUID(*val))
}

func (t TextMarshallerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	val := (*encoding.TextMarshaler)(ptr)
	return val == nil
}

func (t TextMarshallerEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	val := *(*encoding.TextMarshaler)(ptr)
	str, err := val.MarshalText()
	if err != nil {
		panic("Failed to marshal text")
	}
	stream.WriteString(t.serializer.convertString(string(str)))
}

func (t StringEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	val := (*string)(ptr)
	return val == nil
}

func (t StringEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	val := *(*string)(ptr)
	stream.WriteString(t.serializer.convertString(val))
}

func (t StringerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	val := (*fmt.Stringer)(ptr)
	return val == nil
}

func (t StringerEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	val := *(*fmt.Stringer)(ptr)
	stream.WriteString(t.serializer.convertString(val.String()))
}

func (t TimeEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	val := (*time.Time)(ptr)
	return val == nil
}

func (t TimeEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
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
