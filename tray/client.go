package tray

import (
	"github.com/VerifyTests/Verify.Go/utils"
	"log"
	"net"
	"strconv"
)

// Client a client to work with the server interface
type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (d *Client) Connects() bool {
	conn, err := d.getConnection()
	if err != nil {
		return false
	}
	err = conn.Close()
	if err != nil {
		return false
	}
	return true
}

// SendDelete sends delete information to the server
func (d *Client) SendDelete(file string) {
	d.sendDelete(file)
}

// SendMove sends move information to the server
func (d *Client) SendMove(tempFile, targetFile, exe string, arguments []string, canKill bool, processId int32) {
	d.sendMove(tempFile, targetFile, exe, arguments, canKill, processId)
}

func (d *Client) sendDelete(file string) {
	payload := DeletePayload{
		Type: "Delete",
		File: file,
	}

	data, err := serialize(payload)
	if err != nil {
		log.Printf("failed to serialize delete data: %s", err)
	}

	d.sendData(data)
}

func (d *Client) sendMove(temp, target, exe string, arguments []string, canKill bool, processId int32) {
	payload := MovePayload{
		Type:      "Move",
		Target:    utils.File.GetFullPath(target),
		Temp:      utils.File.GetFullPath(temp),
		Exe:       exe,
		Arguments: arguments,
		CanKill:   canKill,
		ProcessId: processId,
	}

	data, err := serialize(payload)
	if err != nil {
		log.Printf("failed to serialize move data: %s", err)
	}

	d.sendData(data)
}

func (d *Client) getConnection() (net.Conn, error) {
	servAddr := "localhost:" + strconv.Itoa(ServerPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (d *Client) sendData(data []byte) {
	conn, err := d.getConnection()
	if err != nil {
		log.Printf("Could not get a connection to the tray app: %s", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		log.Printf("Failed to write data to the tray app: %s", err)
	}

	_ = conn.Close()
}
