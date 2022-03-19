package tray

// MovePayload move file information
type MovePayload struct {
	Type      string
	Temp      string
	Target    string
	Exe       string
	Arguments string
	CanKill   bool
	ProcessId int
}

// DeletePayload delete file information
type DeletePayload struct {
	Type string
	File string
}

// MoveMessageHandler handle the move operation
type MoveMessageHandler func(cmd *MovePayload)

// DeleteMessageHandler handle the delete operation
type DeleteMessageHandler func(cmd *DeletePayload)

// Server struct for handling delete and move operation
type Server struct {
	processor     chan string
	stopped       bool
	MoveHandler   MoveMessageHandler
	DeleteHandler DeleteMessageHandler
}
