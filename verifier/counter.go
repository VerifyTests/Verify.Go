package verifier

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

type countHolder struct {
	currentID     int
	currentGUID   int
	currentTime   int
	idCache       map[interface{}]int
	counterLocker *sync.Mutex
}

func (c *countHolder) GetNextID(input interface{}) int {
	c.counterLocker.Lock()
	defer c.counterLocker.Unlock()

	val, found := c.idCache[input]
	if found {
		return int(val)
	}

	c.currentID++
	c.idCache[input] = c.currentID

	return c.currentID
}

func (c *countHolder) GetNextUUID(input uuid.UUID) int {
	c.counterLocker.Lock()
	defer c.counterLocker.Unlock()

	if val, ok := c.idCache[input]; ok {
		return int(val)
	}

	c.currentGUID++
	c.idCache[input] = c.currentGUID

	return c.currentGUID
}

func (c *countHolder) GetNextTime(input time.Time) int {
	c.counterLocker.Lock()
	defer c.counterLocker.Unlock()

	if val, ok := c.idCache[input]; ok {
		return int(val)
	}

	c.currentTime++
	c.idCache[input] = c.currentTime

	return c.currentTime
}

func startCounter() *countHolder {
	return &countHolder{
		counterLocker: &sync.Mutex{},
		idCache:       make(map[interface{}]int),
	}
}
