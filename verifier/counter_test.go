package verifier

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestGetNextTime(t *testing.T) {
	var count = startCounter()

	current := time.Now()
	first := count.GetNextTime(current)
	assert.Equal(t, 1, first)

	repeat := count.GetNextTime(current) // Same
	assert.Equal(t, 1, repeat)

	second := count.GetNextTime(time.Now())
	assert.Equal(t, 2, second)
}

func TestGetNextUUID(t *testing.T) {
	counter := startCounter()

	guid1, _ := uuid.NewUUID()
	guid2, _ := uuid.NewUUID()

	first := counter.GetNextUUID(guid1)
	assert.Equal(t, 1, first)

	repeat := counter.GetNextUUID(guid1) // Same value
	assert.Equal(t, 1, repeat)

	second := counter.GetNextUUID(guid2)
	assert.Equal(t, 2, second)
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

	assert.Equal(t, 1001, counter.GetNextID(1001))
}
