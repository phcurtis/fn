// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/phcurtis/fn"
)

func Test_pkgcfggroup(t *testing.T) {
	fnpass := "Test_pkgcfggroup - test passed "

	// subTESTS:
	// #1 (PkgCfgDef)            verify PkgCfgDef returns expected values
	// #2 (PkgCfg)               verify individual setter funcs alter what makes up a PkgCfg
	// #3 (SetPkgCfgDef/PkgCfg)) verify rewriting defaults matches what should be in the pkgCfgDef
	// #4 (PkgCfg)               verify setter funcs mods match a corresponding fetched pkgCfg
	// #5 (SetPkgCfg/PkgCfg))    verify setting a pkgcfg matches returned values from PkgCfg
	// on completion of these test restore pkgCfgDef

	// prelims
	var got, want *fn.PkgCfgStruct
	var giowr, wiowr io.Writer
	var gotstr, wantstr, pkgdefwantstr string
	defiowr := fn.LogGetOutputDef()
	note := fmt.Sprintf("\n    Note:os.Stdout:%v \n    Note:os.Stderr:%v", os.Stdout, os.Stderr)
	pkgDefWant := &fn.PkgCfgStruct{
		Logflags:      fn.LflagsDef,
		LogPrefix:     fn.LogPrefixDef,
		LogTraceFlags: fn.TrFlagsDef,
		LogAlignFile:  fn.LogAlignFileDef,
		LogAlignFunc:  fn.LogAlignFuncDef,
	}
	pkgdefwantstr = fmt.Sprintf("%+v", pkgDefWant)

	// #1 (PkgCfgDef)    verify PkgCfgDef returns expected values
	want = pkgDefWant
	wiowr = defiowr
	got, giowr = fn.PkgCfgDef()
	wantstr = pkgdefwantstr
	gotstr = fmt.Sprintf("%+v", got)
	if gotstr != wantstr {
		t.Errorf("#1: PkgCfgDef() incorrect:\n got:%s \nwant:%s\n", gotstr, wantstr)
	}
	if testing.Verbose() {
		log.Println(fnpass + "#1 (PkgCfgDef)")
	}

	// #2 (PkgCfg)       verify individual setter funcs alter what makes up a PkgCfg
	fn.LogSetAlignFile(5)
	fn.LogSetAlignFunc(6)
	fn.LogSetFlags(0xffffffff)
	fn.LogSetPrefix("pRe")
	fn.LogSetTraceFlags(0xdeadbeef) //3735928559
	wiowr = os.Stderr
	fn.LogSetOutput(wiowr)
	got, giowr = fn.PkgCfg()
	gotstr = fmt.Sprintf("%+v", got)
	if gotstr == wantstr {
		t.Errorf("#2: PkgCfg() incorrect:\n got:%s \nwant:%s \n", gotstr, wantstr)
	}
	if giowr != wiowr {
		t.Errorf("#2a: PkgCfg() incorrect:\n gotiowr:%x \nwantowr:%x %s\n",
			giowr, wiowr, note)
	}
	if testing.Verbose() {
		log.Println(fnpass + "#2 (PkgCfg)")
	}

	// #3 (SetPkgCfgDef/PkgCfg) verify rewriting the defaults matches what should be in the pkg defaults
	fn.SetPkgCfgDef(true)
	wiowr = defiowr
	got, giowr = fn.PkgCfg()
	gotstr = fmt.Sprintf("%+v", got)
	if gotstr != wantstr {
		t.Errorf("#3: SetPkgCfgDef() incorrect:\n got:%s \nwant:%s %s\n", gotstr, wantstr, note)
	}
	if giowr != wiowr {
		t.Errorf("#3a: PkgCfgDef() incorrect:\n gotiowr:%x \nwantiowr:%x %s\n", giowr, wiowr, note)
	}
	if testing.Verbose() {
		log.Println(fnpass + "#3 (SetPkgCfgDef/PkgCfg)")
	}

	// #4 (PkgCfg)       verify setter funcs mods  match a corresponding fetched pkgCfg
	fn.SetPkgCfgDef(true)
	wiowr = os.Stderr
	fn.LogSetOutput(wiowr)
	fn.LogSetFlags(0xffff)
	fn.LogSetPrefix("ZyxAbcd")
	fn.LogSetTraceFlags(0xbeef) //dec=48879
	fn.LogSetAlignFile(11)
	fn.LogSetAlignFunc(12)
	got, giowr = fn.PkgCfg()
	want = &fn.PkgCfgStruct{
		Logflags:      0xffff,
		LogPrefix:     "ZyxAbcd",
		LogTraceFlags: 0xbeef, //dec=48879
		LogAlignFile:  11,
		LogAlignFunc:  12,
	}
	got, giowr = fn.PkgCfg()
	gotstr = fmt.Sprintf("%+v", got)
	wantstr = fmt.Sprintf("%+v", want)
	if gotstr != wantstr {
		t.Errorf(" #4: PkgCfgDef() incorrect:\n got:%s \nwant:%s \n", gotstr, wantstr)
	}
	if giowr != wiowr {
		t.Errorf("#4a: PkgCfgDef() incorrect:\n gotiowr:%x \nwantiowr:%x %s\n", giowr, wiowr, note)
	}
	if testing.Verbose() {
		log.Println(fnpass + "#4 (PkgCfg)")
	}

	// #5 (SetPkgCfg/PkgCfg)    verify setting a pkgcfg matches returned values from  PkgCfg
	want.LogAlignFile = 8
	wiowr = os.Stdout
	fn.SetPkgCfg(want, wiowr)
	got, giowr = fn.PkgCfg()
	gotstr = fmt.Sprintf("%+v", got)
	wantstr = fmt.Sprintf("%+v", want)
	if gotstr != wantstr {
		t.Fatalf(" #5: SetPkgCfg() incorrect:\n got:%s \nwant:%s ", gotstr, wantstr)
	}
	if giowr != wiowr {
		t.Fatalf("#5a: SetPkgCfg() incorrect:\n gotiowr:%x \nwantiowr:%x ", giowr, wiowr)
	}
	if testing.Verbose() {
		log.Println(fnpass + "#5 (SetPkgCfg/PkgCfg)")
	}

	// reset settings to package defaults
	fn.SetPkgCfgDef(true)
	if testing.Verbose() {
		log.Println("restored pkgCfgDef")
	}
}
