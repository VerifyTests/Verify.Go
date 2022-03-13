package verifier

import (
	"github.com/google/uuid"
	"sync"
	"testing"
	"time"
)

func TestGetNextTime(t *testing.T) {
	var count = startCounter()

	current := time.Now()
	first := count.GetNextTime(current)
	if first != 1 {
		t.Fatalf("Should generate correct time")
	}

	repeat := count.GetNextTime(current) // Same
	if repeat != 1 {
		t.Fatalf("Should generate correct time")
	}

	second := count.GetNextTime(time.Now())
	if second != 2 {
		t.Fatalf("Should generate correct time")
	}
}

func TestGetNextUUID(t *testing.T) {
	counter := startCounter()

	guid1, _ := uuid.NewUUID()
	guid2, _ := uuid.NewUUID()

	first := counter.GetNextUUID(guid1)
	if first != 1 {
		t.Fatalf("Should generate correct guid")
	}

	repeat := counter.GetNextUUID(guid1) // Same value
	if repeat != 1 {
		t.Fatalf("Should generate correct guid")
	}

	second := counter.GetNextUUID(guid2)
	if second != 2 {
		t.Fatalf("Should generate correct guid")
	}
}

func TestConcurrentNext(t *testing.T) {
	wantedCount := 1000
	var wg sync.WaitGroup
	wg.Add(wantedCount)

	counter := startCounter()

	for i := 0; i < wantedCount; i++ {
		go func(c *countHolder, val int) {
			c.GetNextID(val)
			wg.Done()
		}(counter, i)
	}
	wg.Wait()

	if counter.GetNextID(1001) != 1001 {
		t.Fatalf("Should generate correct id")
	}
}
