package expr

import (
	"testing"

	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/ztype"
)

func TestEval(t *testing.T) {
	tt := zlsgo.NewTest(t)

	type value struct {
		env      ztype.Map
		expected any
	}

	for input, val := range map[string]value{
		`"Hello "+name`: {ztype.Map{"name": "World"}, "Hello World"},
		`name=="World"`: {ztype.Map{"name": "World"}, true},
		`v.name`:        {ztype.Map{"v": ztype.Map{"name": "one"}}, "one"},
		`v[1].name`:     {ztype.Map{"v": ztype.Maps{{"name": "one"}, {"name": "two"}}}, "two"},
		`1+1`:           {nil, 2},
		`1>1`:           {nil, false},
		`2>1`:           {nil, true},
		`1<1+3`:         {nil, true},
		`1+1>1`:         {nil, true},
		`2>=1`:          {nil, true},
		`1<=2`:          {nil, true},
		`2<=2`:          {nil, true},
		`3<=2`:          {nil, false},
	} {
		resp, err := Eval(input, val.env)
		tt.NoError(err, true)
		tt.Equal(val.expected, resp.Value())
	}
}
