package fn_test

import (
	"fmt"
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

func f1() {
	if f1s.cnt < f1s.invoke {
		f1s.cnt++
		f1()
	} else {
		fn.CStk()
		deep := strings.Count(fn.CStk(), "<--") + 1
		if deep != f1s.total {
			panic(fmt.Sprintf("wrong invocations: deep:%d invoke:%d total:%d", deep, f1s.invoke, f1s.total))
		}
	}
}
func f1main(total, invoke int) {
	if invoke < 0 || invoke > fn.MaxLvlCStk {
		panic(fmt.Sprintf("invoke is out of range:%d\n", invoke))
	}
	f1s.total = total
	f1s.invoke = invoke - 1 // since f1main is already 1 deep
	f1s.cnt = 1
	if f1s.invoke-f1s.cnt > 0 {
		f1()
	}
}
func BenchmarkVarious(b *testing.B) {
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
	}
	for _, v := range tests {
		b.Run(v.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				invoke := v.deep - deepAdj
				f1main(v.deep, invoke)
			}
		})
	}

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
	b.Run("LogEnd(LogBeg())..............", func(b *testing.B) {
		func1(fn.LflagsDef, "logEndLogBeg-")
		for i := 0; i < b.N; i++ {
			/* do not do a defer because how benchmark apparatus works
			   it defers all the b.N defer calls until entire loop finishes.
			*/
			fn.LogEnd(fn.LogBeg())
		}
	})
	b.Run("LogEnd(LogBeg())-Discard......", func(b *testing.B) {
		fn.LogSetOutput(ioutil.Discard)
		fn.LogSetFlags(fn.LflagsOff)
		for i := 0; i < b.N; i++ {
			// don't use defer see comment above
			fn.LogEnd(fn.LogBeg())
		}
	})

	b.Run("LogEndDur(LogBegDur())........", func(b *testing.B) {
		func1(fn.LflagsDef, "logEndDurLogBegDur-")
		for i := 0; i < b.N; i++ {
			// don't use defer see comment above
			fn.LogEndDur(fn.LogBegDur())
		}
	})
	b.Run("LogEndD(LogBegD())-Discard....", func(b *testing.B) {
		fn.LogSetOutput(ioutil.Discard)
		fn.LogSetFlags(fn.LflagsOff)
		for i := 0; i < b.N; i++ {
			// don't use defer see comment above
			fn.LogEndDur(fn.LogBegDur())
		}
	})
}
