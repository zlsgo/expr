package expr

import (
	"testing"

	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/ztype"
)

func TestEval(t *testing.T) {
	tt := zlsgo.NewTest(t)

	resp, err := Eval(`1+1`, nil)
	tt.NoError(err, true)
	tt.Equal(2, resp.(int))

	resp, err = Eval(`"Hello "+name`, ztype.Map{"name": "World"})
	tt.NoError(err, true)
	tt.Log("Hello World", resp.(string))
}
