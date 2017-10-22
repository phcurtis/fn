// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/phcurtis/fn"
)

func BenchmarkVarious(b *testing.B) {
	defer func() {
		fn.LogSetFlags(fn.LflagsDef)
		fn.LogSetOutput(os.Stdout)
	}()

	b.Run("fn.Cur............", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fn.Cur()
		}
	})
	b.Run("fn.LvlBase(Lme)...", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fn.LvlBase(fn.Lme)
		}
	})
	b.Run("fn.LvlCStk(Lme)...", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fn.LvlCStk(fn.Lme)
		}
	})
	b.Run("fn.CStk...........", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fn.CStk()
		}
	})
}

type f1Struct struct {
	cnt    int
	invoke int
	total  int
}

var f1s f1Struct

func f1(b *testing.B) {
	if f1s.cnt < f1s.invoke {
		f1s.cnt++
		f1(b)
	} else {
		deep := strings.Count(fn.CStk(), "<--") + 1
		if deep != f1s.total {
			b.Fatalf("wrong invocations: deep:%d invoke:%d total:%d", deep, f1s.invoke, f1s.total)
		}
	}
}

func f1main(total, invoke int, b *testing.B) {
	mintot := total - invoke + 3
	if total < mintot || total > fn.LvlCStkMax {
		b.Fatalf("total:%d is out of range[%d-%d]\n", total, mintot, fn.LvlCStkMax)
	}
	f1s.total = total
	f1s.invoke = invoke - 1 // since f1main is already 1 deep
	f1s.cnt = 1
	if f1s.invoke-f1s.cnt > 0 {
		f1(b)
	}
}

func Benchmark_cstkdepth(b *testing.B) {
	deepAdj := strings.Count(fn.CStk(), "<--") + 1 // add 1 to separator count
	tests := []struct {
		name string
		deep int
	}{
		{"fn.CStk.~10 deep..", 10},
		{"fn.CStk.~20 deep..", 20},
		{"fn.CStk.~30 deep..", 30},
		{"fn.CStk.~40 deep..", 40},
		{"fn.CStk.~50 deep..", 50},
		{"fn.CStk.~100deep..", 100},
		{"fn.CStk.~200deep..", 200},
		{"fn.CStk.~250deep..", 250},
		{"fn.CStk.~500deep..", 500},
	}
	for _, v := range tests {
		b.Run(v.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				invoke := v.deep - deepAdj
				f1main(v.deep, invoke, b)
			}
		})
	}
}

func routeTmpFile(b *testing.B) func() {
	tmpfile, err := ioutil.TempFile("", "fn-benchmark-")
	if err != nil {
		b.Fatal(err)
	}
	if testing.Verbose() {
		log.Printf("\nrouting output to tempfile:%s\n", tmpfile.Name())
	}
	fn.LogSetOutput(tmpfile)
	return func() {
		if testing.Verbose() {
			log.Printf("\nremoving tempfile:%s\n", tmpfile.Name())
		}
		err := os.Remove(tmpfile.Name())
		if err != nil {
			log.Printf("\nerror removing %v err:%v\n", tmpfile.Name(), err)
		}
	}
}

func align(num int, str string) string {
	size := 38
	str = fmt.Sprintf("#%02d:%s", num, str)
	if len(str) < size {
		str = str + strings.Repeat(".", size-len(str))
	}
	return str
}

