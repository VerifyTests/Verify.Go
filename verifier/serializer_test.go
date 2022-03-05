package verifier

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net"
	"strings"
	"testing"
	"time"
)

type Address struct {
	Street  string
	Country string
	Suburb  string
}

type Person struct {
	GivenNames string
	FamilyName string
	Spouse     string
	Address    Address
	Dob        time.Time
}

type Record struct {
	StringValue     string
	GuidStringValue string
	DateStringValue string
}

func TestSerializeAndDeserialize(t *testing.T) {

	person := Person{
		GivenNames: "John",
		FamilyName: "Smith",
		Spouse:     "Jane",
		Dob:        time.Date(2000, 10, 1, 0, 0, 0, 0, time.UTC),
	}

	serializer := getTestSerializer(true, false)
	serialized := serializer.Serialize(person)

	assert.NotEmpty(t, serialized)

	var p Person
	e := json.Unmarshal([]byte(serialized), &p)

	assert.NoError(t, e)
	assert.Equal(t, p.GivenNames, person.GivenNames)
	assert.Equal(t, p.FamilyName, person.FamilyName)
	assert.Equal(t, p.Dob, person.Dob)
	assert.Equal(t, p.Address, person.Address)
}

func TestSerializeScrubbedTime(t *testing.T) {

	person := Person{
		GivenNames: "John",
		FamilyName: "Smith",
		Spouse:     "Jane",
		Dob:        time.Date(2000, 10, 1, 0, 0, 0, 0, time.Local),
	}

	serializer := getTestSerializer(true, true)
	serialized := serializer.Serialize(person)

	assert.NotEmpty(t, serialized)
	assert.Contains(t, serialized, "\"Dob\": \"Time_1\"")
}

func TestSerializeScrubbedZeroTime(t *testing.T) {

	person := Person{
		GivenNames: "John",
		Dob:        time.Time{},
	}

	serializer := getTestSerializer(true, true)
	serialized := serializer.Serialize(person)

	assert.NotEmpty(t, serialized)
	assert.Contains(t, serialized, "\"Dob\": \"Time_Zero\"")
}

func TestSerializeNonScrubbedZeroTime(t *testing.T) {

	person := Person{
		GivenNames: "John",
		Dob:        time.Time{},
	}

	serializer := getTestSerializer(true, false)
	serialized := serializer.Serialize(person)

	assert.NotEmpty(t, serialized)
	assert.Contains(t, serialized, "\"Dob\": \"0001-01-01T00:00:00Z\"")
}

func TestSerializeScrubbedStringTypes(t *testing.T) {

	person := Record{
		StringValue:     "TestValue",
		GuidStringValue: uuid.New().String(),
		DateStringValue: time.Now().Format(time.RFC3339),
	}

	serializer := getTestSerializer(true, true)
	serialized := serializer.Serialize(person)

	assert.NotEmpty(t, serialized)
	assert.Contains(t, serialized, "TestValue")
	assert.Contains(t, serialized, "Guid_1")
	assert.Contains(t, serialized, "Time_1")
}

func TestSerializeNotScrubbedStringTypes(t *testing.T) {

	person := Record{
		StringValue:     "TestValue",
		GuidStringValue: uuid.New().String(),
		DateStringValue: time.Now().Format(time.RFC3339),
	}

	serializer := getTestSerializer(false, false)
	serialized := serializer.Serialize(person)

	assert.NotEmpty(t, serialized)
	assert.Contains(t, serialized, "TestValue")
	assert.NotContains(t, serialized, "Guid_1")
	assert.NotContains(t, serialized, "Time_1")
}

func TestSerializeStringBuilder(t *testing.T) {

	builder := strings.Builder{}
	builder.WriteString("FirstValue")
	builder.WriteString("\tSecondValue")

	serializer := getTestSerializer(true, true)
	serialized := serializer.Serialize(builder)

	assert.NotEmpty(t, serialized)
	assert.Contains(t, serialized, "FirstValue\tSecondValue")
}

func TestSerializeMultiLineString(t *testing.T) {

	builder := strings.Builder{}
	builder.WriteString("FirstValue\r\n")
	builder.WriteString("SecondValue\r\n")
	builder.WriteString("ThirdValue\n")

	serializer := getTestSerializer(true, true)
	serialized := serializer.Serialize(builder)
	t.Logf("Serialized: %s", serialized)

	assert.NotEmpty(t, serialized)
	assert.Contains(t, serialized, "FirstValue")
	assert.Contains(t, serialized, "SecondValue")
	assert.Contains(t, serialized, "ThirdValue")
}

func TestSerializeStringerInterface(t *testing.T) {

	ip := net.IPv4(192, 168, 0, 1)
	stringer := fmt.Stringer(ip) //ip as stringer

	serializer := getTestSerializer(true, true)
	serialized := serializer.Serialize(stringer)

	assert.NotEmpty(t, serialized)
	assert.Contains(t, serialized, "192.168.0.1")
}

func getTestSerializer(scrubGuid, scrubTime bool) *serializer {
	settings := NewSettings()
	if !scrubGuid {
		settings.DontScrubGuids()
	}

	if !scrubTime {
		settings.DontScrubTimes()
	}

	serializer := newSerializer(settings.(*verifySettings), newDataScrubber(startCounter()))
	return serializer
}
