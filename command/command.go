package command

type Method string

const (
	MethodNull = Method("null")
	MethodGet = Method("get")
	MethodCmd = Method("command")
	MethodSet = Method("set")
)

type Command struct {
	Cmd Method
	Params []string
}


func NewCommand() *Command {
	return &Command{
		Params:[]string{},
	}
}

func (c *Command) AppendParam(p string) {
	c.Params = append(c.Params, p)
}

func (c *Command) Parse() {

}

func (c *Command) ToString() string {
	ret := string(c.Cmd) + "\n"
	for _, param := range c.Params {
		ret += string(param) + "\n"
	}
	return ret
}