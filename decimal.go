package expr

import (
	"errors"
	"reflect"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/shopspring/decimal"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/ztype"
)

type (
	decimalPatcher struct {
	}
)

var (
	decimalType = reflect.TypeOf(decimal.Decimal{})
)

func (decimalPatcher) isNumber(typ reflect.Type) bool {
	if typ == decimalType {
		return true
	}

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}

func (d decimalPatcher) Visit(node *ast.Node) {
	if n, ok := (*node).(*ast.BinaryNode); ok {
		if !d.isNumber(n.Left.Type()) || !d.isNumber(n.Right.Type()) {
			return
		}

		var callee *ast.IdentifierNode
		switch n.Operator {
		case "+":
			callee = &ast.IdentifierNode{Value: "_Add"}
		case "-":
			callee = &ast.IdentifierNode{Value: "_Sub"}
		case "*":
			callee = &ast.IdentifierNode{Value: "_Mul"}
		case "/":
			callee = &ast.IdentifierNode{Value: "_Div"}
		case "%":
			callee = &ast.IdentifierNode{Value: "_Mod"}
		case ">":
			callee = &ast.IdentifierNode{Value: "_Gt"}
		case "<":
			callee = &ast.IdentifierNode{Value: "_Lt"}
		case "==":
			callee = &ast.IdentifierNode{Value: "_Eq"}
		case ">=":
			callee = &ast.IdentifierNode{Value: "_Gte"}
		case "<=":
			callee = &ast.IdentifierNode{Value: "_Lte"}
		default:
			zlog.Debug(n.Operator)
			return
		}

		ast.Patch(node, &ast.CallNode{
			Callee:    callee,
			Arguments: []ast.Node{n.Left, n.Right},
		})
		(*node).SetType(decimalType)
	}
}

func toDecimal(params ...any) (x, y decimal.Decimal, err error) {
	if len(params) != 2 {
		err = errors.New("params length must be 2")
		return
	}

	var (
		ok bool
	)

	if x, ok = params[0].(decimal.Decimal); !ok {
		if x, err = decimal.NewFromString(ztype.ToString(params[0])); err != nil {
			return
		}
	}

	if y, ok = params[1].(decimal.Decimal); !ok {
		if y, err = decimal.NewFromString(ztype.ToString(params[1])); err != nil {
			return
		}
	}

	return
}

var decimalHandler = []expr.Option{
	expr.Function(
		"_Eq",
		func(params ...any) (any, error) {
			x, y, err := toDecimal(params...)
			if err != nil {
				return nil, err
			}
			return x.Cmp(y) == 0, nil
		},
		new(func(x, y any) (decimal.Decimal, error)),
	),
	expr.Function(
		"_Gt",
		func(params ...any) (any, error) {
			x, y, err := toDecimal(params...)
			if err != nil {
				return nil, err
			}
			return x.GreaterThan(y), nil
		},
		new(func(x, y any) (decimal.Decimal, error)),
	),
	expr.Function(
		"_Gte",
		func(params ...any) (any, error) {
			x, y, err := toDecimal(params...)
			if err != nil {
				return nil, err
			}
			return !(x.Cmp(y) == -1), nil
		},
		new(func(x, y any) (decimal.Decimal, error)),
	),
	expr.Function(
		"_Lt",
		func(params ...any) (any, error) {
			x, y, err := toDecimal(params...)
			if err != nil {
				return nil, err
			}
			return x.LessThan(y), nil
		},
		new(func(x, y any) (decimal.Decimal, error)),
	),
	expr.Function(
		"_Lte",
		func(params ...any) (any, error) {
			x, y, err := toDecimal(params...)
			if err != nil {
				return nil, err
			}
			return !(x.Cmp(y) == 1), nil
		},
		new(func(x, y any) (decimal.Decimal, error)),
	),
	expr.Function(
		"_Mul",
		func(params ...any) (any, error) {
			x, y, err := toDecimal(params...)
			if err != nil {
				return nil, err
			}
			return x.Mul(y), nil
		},
		new(func(x, y any) (decimal.Decimal, error)),
	),
	expr.Function(
		"_Div",
		func(params ...any) (any, error) {
			x, y, err := toDecimal(params...)
			if err != nil {
				return nil, err
			}

			return x.Div(y), nil
		},
		new(func(x, y any) (decimal.Decimal, error)),
	),
	expr.Function(
		"_Add",
		func(params ...any) (any, error) {
			x, y, err := toDecimal(params...)
			if err != nil {
				return nil, err
			}
			return x.Add(y), nil
		},
		new(func(x, y any) (decimal.Decimal, error)),
	),
	expr.Function(
		"_Sub",
		func(params ...any) (any, error) {
			x, y, err := toDecimal(params...)
			if err != nil {
				return nil, err
			}
			return x.Sub(y), nil
		},
		new(func(x, y any) (decimal.Decimal, error)),
	),
	expr.Function(
		"_Mod",
		func(params ...any) (any, error) {
			x, y, err := toDecimal(params...)
			if err != nil {
				return nil, err
			}
			return x.Mod(y), nil
		},
		new(func(x, y any) (decimal.Decimal, error)),
	),
	expr.Patch(decimalPatcher{}),
}

func CompileDecimal(input string, env any, ops ...any) (*Program, error) {
	o := make([]expr.Option, 0, len(ops)+1)
	ops = append(ops, env)
	for i := range ops {
		if ops[i] == nil {
			continue
		}
		if op, ok := ops[i].(expr.Option); ok {
			o = append(o, op)
		} else {
			o = append(o, expr.Env(ops[i]))
		}
	}

	program, err := expr.Compile(input, append(decimalHandler, o...)...)
	if err != nil {
		return nil, err
	}

	return &Program{program: program}, nil
}

func EvalDecimal(input string, env any) (value ztype.Type, err error) {
	var p *Program
	p, err = CompileDecimal(input, env)
	if err != nil {
		return
	}

	return Run(p, env)
}
