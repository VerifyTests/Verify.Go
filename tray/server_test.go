package tray

import (
	"github.com/VerifyTests/Verify.Go/verifier"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var deleteMsg = `
{
"Type":"Delete",
"File":"Foo"
}`

var moveMsg = `{
"Type":"Move",
"CanKill":true,
"Exe":"theExe",
"Arguments":"TheArguments \"s\"",
"ProcessId":10
}`

var testSettings verifier.VerifySettings

func init() {
	testSettings = verifier.NewSettings()
	testSettings.UseDirectory("../_testdata")
}

func TestTracker_ReceivingMove(t *testing.T) {

	var receivedDelete *DeletePayload
	var receivedMove *MovePayload
	var server = NewServer()
	defer server.Stop()

	var receiveChan = make(chan string)

	server.DeleteHandler = func(cmd *DeletePayload) {
		receivedDelete = cmd
	}
	server.MoveHandler = func(cmd *MovePayload) {
		receivedMove = cmd
	}

	go func() {
		receiveChan <- moveMsg
	}()

	go func() {
		server.startProcessor(receiveChan)
	}()

	shouldReceive := func() bool {
		return receivedMove != nil && receivedDelete == nil
	}

	assert.Eventually(t, shouldReceive, 5*time.Second, 2*time.Second)

	verifier.VerifyWithSetting(t, testSettings, receivedMove)
}

func TestTracker_ReceivingDelete(t *testing.T) {

	var receivedDelete *DeletePayload
	var receivedMove *MovePayload
	var receiveChan = make(chan string)
	var server = NewServer()
	defer server.Stop()

	server.processor = receiveChan
	server.DeleteHandler = func(cmd *DeletePayload) {
		receivedDelete = cmd
	}
	server.MoveHandler = func(cmd *MovePayload) {
		receivedMove = cmd
	}

	go func() {
		receiveChan <- deleteMsg
	}()

	go func() {
		server.startProcessor(receiveChan)
	}()

	shouldReceive := func() bool {
		return receivedDelete != nil && receivedMove == nil
	}

	assert.Eventually(t, shouldReceive, 5*time.Second, 2*time.Second)

	verifier.VerifyWithSetting(t, testSettings, receivedDelete)
}