func BenchmarkLog(b *testing.B) {
	fn.SetPkgCfgDef(false)

	const (
		LTF = iota
		LTFmembuf
		LTFDiscardLfdef
		LTFDiscardLfoff
		LTFTign
		LCTFYes
		LCTFNo
		LCTFYesTign

		LTMF
		LTMFmembuf
		LTMFDiscardLfdef
		LTMFDiscardLfoff
		LTMFTign
		LTFYesTign
		LCTMFYes
		LCTMFNo
		LCTMFYesTign
	)
	tests := []struct {
		num      int
		ftype    int
		logflags int
		trflags  int
	}{
		{1, LCTMFYes, fn.LflagsDef, fn.TrFlagsDef},
		{2, LCTFYes, fn.LflagsDef, fn.TrFlagsDef},

		{3, LTMF, fn.LflagsDef, fn.TrFlagsDef},
		{4, LTF, fn.LflagsDef, fn.TrFlagsDef},

		{5, LTMFmembuf, fn.LflagsDef, fn.TrFlagsDef},
		{6, LTFmembuf, fn.LflagsDef, fn.TrFlagsDef},

		{7, LTMFDiscardLfdef, fn.LflagsDef, fn.TrFlagsDef},
		{8, LTFDiscardLfdef, fn.LflagsDef, fn.TrFlagsDef},

		{9, LTMFDiscardLfoff, fn.LflagsOff, fn.TrFlagsDef},
		{10, LTFDiscardLfoff, fn.LflagsOff, fn.TrFlagsDef},

		{11, LCTMFYesTign, fn.LflagsDef, fn.TrFlagsDef | fn.Trlogignore},
		{12, LCTFYesTign, fn.LflagsDef, fn.TrFlagsDef | fn.Trlogignore},

		{13, LTMFTign, fn.LflagsDef, fn.TrFlagsDef | fn.Trlogignore},
		{14, LTFTign, fn.LflagsDef, fn.TrFlagsDef | fn.Trlogignore},

		{15, LCTMFNo, fn.LflagsDef, fn.TrFlagsDef},
		{16, LCTFNo, fn.LflagsDef, fn.TrFlagsDef},
	}
	for _, v := range tests {
		fn.SetPkgCfgDef(true) // set pkg config to default state
		fn.LogSetFlags(v.logflags)
		fn.LogSetTraceFlags(v.trflags)
		var name string
		switch v.ftype {
		case LTF, LTFDiscardLfdef, LTFDiscardLfoff, LTFTign:
			var suffix string
			if v.ftype == LTFDiscardLfdef {
				suffix = "-discard-lfdef"
				fn.LogSetOutput(ioutil.Discard)
			} else if v.ftype == LTFDiscardLfoff {
				suffix = "-discard-lfoff"
				fn.LogSetOutput(ioutil.Discard)
			} else if v.ftype == LTFTign {
				suffix = "-tign=true"
			}
			name = align(v.num, `LogTrace`+suffix)
			b.Run(name, func(b *testing.B) {
				if len(suffix) == 0 {
					defer routeTmpFile(b)()
				}
				for i := 0; i < b.N; i++ {
					fn.LogTrace()()
				}
			})

		case LTMF, LTMFDiscardLfdef, LTMFDiscardLfoff, LTMFTign:
			var suffix string
			if v.ftype == LTMFDiscardLfdef {
				suffix = "-discard-lfdef"
				fn.LogSetOutput(ioutil.Discard)
			} else if v.ftype == LTMFDiscardLfoff {
				suffix = "-discard-lfoff"
				fn.LogSetOutput(ioutil.Discard)
			} else if v.ftype == LTMFTign {
				suffix = "-tign=true"
			}
			name = align(v.num, `LogTraceMsgs`+suffix)
			b.Run(name, func(b *testing.B) {
				if len(suffix) == 0 {
					defer routeTmpFile(b)()
				}
				for i := 0; i < b.N; i++ {
					fn.LogTraceMsgs("msg1")("msg2")
				}
			})

		case LCTFNo, LCTFYes, LCTFYesTign:
			do := (v.ftype == LCTFYes) || (v.ftype == LCTFYesTign)
			name = align(v.num, fmt.Sprintf("LogCondTrace<%t>tign=%t",
				do, (v.ftype == LCTFYesTign)))
			b.Run(name, func(b *testing.B) {
				defer routeTmpFile(b)()
				for i := 0; i < b.N; i++ {
					fn.LogCondTrace(do)()
				}
			})

		case LCTMFNo, LCTMFYes, LCTMFYesTign:
			do := (v.ftype == LCTMFYes) || (v.ftype == LCTMFYesTign)
			name = align(v.num, fmt.Sprintf("LogCondTraceMsgs<%t>tign=%t",
				do, (v.ftype == LCTMFYesTign)))
			b.Run(name, func(b *testing.B) {
				defer routeTmpFile(b)()
				for i := 0; i < b.N; i++ {
					fn.LogCondTraceMsgs(do, "msg1")("msg2")
				}
			})

		case LTFmembuf:
			name = align(v.num, "LogTrace-membuf")
			buf := bytes.NewBufferString("")
			fn.LogSetOutput(buf)
			b.Run(name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					fn.LogTrace()()
				}
			})

		case LTMFmembuf:
			name = align(v.num, "LogTraceMsgs-membuf")
			buf := bytes.NewBufferString("")
			fn.LogSetOutput(buf)
			b.Run(name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					fn.LogTraceMsgs("msg1")("msg2")
				}
			})

		default:
			log.Panic("unknown switch case in: " + fn.Cur())
		}
	}

	fn.SetPkgCfgDef(true) // set pkg config to default state
}
