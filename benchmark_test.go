// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/phcurtis/fn"
)

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

func BenchmarkLog(b *testing.B) {
	fn.SetPkgCfgDef(false)
	fn.LogSetTraceFlags(fn.TrFlagsDef)
	func1 := func(lflags int, tempbn string) {
		fn.LogSetFlags(lflags)
		tmpfile, err := ioutil.TempFile("", tempbn)
		if err != nil {
			b.Fatal(err)
		}
		if testing.Verbose() {
			log.Printf("\nrouting output to tempfile:%s\n", tmpfile.Name())
		}
		fn.LogSetOutput(tmpfile)
		defer func() {
			if testing.Verbose() {
				log.Printf("\nremoving tempfile:%s\n", tmpfile.Name())
			}
			err := os.Remove(tmpfile.Name())
			if err != nil {
				log.Printf("\nerror removing %v err:%v\n", tmpfile.Name(), err)
			}
		}()
	}

	f := func(str string) string {
		size := 38
		if len(str) < size {
			str = str + strings.Repeat(".", size-len(str))
		}
		return str
	}

	b.Run(f("LogTrace()()"), func(b *testing.B) {
		func1(fn.LflagsDef, "logTrace-")
		for i := 0; i < b.N; i++ {
			// don't do a defer HERE because how benchmark apparatus works
			//   -- it defers all the b.N defer calls until entire loop finishes
			//   which makes timing huge.
			//
			fn.LogTrace()()
		}
	})
	b.Run(f(`LogTraceMsgs("msg1")("msg2")`), func(b *testing.B) {
		func1(fn.LflagsDef, "logTrace-")
		for i := 0; i < b.N; i++ {
			// don't use defer HERE see comment above
			fn.LogTraceMsgs("msg1")("msg2")
		}
	})

	b.Run(f("LogTrace()()-Discard"), func(b *testing.B) {
		fn.LogSetOutput(ioutil.Discard)
		fn.LogSetFlags(fn.LflagsOff)
		for i := 0; i < b.N; i++ {
			// don't use defer HERE see comment above
			fn.LogTrace()()
		}
	})
	b.Run(f(`LogTraceMsgs("msg1")("msg2")-Discard`), func(b *testing.B) {
		fn.LogSetOutput(ioutil.Discard)
		fn.LogSetFlags(fn.LflagsOff)
		for i := 0; i < b.N; i++ {
			// don't use defer HERE see comment above
			fn.LogTraceMsgs("msg1")("msg2")
		}
	})

	b.Run(f(`LogTraceMsgs("msg1")("msg2")-Trlogoff`), func(b *testing.B) {
		fn.LogSetTraceFlags(fn.TrFlagsDef | fn.Trlogoff)
		fn.LogSetFlags(fn.LflagsOff)
		for i := 0; i < b.N; i++ {
			// don't use defer HERE see comment above
			fn.LogTraceMsgs("msg1")("msg2")
		}
	})

	b.Run(f(`LogTraceMsgs("msg1")("msg2")-toMembuf`), func(b *testing.B) {
		fn.LogSetTraceFlags(fn.TrFlagsDef)
		buf := bytes.NewBufferString("")
		fn.LogSetOutput(buf)
		fn.LogSetFlags(fn.LflagsDef)

		//func1(fn.LflagsDef, "logTrace-")
		for i := 0; i < b.N; i++ {
			// don't use defer HERE see comment above
			fn.LogTraceMsgs("msg1")("msg2")
		}
	})
}
