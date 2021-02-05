package command

type Result struct {
	Ret []string
}

func NewResult() *Result {
	return &Result{Ret: []string{}}
}

func (r *Result) AppendRet(s string) {
	r.Ret = append(r.Ret, s)
}