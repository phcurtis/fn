// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn_test

import (
	"fmt"
	"strconv"
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

func TestLvlInfo(t *testing.T) {
	filenameshort := "fn_test.go"
	filenamenogps := "github.com/phcurtis/fn" + "/" + filenameshort
	funcname := "TestLvlInfo.func1"
	baseFN := pkgName + "." + funcname
	fullFN := "github.com/phcurtis/" + baseFN
	parens := "()"

	tests := []struct {
		name   string
		lvl    int
		iflags int
		wfile  string
		wname  string
		prefix bool
	}{
		{"1-Ifnbase|Ifilenogps", 0, fn.Ifnbase | fn.Ifilenogps,
			filenamenogps, baseFN + parens, false},
		{"2-Ifileshort", 0, fn.Ifileshort,
			filenameshort, fullFN + parens, false},
		{"3-Ifnbase", 5, fn.IflagsDef, "???", fn.CStkEndPfix, true},
	}
	for _, v := range tests {
		var file string
		var line int
		var name string
		t.Run(v.name, func(t *testing.T) {
			file, line, name = fn.LvlInfo(v.lvl, v.iflags)
			if file != v.wfile {
				t.Errorf("%s: fn.LvlInfo(%d,%d): \n filegot:%s \nfilewant:%s \n",
					v.name, v.lvl, v.iflags, file, v.wfile)
			} else if name != v.wname && (!v.prefix) {
				t.Errorf("%s: fn.LvlInfo(%d,%d): \n namegot:%s \nnamewant:%s \n",
					v.name, v.lvl, v.iflags, name, v.wname)
			} else if line < 1 && v.lvl < 1 {
				t.Errorf("%s: fn.LvlInfo(%d,%d): \n linegot:%d \nlinewant:> 0\n",
					v.name, v.lvl, v.iflags, line)
			} else if v.prefix && strings.HasPrefix(v.wname, name) {
				t.Errorf("%s: fn.LvlInfo(%d,%d): \n prefixnamegot:%s \nprefixnamewant:%s \n",
					v.name, v.lvl, v.iflags, name, v.wname)
			}
		})
	}
}

func Test_lvlinfostrings(t *testing.T) {
	//filenamegps := "/home/paul/go/src/"
	filenameshort := "fn_test.go"
	filenamenogps := "github.com/phcurtis/fn" + "/" + filenameshort
	funcname := "fn_test.Test_lvlinfostrings()"

	tests := []struct {
		name     string
		got      string
		wantfile string
		wantfunc string
	}{
		{"fn.LvlInfoStr(0)..:", fn.LvlInfoStr(0, fn.IflagsDef), filenamenogps, funcname},
		{"fn.LvlInfoCmn(0)..:", fn.LvlInfoCmn(0), filenamenogps, funcname},
		{"fn.LvlInfoShort(0):", fn.LvlInfoShort(0), filenameshort, funcname},
	}
	for _, v := range tests {
		parts := strings.Split(v.got, ":")
		if len(parts) != 3 {
			t.Errorf("number parts wrong")
		}
		if parts[0] != v.wantfile {
			t.Errorf("gotfilename:%q wantfilenamebad line number:%q", parts[0], v.wantfile)
		}
		line, err := strconv.Atoi(parts[1])
		if err != nil || line < 1 {
			t.Errorf("bad line number:%q", parts[1])
		}
		if parts[2] != v.wantfunc {
			t.Errorf("gotfuncname:%q wantfuncname:%q", parts[2], v.wantfunc)
		}
	}
}
