package verifier

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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

	if len(serialized) == 0 {
		t.Fatalf("serialized string should not be empty")
	}

	var p Person
	e := json.Unmarshal([]byte(serialized), &p)

	if e != nil {
		t.Fatalf("Should be able to unmarshal: %s", e)
	}
	if person.GivenNames != p.GivenNames {
		t.Fatalf("GivenNames was not equal: %s", p.GivenNames)
	}
	if person.FamilyName != p.FamilyName {
		t.Fatalf("FamilyName was not equal: %s", p.GivenNames)
	}
	if person.Address != p.Address {
		t.Fatalf("Address was not equal: %s", p.GivenNames)
	}
	if person.Dob != p.Dob {
		t.Fatalf("Dob was not equal: %s", p.GivenNames)
	}
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

	if len(serialized) == 0 {
		t.Fatalf("Should have serialized into string")
	}
	if !strings.Contains(serialized, "\"Dob\": \"Time_1\"") {
		t.Fatalf("Should have scrubbed Dob")
	}
}

func TestSerializeScrubbedZeroTime(t *testing.T) {

	person := Person{
		GivenNames: "John",
		Dob:        time.Time{},
	}

	serializer := getTestSerializer(true, true)
	serialized := serializer.Serialize(person)

	if len(serialized) == 0 {
		t.Fatalf("Should have serialized into string")
	}
	if !strings.Contains(serialized, "\"Dob\": \"Time_Zero\"") {
		t.Fatalf("Should have scrubbed Dob as zero time")
	}
}

func TestSerializeNonScrubbedZeroTime(t *testing.T) {

	person := Person{
		GivenNames: "John",
		Dob:        time.Time{},
	}

	serializer := getTestSerializer(true, false)
	serialized := serializer.Serialize(person)

	if len(serialized) == 0 {
		t.Fatalf("Should have serialized into string")
	}
	if !strings.Contains(serialized, "\"Dob\": \"0001-01-01T00:00:00Z\"") {
		t.Fatalf("Should have serialized Dob value")
	}
}

func TestSerializeScrubbedStringTypes(t *testing.T) {

	person := Record{
		StringValue:     "TestValue",
		GuidStringValue: uuid.New().String(),
		DateStringValue: time.Now().Format(time.RFC3339),
	}

	serializer := getTestSerializer(true, true)
	serialized := serializer.Serialize(person)

	if len(serialized) == 0 {
		t.Fatalf("Should have serialized into string")
	}
	if !strings.Contains(serialized, "TestValue") {
		t.Fatalf("Should have serialized 'TestValue'")
	}
	if !strings.Contains(serialized, "Guid_1") {
		t.Fatalf("Should have serialized 'Guid_1'")
	}
	if !strings.Contains(serialized, "Time_1") {
		t.Fatalf("Should have serialized 'Time_1'")
	}
}

func TestSerializeNotScrubbedStringTypes(t *testing.T) {

	person := Record{
		StringValue:     "TestValue",
		GuidStringValue: uuid.New().String(),
		DateStringValue: time.Now().Format(time.RFC3339),
	}

	serializer := getTestSerializer(false, false)
	serialized := serializer.Serialize(person)

	if len(serialized) == 0 {
		t.Fatalf("Should have serialized into string")
	}
	if !strings.Contains(serialized, "TestValue") {
		t.Fatalf("Should have serialized 'TestValue'")
	}
	if strings.Contains(serialized, "Guid_1") {
		t.Fatalf("Should have serialized 'Guid_1'")
	}
	if strings.Contains(serialized, "Time_1") {
		t.Fatalf("Should have serialized 'Time_1'")
	}
}

func TestSerializeStringBuilder(t *testing.T) {

	builder := strings.Builder{}
	builder.WriteString("FirstValue")
	builder.WriteString("\tSecondValue")

	serializer := getTestSerializer(true, true)
	serialized := serializer.Serialize(builder)

	if len(serialized) == 0 {
		t.Fatalf("Should have serialized into string")
	}
	if !strings.Contains(serialized, "FirstValue\tSecondValue") {
		t.Fatalf("Should have correct serialized values")
	}
}

func TestSerializeMultiLineString(t *testing.T) {

	builder := strings.Builder{}
	builder.WriteString("FirstValue\r\n")
	builder.WriteString("SecondValue\r\n")
	builder.WriteString("ThirdValue\n")

	serializer := getTestSerializer(true, true)
	serialized := serializer.Serialize(builder)
	t.Logf("Serialized: %s", serialized)

	if len(serialized) == 0 {
		t.Fatalf("Should have serialized into string")
	}
	if !strings.Contains(serialized, "FirstValue") {
		t.Fatalf("Should have correct serialized values")
	}
	if !strings.Contains(serialized, "SecondValue") {
		t.Fatalf("Should have correct serialized values")
	}
	if !strings.Contains(serialized, "ThirdValue") {
		t.Fatalf("Should have correct serialized values")
	}
}

func TestSerializeStringerInterface(t *testing.T) {

	ip := net.IPv4(192, 168, 0, 1)
	stringer := fmt.Stringer(ip) //ip as stringer

	serializer := getTestSerializer(true, true)
	serialized := serializer.Serialize(stringer)

	if len(serialized) == 0 {
		t.Fatalf("Should have serialized into string")
	}
	if !strings.Contains(serialized, "192.168.0.1") {
		t.Fatalf("Should have correct serialized values")
	}
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
