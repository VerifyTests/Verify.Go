package tray

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestClient_SendDeleteAndSendMove_Integration(t *testing.T) {
	runTests, _ := strconv.ParseBool(os.Getenv("RUN_INTEGRATION_TESTS"))
	if !runTests {
		t.Skip("Skipping integration tests")
	}

	var deleted *DeletePayload
	var moved *MovePayload

	server := NewServer()
	server.DeleteHandler = func(cmd *DeletePayload) {
		deleted = cmd
	}
	server.MoveHandler = func(cmd *MovePayload) {
		moved = cmd
	}

	server.Start()
	defer server.Stop()

	client := NewClient()
	client.AddDelete("testFile.txt")
	client.AddMove("test.received.txt", "test.verified.txt", "", nil, false, 0)

	assert.Eventually(t, func() bool {
		return deleted != nil &&
			strings.Contains(deleted.File, "testFile.txt")
	}, 9*time.Second, 3*time.Second)

	assert.Eventually(t, func() bool {
		return moved != nil &&
			strings.Contains(moved.Temp, "test.received.txt") &&
			strings.Contains(moved.Target, "test.verified.txt")
	}, 9*time.Second, 3*time.Second)
}
