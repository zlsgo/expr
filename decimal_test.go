package expr

import (
	"testing"

	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/ztype"
)

func TestCompileDecimal(t *testing.T) {
	tt := zlsgo.NewTest(t)

	data := ztype.Map{
		"GMV":      0.1,
		"Platform": 3,
	}

	pDecimal, err := CompileDecimal(`GMV*Platform`, data)
	tt.NoError(err, true)
	respDecimal, err := Run(pDecimal, data)
	tt.NoError(err)
	tt.Equal(0.3, ztype.ToFloat64(respDecimal))

	p, err := Compile(`GMV*Platform`, data)
	tt.NoError(err, true)
	resp, err := Run(p, data)
	tt.NoError(err)
	tt.Equal(0.30000000000000004, ztype.ToFloat64(resp))

	tt.Log("正常计算", resp, "高精确计算", respDecimal)
}
