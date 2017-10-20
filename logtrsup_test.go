// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn_test

import (
	"io"
	"os"
	"testing"

	"github.com/phcurtis/fn"
)

func Test_logpairedfuncsgroup(t *testing.T) {
	tests := []struct {
		name   string
		fnt    interface{}
		fnt2   interface{}
		inwant interface{}
		res    bool
	}{
		{"LogAlignFile(15).........", fn.LogSetAlignFile, fn.LogAlignFile, 15, true},
		{"LogAlignFile(0)..........", fn.LogSetAlignFile, fn.LogAlignFile, 0, true},
		{`LogAlignFile('max'+1)....`, fn.LogSetAlignFile, fn.LogAlignFile, fn.LogAlignFileMax + 1, false},
		{`LogAlignFile(-1).........`, fn.LogSetAlignFile, fn.LogAlignFile, -1, false},

		{"LogAlignFunc(15).........", fn.LogSetAlignFunc, fn.LogAlignFunc, 15, true},
		{"LogAlignFunc(0)..........", fn.LogSetAlignFunc, fn.LogAlignFunc, 0, true},
		{`LogAlignFunc('max'+1)....`, fn.LogSetAlignFunc, fn.LogAlignFunc, fn.LogAlignFuncMax + 1, false},
		{`LogAlignFunc(-1).........`, fn.LogSetAlignFunc, fn.LogAlignFunc, -1, false},

		{"LogFlags(0xdeadbeef).....", fn.LogSetFlags, fn.LogFlags, 0xdeadbeef, true},
		{"LogGetOutput(os.Stderr)..", fn.LogSetOutput, fn.LogGetOutput, os.Stderr, true},
		{"LogGetOutput(os.Stdout)..", fn.LogSetOutput, fn.LogGetOutput, os.Stdout, true},
		{`LogPrefix("pReFIX")......`, fn.LogSetPrefix, fn.LogPrefix, "pReFIX", true},
		{"LogTraceFlags(0xbeefdead)", fn.LogSetTraceFlags, fn.LogTraceFlags, 0xbeefdead, true},
	}

	for i, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			switch f := v.fnt.(type) {
			case func(int):
				want := v.inwant.(int)
				f(want)
				f2 := v.fnt2.(func() int)
				got := f2()
				if (got == want) != v.res {
					t.Errorf("%s \n  got:0x%x \n want:0x%x wantres:%t\n", v.name, got, want, v.res)
				}
			case func(io.Writer):
				want := v.inwant.(io.Writer)
				f(want)
				f2 := v.fnt2.(func() io.Writer)
				got := f2()
				if (got == want) != v.res {
					t.Errorf("%s \n  got:%s \n want:%s wantres:%t\n", v.name, got, want, v.res)
				}
			case func(string):
				want := v.inwant.(string)
				f(want)
				f2 := v.fnt2.(func() string)
				got := f2()
				if (got == want) != v.res {
					t.Errorf("%s \n  got:%s \n want:%s wantres:%t\n", v.name, got, want, v.res)
				}
			default:
				t.Fatalf("item in list is unsupport func type, i:=%d\n", i)
			}
		})
	}
}
