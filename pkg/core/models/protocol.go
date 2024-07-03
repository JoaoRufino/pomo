package models

type Protocol struct {
	Cid     CmdID
	Payload Payload
	Error   *string `json:"omitempty"`
}

type Payload interface {
}

type CmdID int

const (
	Cmd_CreateTask CmdID = iota
	Cmd_CreatePomodoro
	Cmd_DeleteTask
	Cmd_GetList
	Cmd_GetServerStatus
	Cmd_GetTask
	Cmd_UpdateStatus
)

const (
	ErrWrongMessageType = "wrong message type provided: got %d but wanted %d"
	ErrWrongDataType    = "wrong data type: got type %T but wanted %T"
)
