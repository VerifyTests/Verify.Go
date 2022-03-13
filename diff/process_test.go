package diff

import (
	"testing"
)

type testEnvReader struct {
	lookup map[string]string
}

func newTestReader() testEnvReader {
	return testEnvReader{
		lookup: make(map[string]string),
	}
}

func (t testEnvReader) LookupEnv(key string) (string, bool) {
	val, found := t.lookup[key]
	return val, found
}

func TestFetchingProcessesLists(t *testing.T) {
	p := processCleaner{}
	procs := p.findAllProcess()

	if procs == nil || len(procs) == 0 {
		t.Fatalf("processes should not be empty")
	}
}
