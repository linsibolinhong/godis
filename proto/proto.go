package proto

import "github.com/linsibolinhong/godis/command"

type Proto interface {
	ReadCommand() (*command.Command, error)
	WriteResult(result *command.Result) error
	WriteError(err error) error
}