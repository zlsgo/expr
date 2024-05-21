package expr

import (
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
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

func Run(p *Program, env any) (any, error) {
	return vm.Run(p.program, env)
}

func Eval(input string, env any) (any, error) {
	p, err := Compile(input)
	if err != nil {
		return nil, err
	}
	return Run(p, env)
}
