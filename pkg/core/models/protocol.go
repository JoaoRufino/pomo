package models

type Protocol struct {
	Cid     CmdID
	Payload Payload
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
	Cmd_UpdateTask
)
