package command

type Method string

const (
	MethodNull = Method("null")
	MethodGet = Method("get")
)

type Command struct {
	cmd Method
	params [][]byte
}


func NewCommand() *Command {
	return &Command{
		params:[][]byte{},
	}
}

func (c *Command) AppendParam(p []byte) {
	c.params = append(c.params, p)
}

func (c *Command) Parse() {

}

func (c *Command) ToString() string {
	ret := string(c.cmd) + "\n"
	for _, param := range c.params {
		ret += string(param) + "\n"
	}
	return ret
}