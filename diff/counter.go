package diff

import (
	"strconv"
	"sync"
)

const defaultMaxInstance = 5

type instanceCounter struct {
	launchedInstances   int
	maxInstanceToLaunch int
	locker              sync.Mutex
}

func newInstanceCounter(reader EnvReader) *instanceCounter {
	counter := instanceCounter{}
	counter.maxInstanceToLaunch = getMaxInstances(reader)
	return &counter
}

func (c *instanceCounter) ReachedMax() bool {
	c.locker.Lock()
	defer c.locker.Unlock()

	c.launchedInstances += 1
	return c.launchedInstances > c.maxInstanceToLaunch
}

func getMaxInstances(reader EnvReader) int {
	variable, found := reader.LookupEnv("DiffEngine_MaxInstances")
	if !found {
		return defaultMaxInstance
	}

	parsedCount, err := strconv.ParseInt(variable, 10, 16)
	if err != nil {
		panic("Could not parse the DiffEngine_MaxInstances environment variable: " + variable)
	}

	return int(parsedCount)
}
