package cache

import (
	"fmt"
	"github.com/linsibolinhong/godis/command"
)

type simpleCache struct {
	kv map[string]string
	cmdMap map[command.Method]CommandFunc
}

func NewSimpleCache() Cache {
	c := &simpleCache{
		kv:     map[string]string{},
		cmdMap: map[command.Method]CommandFunc{},
	}

	c.cmdMap[command.MethodCmd] = c.Ping
	c.cmdMap[command.MethodNull] = c.Ping
	c.cmdMap[command.MethodGet] = c.Get
	c.cmdMap[command.MethodSet] = c.Set
	return c
}

func (c *simpleCache) Command(cmd *command.Command) (*command.Result, error) {
	if cmd == nil {
		return nil, fmt.Errorf("param is nil")
	}

	fn, found := c.cmdMap[cmd.Cmd]
	if !found {
		return nil, fmt.Errorf("unknown method %v", cmd.Cmd)
	}

	return fn(cmd)
}

func (c *simpleCache) Ping(cmd *command.Command) (*command.Result, error) {
	return nil, nil
}

func (c *simpleCache) Get(cmd *command.Command) (*command.Result, error) {
	val, found := c.kv[cmd.Params[0]]
	if !found {
		return nil, fmt.Errorf("key not exist")
	}
	ret := command.NewResult()
	ret.AppendRet(val)
	return ret, nil
}

func (c *simpleCache) Set(cmd *command.Command) (*command.Result, error) {
	key := cmd.Params[0]
	val := cmd.Params[1]
	c.kv[key] =val
	ret := command.NewResult()
	return ret, nil
}