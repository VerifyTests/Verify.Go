package diff

import (
	"testing"
)

type TestEnvReader struct {
	Key   string
	Value string
}

func (t *TestEnvReader) LookupEnv(key string) (string, bool) {
	if t.Key == key {
		return t.Value, true
	}
	return "", false
}

func TestRunner(t *testing.T) {
	table := []struct {
		key      string
		value    string
		expected bool
	}{
		{string(Appveyor), "true", true},
		{string(Appveyor), "TRUE", true},
		{string(MyGet), "true", false},
		{string(MyGet), "myget", true},
		{string(GitHub), "address", true},
		{string(GitHub), "", false},
		{string(AzureDevOps), "address", true},
		{string(TeamCity), "address", true},
		{string(MyGet), "address", false},
		{string(MyGet), "true", false},
		{string(MyGet), "", false},
		{string(MyGet), "false", false},
		{string(MyGet), "myget", true},
		{string(GitLab), "address", true},
		{string(GitLab), "", false},
		{"DiffEngine_Disabled", "", false},
		{"DiffEngine_Disabled", "false", false},
		{"DiffEngine_Disabled", "true", true},
	}

	e := func(env string, value string) TestEnvReader {
		return TestEnvReader{
			Key:   env,
			Value: value,
		}
	}

	for _, row := range table {
		reader := e(row.key, row.value)
		runner := newRunner(&reader)
		runner.logger.EnableLogging()

		check := checkDisabled(&reader)

		if row.expected != runner.disabled {
			t.Fatalf("failed to get expected disabled status")
		}
		if row.expected != check {
			t.Fatalf("failed to get correct check status")
		}
	}
}
