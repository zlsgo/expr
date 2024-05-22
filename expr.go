package expr

import (
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/sohaha/zlsgo/ztype"
)

type Program struct {
	program *vm.Program
}

func Compile(input string, env ...any) (*Program, error) {
	ops := []expr.Option{}
	for i := range env {
		if o, ok := env[i].(expr.Option); ok {
			ops = append(ops, o)
		} else {
			ops = append(ops, expr.Env(env[i]))
		}
	}
	program, err := expr.Compile(input, ops...)
	if err != nil {
		return nil, err
	}

	return &Program{program: program}, nil
}

func Run(p *Program, env any) (value ztype.Type, err error) {
	var resp any
	resp, err = vm.Run(p.program, env)
	if err != nil {
		return
	}

	value = ztype.New(resp)
	return
}

func Eval(input string, env any) (value ztype.Type, err error) {
	var p *Program
	p, err = Compile(input)
	if err != nil {
		return
	}
	return Run(p, env)
}
