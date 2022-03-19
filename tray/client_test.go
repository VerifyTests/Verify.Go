//go:build integration
// +build integration

package tray

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClient_SendDeleteAndSendMove(t *testing.T) {
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
	client.SendDelete("testFile.txt")
	client.SendMove("test.received.txt", "test.verified.txt", "", "", false, 0)

	assert.Eventually(t, func() bool {
		return deleted != nil &&
			deleted.File == "testFile.txt"
	}, 9*time.Second, 3*time.Second)

	assert.Eventually(t, func() bool {
		return moved != nil &&
			moved.Temp == "test.received.txt" &&
			moved.Target == "test.verified.txt"
	}, 60*time.Second, 3*time.Second)
}
