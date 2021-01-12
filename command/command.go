package command

type Command struct {
	Params [][]byte
}


func NewCommand() *Command {
	return &Command{
		Params:[][]byte{},
	}
}


func (c *Command) ToString() string {
	ret := ""
	for _, param := range c.Params {
		ret += string(param) + "\n"
	}
	return ret
}