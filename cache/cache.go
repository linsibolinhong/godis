package cache

import "github.com/linsibolinhong/godis/command"

type CommandFunc func(cmd *command.Command) (*command.Result, error)

type Cache interface {
	Command(cmd *command.Command) (*command.Result, error)
}
