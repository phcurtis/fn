// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/phcurtis/fn"
)

var pkgCfgDefWant = &fn.PkgCfgStruct{
	LogFlags:      fn.LflagsDef,
	LogPrefix:     fn.LogPrefixDef,
	LogTraceFlags: fn.TrFlagsDef,
	LogAlignFile:  fn.LogAlignFileDef,
	LogAlignFunc:  fn.LogAlignFuncDef,
}
var pkgCfgDefWantstr = fmt.Sprintf("%+v", pkgCfgDefWant)
var noteStdio = fmt.Sprintf("\n    Note:os.Stdout:%v \n    Note:os.Stderr:%v", os.Stdout, os.Stderr)

func TestPkgCfgDef(t *testing.T) {
	// verify PkgCfgDef returns expected values

	fn.SetPkgCfgDef(true) // set pkg config to (what should be) a known default state
	got, giowr := fn.PkgCfgDef()
	gotstr := fmt.Sprintf("%+v", got)
	wantstr := pkgCfgDefWantstr
	if gotstr != wantstr {
		t.Errorf("PkgCfgDef() incorrect:\n got:%s \nwant:%s\n", gotstr, wantstr)
	}
	wiowr := fn.LogGetOutputDef()
	if giowr != wiowr {
		t.Errorf("PkgCfgDef() incorrect:\n got-iowr:%s \nwant-iowr:%s\n", giowr, wiowr)
	}

	fn.SetPkgCfgDef(true) // set pkg config to (what should be) a known default state
}

func TestPkgCfg(t *testing.T) {
	// verify individual setter funcs alter what makes up a PkgCfg

	fn.SetPkgCfgDef(true) // set pkg config to (what should be) a known default state
	fn.LogSetAlignFile(5)
	fn.LogSetAlignFunc(6)
	fn.LogSetFlags(0xffffffff)
	fn.LogSetPrefix("pRe")
	fn.LogSetTraceFlags(0xdeadbeef) //3735928559
	var wiowr io.Writer = os.Stderr
	fn.LogSetOutput(wiowr)
	got, giowr := fn.PkgCfg()
	gotstr := fmt.Sprintf("%+v", got)
	wantstr := pkgCfgDefWantstr
	if gotstr == wantstr {
		t.Errorf("PkgCfg() incorrect:\n got:%s \nwant:%s \n", gotstr, wantstr)
	}
	if giowr != wiowr {
		t.Errorf("PkgCfg() incorrect:\n gotiowr:%x \nwantowr:%x %s\n",
			giowr, wiowr, noteStdio)
	}
	fn.SetPkgCfgDef(true) // set pkg config to (what should be) a known default state
}

func TestSetPkgCfgDef(t *testing.T) {
	// verify rewriting the defaults matches what should be in the pkg defaults

	fn.SetPkgCfgDef(true) // set pkg config to (what should be) a known default state
	wiowr := fn.LogGetOutputDef()
	got, giowr := fn.PkgCfg()
	gotstr := fmt.Sprintf("%+v", got)
	wantstr := pkgCfgDefWantstr
	if gotstr != wantstr {
		t.Errorf("SetPkgCfgDef() incorrect:\n got:%s \nwant:%s %s\n", gotstr, wantstr, noteStdio)
	}
	if giowr != wiowr {
		t.Errorf("PkgCfgDef() incorrect:\n gotiowr:%x \nwantiowr:%x %s\n", giowr, wiowr, noteStdio)
	}
	fn.SetPkgCfgDef(true) // set pkg config to (what should be) a known default state
}

func Test_fetchedpkgconfig(t *testing.T) {
	//  verify setter funcs mods  match a corresponding fetched pkgCfg
	fn.SetPkgCfgDef(true) // set pkg config to (what should be) a known default state
	wiowr := io.Writer(os.Stderr)
	fn.LogSetOutput(wiowr)
	fn.LogSetFlags(0xffff)
	fn.LogSetPrefix("ZyxAbcd")
	fn.LogSetTraceFlags(0xbeef)
	fn.LogSetAlignFile(11)
	fn.LogSetAlignFunc(12)
	want := &fn.PkgCfgStruct{
		LogFlags:      0xffff,
		LogPrefix:     "ZyxAbcd",
		LogTraceFlags: 0xbeef, //dec=48879
		LogAlignFile:  11,
		LogAlignFunc:  12,
	}
	got, giowr := fn.PkgCfg()
	gotstr := fmt.Sprintf("%+v", got)
	wantstr := fmt.Sprintf("%+v", want)
	if gotstr != wantstr {
		t.Errorf("PkgCfgDef() incorrect:\n got:%s \nwant:%s \n", gotstr, wantstr)
	}
	if giowr != wiowr {
		t.Errorf("PkgCfgDef() incorrect:\n gotiowr:%x \nwantiowr:%x %s\n", giowr, wiowr, noteStdio)
	}

	fn.SetPkgCfgDef(true) // set pkg config to (what should be) a known default state
}

func Test_pkgcfg_matches_setpkgcfg(t *testing.T) {
	// verify setting a pkgcfg matches returned values from  PkgCfg

	fn.SetPkgCfgDef(true) // set pkg config to (what should be) a known default state
	want, _ := fn.PkgCfgDef()
	wiowr := fn.LogGetOutputDef()
	want.LogAlignFile = 8
	fn.SetPkgCfg(want, wiowr)
	got, giowr := fn.PkgCfg()
	gotstr := fmt.Sprintf("%+v", got)
	wantstr := fmt.Sprintf("%+v", want)
	if gotstr != wantstr {
		t.Fatalf("SetPkgCfg() incorrect:\n got:%s \nwant:%s ", gotstr, wantstr)
	}
	if giowr != wiowr {
		t.Fatalf("SetPkgCfg() incorrect:\n gotiowr:%x \nwantiowr:%x ", giowr, wiowr)
	}

	fn.SetPkgCfgDef(true) // set pkg config to (what should be) a known default state
}
