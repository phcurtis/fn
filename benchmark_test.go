// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn_test

import (
	"bytes"
	"fmt"
	"io"
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
	mintot := 4
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
	outdef := fn.LogGetOutputDef()
	outdis := ioutil.Discard
	membuf := bytes.NewBufferString("")

	tests := []struct {
		num      int
		ftype    int
		logflags int
		trflags  int
		suffix   string
		iowr     io.Writer
	}{
		{1, LCTMFYes, fn.LflagsDef, fn.TrFlagsDef, "<true>tign=false", outdef},
		{2, LCTFYes, fn.LflagsDef, fn.TrFlagsDef, "<true>tign=false", outdef},

		{3, LTMF, fn.LflagsDef, fn.TrFlagsDef, "", outdef},
		{4, LTF, fn.LflagsDef, fn.TrFlagsDef, "", outdef},

		{5, LTMFmembuf, fn.LflagsDef, fn.TrFlagsDef, "-membuf", membuf},
		{6, LTFmembuf, fn.LflagsDef, fn.TrFlagsDef, "-membuf", membuf},

		{7, LTMFDiscardLfdef, fn.LflagsDef, fn.TrFlagsDef, "-discard-lfdef", outdis},
		{8, LTFDiscardLfdef, fn.LflagsDef, fn.TrFlagsDef, "-discard-lfdef", outdis},

		{9, LTMFDiscardLfoff, fn.LflagsOff, fn.TrFlagsDef, "-discard-lfoff", outdis},
		{10, LTFDiscardLfoff, fn.LflagsOff, fn.TrFlagsDef, "-discard-lfoff", outdis},

		{11, LCTMFYesTign, fn.LflagsDef, fn.TrFlagsDef | fn.Trlogignore, "<true>tign=true", outdef},
		{12, LCTFYesTign, fn.LflagsDef, fn.TrFlagsDef | fn.Trlogignore, "<true>tign=true", outdef},

		{13, LTMFTign, fn.LflagsDef, fn.TrFlagsDef | fn.Trlogignore, "-tign=true", outdef},
		{14, LTFTign, fn.LflagsDef, fn.TrFlagsDef | fn.Trlogignore, "-tign=true", outdef},

		{15, LCTMFNo, fn.LflagsDef, fn.TrFlagsDef, "<false>-tign=false", outdef},
		{16, LCTFNo, fn.LflagsDef, fn.TrFlagsDef, "<false>-tign=false", outdef},
	}
	for _, v := range tests {
		fn.SetPkgCfgDef(true) // set pkg config to default state
		fn.LogSetFlags(v.logflags)
		fn.LogSetTraceFlags(v.trflags)
		fn.LogSetOutput(v.iowr)
		var name string
		switch v.ftype {
		case LTF, LTFDiscardLfdef, LTFDiscardLfoff, LTFTign, LTFmembuf:
			name = align(v.num, "LogTrace"+v.suffix)
			b.Run(name, func(b *testing.B) {
				if v.iowr == outdef {
					defer routeTmpFile(b)()
				}
				for i := 0; i < b.N; i++ {
					fn.LogTrace()()
				}
			})

		case LTMF, LTMFDiscardLfdef, LTMFDiscardLfoff, LTMFTign, LTMFmembuf:
			name = align(v.num, "LogTraceMsgs"+v.suffix)
			b.Run(name, func(b *testing.B) {
				if v.iowr == outdef {
					defer routeTmpFile(b)()
				}
				for i := 0; i < b.N; i++ {
					fn.LogTraceMsgs("msg1")("msg2")
				}
			})

		case LCTFNo, LCTFYes, LCTFYesTign:
			name = align(v.num, "LogCondTrace"+v.suffix)
			b.Run(name, func(b *testing.B) {
				defer routeTmpFile(b)()
				for i := 0; i < b.N; i++ {
					fn.LogCondTrace(v.ftype != LCTFNo)()
				}
			})

		case LCTMFNo, LCTMFYes, LCTMFYesTign:
			name = align(v.num, "LogCondTraceMsgs"+v.suffix)
			b.Run(name, func(b *testing.B) {
				defer routeTmpFile(b)()
				for i := 0; i < b.N; i++ {
					fn.LogCondTraceMsgs(v.ftype != LCTMFNo, "msg1")("msg2")
				}
			})

		default:
			log.Panic("unknown switch case in: " + fn.Cur())
		}
	}

	fn.SetPkgCfgDef(true) // set pkg config to default state
}
