package api_tests_test

import (
	"github.com/VerifyTests/Verify.Go/verifier"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"testing"
	"time"
)

type Person struct {
	GivenNames string    `json:"given_names"`
	FamilyName string    `json:"family_name"`
	Spouse     string    `json:"spouse"`
	Address    Address   `json:"address"`
	Children   []string  `json:"children"`
	Title      Title     `json:"title"`
	ID         uuid.UUID `json:"id"`
	Dob        time.Time `json:"dob"`
}

type Address struct {
	Street  string  `json:"street"`
	Suburb  *string `json:"suburb,omitempty"`
	Country *string `json:"country,omitempty"`
}

type Title int

const (
	Mr  Title = 1
	Mrs       = 2
)

func (t Title) String() string {
	switch t {
	case Mr:
		return "Mr."
	case Mrs:
		return "Mrs."
	}
	return "Undefined Title"
}

func DefaultTestConfiguration() []verifier.VerifyConfigure {
	return []verifier.VerifyConfigure{
		verifier.UseDirectory("../_testdata"),
		verifier.DisableDiff(),
	}
}

func NewTestVerifier(t *testing.T) verifier.Verifier {
	cfg := DefaultTestConfiguration()
	return verifier.NewVerifier(t).Configure(cfg...)
}

func TestVerifyingNilObject(t *testing.T) {
	NewTestVerifier(t).Verify(nil)
}

func TestSimpleString(t *testing.T) {
	NewTestVerifier(t).Verify("Foo")
}

func TestVerifyingStructs(t *testing.T) {
	var country = "USA"
	var person = Person{
		ID:         uuid.MustParse("ebced679-45d3-4653-8791-3d969c4a986c"),
		Title:      Mr,
		FamilyName: "Smith",
		GivenNames: "John",
		Dob:        time.Date(2022, 02, 01, 1, 2, 3, 4, time.Local),
		Spouse:     "Jill",
		Address: Address{
			Street:  "4 Puddle Lane",
			Country: &country,
		},
		Children: []string{"Sam", "Mary"},
	}
	NewTestVerifier(t).Verify(person)
}

func TestVerifyingNullableStructs(t *testing.T) {
	address := &Address{
		Street: "Test Street",
	}
	NewTestVerifier(t).Verify(address)
}

func TestVerifyingMultiLineStrings(t *testing.T) {
	v := NewTestVerifier(t).Configure(
		verifier.ScrubLinesWithReplace(func(line string) string {
			if strings.Contains(line, "LineE") {
				return "NoMoreLineE"
			}
			return line
		}),
		verifier.ScrubLines(func(line string) bool {
			return strings.Contains(line, "J")
		}),
		verifier.ScrubLinesContaining("b", "D"),
		verifier.ScrubLinesContainingAnyCase("h"),
	)

	v.Verify("LineA\nLineB\nLineC\nLineD\nLineE\nLineF\nLineG\nLineH\nLineI\nLineJ")
}

type ToBeScrubbed struct {
	RowVersion string
}

func TestScrubbingAfterMarshalling(t *testing.T) {
	var target = ToBeScrubbed{
		RowVersion: "7D3",
	}

	v := NewTestVerifier(t).Configure(
		verifier.AddScrubber(func(target string) string {
			return strings.Replace(target, "7D3", "TheRowVersion", 1)
		}))

	v.Verify(target)
}

type nonPublic struct {
	field string
}

func TestVerifyingPrivateStructs(t *testing.T) {
	target := nonPublic{}
	NewTestVerifier(t).Verify(target)
}

func TestRemovingEmptyLines(t *testing.T) {
	target := `
LineA

LineC`

	v := NewTestVerifier(t).Configure(verifier.ScrubEmptyLines())
	v.Verify(target)
}

func TestVerifySlicesOfStructs(t *testing.T) {
	persons := []uuid.UUID{
		uuid.New(), uuid.New(),
	}
	NewTestVerifier(t).Verify(persons)
}

func TestVerifyArrayOfStructs(t *testing.T) {
	persons := [2]Person{
		{GivenNames: "John", Title: Mr},
		{GivenNames: "Jill", Title: Mrs},
	}

	NewTestVerifier(t).Verify(persons)
}

func TestVerifySliceOfTime(t *testing.T) {
	times := []time.Time{
		time.Date(2020, 01, 02, 1, 2, 3, 4, time.Local),
		time.Date(2021, 01, 02, 1, 2, 3, 4, time.Local),
	}

	NewTestVerifier(t).Verify(times)
}

func TestVerifySliceOfStrings(t *testing.T) {
	times := []string{
		"This is a text",
		"Second string",
		uuid.NewString(),
	}

	NewTestVerifier(t).Verify(times)
}

func TestVerifyMaps(t *testing.T) {
	t.Skipf("Brittle, skipping for now.")
	current := time.Now()
	target := map[string]interface{}{
		"FirstString":      "String value",
		"SecondUuidString": uuid.NewString(),
		"ThirdUuid":        uuid.New(),
		"FourthTime":       current,
		"FifthTimeString":  current.Format(time.RFC3339),
	}

	NewTestVerifier(t).Verify(target)
}

func TestUsingTableTests(t *testing.T) {
	type test struct {
		input    interface{}
		testName string
	}

	tests := []test{
		{testName: "guids", input: uuid.New()},
		{testName: "time", input: time.Now()},
		{testName: "integers", input: strconv.Itoa(10)},
	}

	v := NewTestVerifier(t).Configure(
		verifier.ScrubInlineGuids(),
		verifier.ScrubInlineTime(time.RFC3339))

	for _, tc := range tests {
		v = v.Configure(verifier.TestCase(tc.testName))
		v.Verify(tc.input)
	}
}

func TestUsingTableTestsWithSubTests(t *testing.T) {
	type test struct {
		input    interface{}
		testName string
	}

	tests := []test{
		{testName: "guids", input: uuid.New()},
		{testName: "time", input: time.Now()},
		{testName: "integers", input: strconv.Itoa(10)},
	}

	v := NewTestVerifier(t).Configure(
		verifier.ScrubInlineGuids(),
		verifier.ScrubInlineTime(time.RFC3339))

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			v = v.Configure(verifier.TestCase(tc.testName))
			v.Verify(tc.input)
		})
	}
}

func TestUsingTableTestsWithSubTestsAndExplicitTestCaseName(t *testing.T) {
	type test struct {
		input    interface{}
		testName string
	}

	tests := []test{
		{testName: "with guids", input: uuid.New()},
		{testName: "with time", input: time.Now()},
		{testName: "with integers", input: strconv.Itoa(10)},
	}

	v := NewTestVerifier(t).Configure(
		verifier.ScrubInlineGuids(),
		verifier.ScrubInlineTime(time.RFC3339))

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			v = v.Configure(verifier.TestCase("case-" + tc.testName[5:]))
			v.Verify(tc.input)
		})
	}
}
