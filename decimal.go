package expr

import (
	"errors"
	"reflect"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/shopspring/decimal"
	"github.com/sohaha/zlsgo/ztype"
)

type (
	d struct {
		Value float64
	}
	decimalPatcher struct{}
)

var decimalType = reflect.TypeOf(d{})

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
		default:
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
		px float64
		py float64
	)

	if v, ok := params[0].(d); !ok {
		px = ztype.ToFloat64(params[0])
	} else {
		px = v.Value
	}
	if v, ok := params[1].(d); !ok {
		py = ztype.ToFloat64(params[1])
	} else {
		py = v.Value
	}

	return decimal.NewFromFloat(px), decimal.NewFromFloat(py), nil
}

var decimalHandler = []expr.Option{
	expr.Function(
		"_Mul",
		func(params ...any) (any, error) {
			x, y, err := toDecimal(params...)
			if err != nil {
				return nil, err
			}
			f, _ := x.Mul(y).Float64()
			return f, nil
		},
	),
	expr.Function(
		"_Div",
		func(params ...any) (any, error) {
			x, y, err := toDecimal(params...)
			if err != nil {
				return nil, err
			}
			f, _ := x.Div(y).Float64()
			return f, nil
		},
	),
	expr.Function(
		"_Add",
		func(params ...any) (any, error) {
			x, y, err := toDecimal(params...)
			if err != nil {
				return nil, err
			}
			f, _ := x.Add(y).Float64()
			return f, nil
		},
	),
	expr.Function(
		"_Sub",
		func(params ...any) (any, error) {
			x, y, err := toDecimal(params...)
			if err != nil {
				return nil, err
			}
			f, _ := x.Sub(y).Float64()
			return f, nil
		},
	),
	expr.Patch(decimalPatcher{}),
}

func CompileDecimal(code string, env ...any) (*Program, error) {
	ops := []expr.Option{}
	for i := range env {
		if o, ok := env[i].(expr.Option); ok {
			ops = append(ops, o)
		} else {
			ops = append(ops, expr.Env(env[i]))
		}
	}
	program, err := expr.Compile(code, append(decimalHandler, ops...)...)
	if err != nil {
		return nil, err
	}

	return &Program{program: program}, nil
}
