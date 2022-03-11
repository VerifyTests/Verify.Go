package api_tests_test

import (
	"github.com/VerifyTests/Verify.Go/verifier"
	"github.com/google/uuid"
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

func NewTestSettings() verifier.VerifySettings {
	defaultSettings := verifier.NewSettings()
	defaultSettings.UseDirectory("../_testdata")
	return defaultSettings
}

func TestVerifyingNilObject(t *testing.T) {
	verifier.VerifyWithSetting(t, NewTestSettings(), nil)
}

func TestSimpleString(t *testing.T) {
	settings := NewTestSettings()
	settings.DisableDiff()
	verifier.VerifyWithSetting(t, settings, "Foo")
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
	verifier.VerifyWithSetting(t, NewTestSettings(), person)
}

func TestVerifyingNullableStructs(t *testing.T) {
	address := &Address{
		Street: "Test Street",
	}
	verifier.VerifyWithSetting(t, NewTestSettings(), address)
}

func TestVerifyingMultiLineStrings(t *testing.T) {
	settings := NewTestSettings()

	settings.ScrubLinesWithReplace(func(line string) string {
		if strings.Contains(line, "LineE") {
			return "NoMoreLineE"
		}
		return line
	})
	settings.ScrubLines(func(line string) bool {
		return strings.Contains(line, "J")
	})
	settings.ScrubLinesContaining("b", "D")
	settings.ScrubLinesContainingAnyCase("h")

	verifier.VerifyWithSetting(t, settings, "LineA\nLineB\nLineC\nLineD\nLineE\nLineF\nLineG\nLineH\nLineI\nLineJ")
}

type ToBeScrubbed struct {
	RowVersion string
}

func TestScrubbingAfterMarshalling(t *testing.T) {
	var target = ToBeScrubbed{
		RowVersion: "7D3",
	}

	settings := NewTestSettings()
	settings.AddScrubber(func(target string) string {
		return strings.Replace(target, "7D3", "TheRowVersion", 1)
	})

	verifier.VerifyWithSetting(t, settings, target)
}

type nonPublic struct {
	field string
}

func TestVerifyingPrivateStructs(t *testing.T) {
	target := nonPublic{}
	verifier.VerifyWithSetting(t, NewTestSettings(), target)
}

func TestRemovingEmptyLines(t *testing.T) {
	target := `
LineA

LineC`
	settings := NewTestSettings()
	settings.ScrubEmptyLines()
	verifier.VerifyWithSetting(t, settings, target)
}

func TestVerifySlicesOfStructs(t *testing.T) {
	persons := []uuid.UUID{
		uuid.New(), uuid.New(),
	}
	verifier.VerifyWithSetting(t, NewTestSettings(), persons)
}

func TestVerifyArrayOfStructs(t *testing.T) {
	persons := [2]Person{
		{GivenNames: "John", Title: Mr},
		{GivenNames: "Jill", Title: Mrs},
	}

	verifier.VerifyWithSetting(t, NewTestSettings(), persons)
}

func TestVerifySliceOfTime(t *testing.T) {
	times := []time.Time{
		time.Date(2020, 01, 02, 1, 2, 3, 4, time.Local),
		time.Date(2021, 01, 02, 1, 2, 3, 4, time.Local),
	}

	verifier.VerifyWithSetting(t, NewTestSettings(), times)
}

func TestVerifySliceOfStrings(t *testing.T) {
	times := []string{
		"This is a text",
		"Second string",
		uuid.NewString(),
	}

	verifier.VerifyWithSetting(t, NewTestSettings(), times)
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

	verifier.VerifyWithSetting(t, NewTestSettings(), target)
}