// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/phcurtis/fn"
)

const pkgName = "fn_test"
const baseName = "github.com/phcurtis/" + pkgName + "."

func Test_fngroup(t *testing.T) {
	fns := "Test_fngroup.func1"
	tests := []struct {
		name string
		fnt  interface{}
		arg1 int
		want string
	}{
		{"fn.Cur().....", fn.Cur, 0, baseName + fns},
		{"fn.Lvl(0)....", fn.Lvl, 0, baseName + fns},
		{"fn.LvlBase(0)", fn.LvlBase, 0, pkgName + "." + fns},
		{"fn.CurBase(0)", fn.CurBase, 0, pkgName + "." + fns},
	}
	for i, v := range tests {
		var got, argstr string
		t.Run(v.name, func(t *testing.T) {
			switch f := v.fnt.(type) {
			case func() string:
				got = f()
				argstr = "()"
			case func(int) string:
				got = f(v.arg1)
				argstr = fmt.Sprintf("(%d)", v.arg1)
			default:
				t.Fatalf("item in list is unsupport func type, i:=%d\n", i)
			}
			if got != v.want {
				t.Errorf("%s%s \n  got:%s \n want:%s\n", v.name, argstr, got, v.want)
			}
		})
	}
}

func a2() string        { return fn.CStk() }
func a1() string        { return a2() }
func b2(lvl int) string { return fn.LvlCStk(lvl) }
func b1(lvl int) string { return b2(lvl) }

func Test_cstkgroup(t *testing.T) {
	fns := "Test_cstkgroup.func1<--"
	tests := []struct {
		name     string
		fnt      interface{}
		arg1     int
		wantPfix string
	}{
		{"...CStk(0)-2deep", a1, 0, baseName + "a2<--" + baseName + "a1<--" + baseName + fns},
		{"LvlCStk(0)-2deep", b1, 0, baseName + "b2<--" + baseName + "b1<--" + baseName + fns},
		{"LvlCStk(1)-2deep", b1, 1, baseName + "b1<--" + baseName + fns},
		{"LvlCStk(2)-2deep", b1, 2, baseName + fns},
	}
	var got string
	for i, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			switch f := v.fnt.(type) {
			case func() string:
				got = f()
			case func(int) string:
				got = f(v.arg1)
			default:
				t.Fatalf("item in list is unsupport func type, i:=%d\n", i)
			}
			if !strings.HasPrefix(got, v.wantPfix) {
				t.Errorf("%s:\n     got:%s \nwantPfix:%s \n", v.name, got, v.wantPfix)
			}
		})
	}
}
