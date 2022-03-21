package tray

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
)

// ServerPort the port server listens to
const ServerPort = 3492

func NewServer() *Server {
	srv := &Server{
		processor: make(chan string, 1),
		//processorStop: make(chan bool),
		//listenerStop:  make(chan bool),
	}
	return srv
}

func (t *Server) Start() {
	log.Printf("Starting the server")
	go t.startReceiver()
	go t.startProcessor(t.processor)
}

func (t *Server) Stop() {
	log.Printf("Stopping the server")
	t.stopped = true
	close(t.processor)
}

func (t *Server) startProcessor(processor <-chan string) {
	for {
		if t.stopped {
			break
		}

		message := <-processor

		if len(message) == 0 {
			continue
		}

		if strings.Contains(message, "\"Type\":\"Move\"") {
			moveCommand := MovePayload{}
			deserialize(message, &moveCommand)
			t.moveFile(&moveCommand)
		} else if strings.Contains(message, "\"Type\":\"Delete\"") {
			deleteCommand := DeletePayload{}
			deserialize(message, &deleteCommand)
			t.deleteFile(&deleteCommand)
		} else {
			log.Printf("Unknown payload: %s", message)
		}
	}
}

func (t *Server) startReceiver() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(ServerPort))
	if err != nil {
		panic(fmt.Sprintf("Failed to start receiver: %s", err))
	}

	log.Printf("Server is ready")

	for {
		if t.stopped {
			break
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s", err)
			continue
		}

		log.Printf("Accepted connection")
		go func(reader io.Reader) {
			bytes, err := ioutil.ReadAll(reader)
			if err != nil {
				log.Printf("Failed to read: %s", err)
				return
			}
			message := string(bytes)

			t.processor <- message
		}(conn)
	}

	log.Printf("Closed receiver")
}

func (t *Server) deleteFile(command *DeletePayload) {
	log.Printf("Delete: %+v", command)
	t.DeleteHandler(command)
}

func (t *Server) moveFile(command *MovePayload) {
	log.Printf("Move: %+v", command)
	t.MoveHandler(command)
}

func deserialize(payload string, obj interface{}) {
	err := json.Unmarshal([]byte(payload), obj)
	if err != nil {
		fmt.Println(err)
	}
}

func serialize(obj interface{}) ([]byte, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
