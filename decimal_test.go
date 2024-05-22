package expr

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/ztype"
)

func TestDecimal(t *testing.T) {
	tt := zlsgo.NewTest(t)

	data := ztype.Map{
		"GMV":      0.1,
		"Platform": 3,
	}

	pDecimal, err := CompileDecimal(`GMV*Platform`, data)
	tt.NoError(err, true)
	respDecimal, err := Run(pDecimal, data)
	tt.NoError(err)
	tt.Equal("0.3", respDecimal.String())

	p, err := Compile(`GMV*Platform`, data)
	tt.NoError(err, true)
	resp, err := Run(p, data)
	tt.NoError(err)
	tt.Equal("0.30000000000000004", resp.String())

	tt.Log("正常计算", resp.String(), "高精确计算", respDecimal.String())

	data3 := ztype.Map{
		"GMV":      decimal.NewFromFloat(0.1),
		"Platform": decimal.NewFromInt(3),
	}
	pDecimal3, err := CompileDecimal(`GMV*Platform`, data3)
	tt.NoError(err)
	respDecimal3, err := Run(pDecimal3, data3)
	tt.NoError(err)
	tt.Equal(respDecimal, respDecimal3)
}

func TestDecimalExpr(t *testing.T) {
	tt := zlsgo.NewTest(t)

	data := ztype.Map{
		"GMV":      0.1,
		"Platform": 3,
	}

	resp, err := EvalDecimal(`GMV+Platform`, data)
	tt.NoError(err, true)
	tt.Log(resp.Value())

	resp, err = EvalDecimal(`GMV-Platform`, data)
	tt.NoError(err, true)
	tt.Log(resp.Value())

	resp, err = EvalDecimal(`GMV*Platform`, data)
	tt.NoError(err, true)
	tt.Log(resp.Value())

	resp, err = EvalDecimal(`GMV/Platform`, data)
	tt.NoError(err, true)
	tt.Log(resp.Value())

	resp, err = EvalDecimal(`GMV%Platform`, data)
	tt.NoError(err, true)
	tt.Log(resp.Value())

	type value struct {
		env      ztype.Map
		expected any
	}

	a, _ := decimal.NewFromString("8.32500499700179891")
	b, _ := decimal.NewFromString("8.3250049970017989")
	c, _ := decimal.NewFromString("1.000000000000000000321")
	d, _ := decimal.NewFromString("9.325004997001798910321")
	for input, val := range map[string]value{
		`"Hello "+name`: {ztype.Map{"name": "World"}, "Hello World"},
		`name=="World"`: {ztype.Map{"name": "World"}, true},
		`v.name`:        {ztype.Map{"v": ztype.Map{"name": "one"}}, "one"},
		`v[1].name`:     {ztype.Map{"v": ztype.Maps{{"name": "one"}, {"name": "two"}}}, "two"},
		`1+1`:           {nil, decimal.NewFromInt(2)},
		`1>1`:           {nil, false},
		`2>1`:           {nil, true},
		`1<1+3`:         {nil, true},
		`1+1>1`:         {nil, true},
		`2>=1`:          {nil, true},
		`1<=2`:          {nil, true},
		`2<=2`:          {nil, true},
		`3<=2`:          {nil, false},
		`a<b`:           {ztype.Map{"a": a, "b": b}, false},
		`a>b`:           {ztype.Map{"a": a, "b": b}, true},
		`a>=b`:          {ztype.Map{"a": a, "b": b}, true},
		`a<=b`:          {ztype.Map{"a": a, "b": b}, false},
		`a+c`:           {ztype.Map{"a": a, "c": c}, d},
	} {
		resp, err := EvalDecimal(input, val.env)
		tt.NoError(err, true)
		if ztype.ToString(val.expected) != ztype.ToString(resp.Value()) {
			tt.Fatal(input+": ", val.expected, "!=", resp.Value())
		}
	}
}
