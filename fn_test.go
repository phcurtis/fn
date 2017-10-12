package fn_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/phcurtis/fn"
)

const pkgName = "fn_test"
const baseName = "github.com/phcurtis/" + pkgName + "."

func TestCur(t *testing.T) {
	gotname := fn.Cur()
	wantname := baseName + "TestCur"
	if gotname != wantname {
		t.Errorf("fn.Cur():\n got:%s \nwant:%s ", gotname, wantname)
	}
}

func TestLvl(t *testing.T) {
	gotname := fn.Lvl(0)
	wantname := baseName + "TestLvl"
	if gotname != wantname {
		t.Errorf("fn.Lvl(0):\n got:%s \nwant:%s ", gotname, wantname)
	}
}

func TestLvlBase(t *testing.T) {
	gotname := fn.LvlBase(0)
	wantname := pkgName + ".TestLvlBase"
	if gotname != wantname {
		t.Errorf("fn.LvlBase(0): got:%s \n want:%s \n", gotname, wantname)
	}
}

func TestCurBase(t *testing.T) {
	gotname := fn.CurBase()
	wantname := pkgName + ".TestCurBase"
	if gotname != wantname {
		t.Errorf("fn.CurBase(): got:%s \n want:%s \n", gotname, wantname)
	}
}

func a2() string { return fn.CStk() }
func a1() string { return a2() }

func TestCStk(t *testing.T) {
	const wantPfix = baseName + "a2<--" + baseName + "a1<--" + baseName + "TestCStk<--"

	got := a1()
	if !strings.HasPrefix(got, wantPfix) {
		t.Errorf("TestCStk.a2():\n       got:%s \nwantPrefix:%s\n", got, wantPfix)
	}
}

func b2(lvl int) string { return fn.LvlCStk(lvl) }
func b1(lvl int) string { return b2(lvl) }

func TestLvlCStk(t *testing.T) {
	tests := []struct {
		name     string
		lvl      int
		wantPfix string
	}{
		{"Lvl0", 0, baseName + "b2<--" + baseName + "b1<--" + baseName + "TestLvlCStk.func1<--"},
		{"Lvl1", 1, baseName + "b1<--" + baseName + "TestLvlCStk.func1<--"},
		{"Lvl2", 2, baseName + "TestLvlCStk.func1<--"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := b1(test.lvl)
			if !strings.HasPrefix(got, test.wantPfix) {
				t.Errorf("TestLvlCStk.b2(%d):\n     got:%s \nwantPfix:%s ", test.lvl, got, test.wantPfix)
			}
		})
	}
}

func TestLogSetPrefix(t *testing.T) {
	tests := []struct {
		name string
		set  string
		cmp  string
		res  bool
	}{
		{"t1", "LogFn1:", "LogFn1:", true},
		{"t2", "XogFn1:", "XogFn1:", true},
		{"t3", "LogFn1:", "XogFn1:", false},
		{"t4", "LogFn:", "XogFn1:", false},
		{"t5", "LogFn1:", "LogFn1:", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn.LogSetPrefix(test.set)
			got := fn.LogPrefix()
			if got == test.cmp != test.res {
				t.Errorf("LogSetPrefix(%q) == %s expectedResult:%t\n",
					test.set, test.cmp, test.res)
			}
		})
	}
}

func readStdoutCapLine(b *bytes.Buffer) string {
	const lfeed = '\n'
	line, err := b.ReadString(lfeed)
	if err != nil {
		if err == io.EOF {
			return ""
		}
		panic(err)
	}
	lena := len(line) - 1
	if lena >= 0 && line[lena] == lfeed {
		line = line[:lena]
	}
	return line
}

func Test_logfuncs(t *testing.T) {
	fullFuncName := baseName + "Test_logfuncs"

	fn.LogSetFlags(fn.LflagsOff) // turn off stuff so don't have to fool with time, etc
	fn.LogSetPrefix("Prefix:")   // set prefix to known value

	// set fn.Log outputs to be sent to bytes.Buffer (in memory)
	buf := bytes.NewBufferString("")
	fn.LogSetOutput(buf)

	// capture to buf: LogBeg/LogEnd both log.logger stuff and func returned stuff
	logbeg := fn.LogBeg()
	fmt.Fprintln(buf, logbeg)
	logend := fn.LogEnd(logbeg)
	fmt.Fprintln(buf, logend)

	// capture to buf: LogBegDur/LogEndDur both log.logger stuff and func returned stuff
	logbegDur, tim := fn.LogBegDur()
	fmt.Fprintln(buf, logbegDur)
	logendDur, dur := fn.LogEndDur(logbegDur, tim)
	fmt.Fprintf(buf, "%s Dur:%v\n", logendDur, dur)

	tests := []struct {
		name string
		dur  bool
		want string
	}{
		{"LogBeg()----Log.", false, "Prefix:Beg:"},
		{"LogBeg()----Func", false, ""},
		{"LogEnd()----Log.", false, "Prefix:End:"},
		{"LogEnd()----Func", false, ""},
		{"LogBegDur()-Log.", false, "Prefix:BegDur:"},
		{"LogBegDur()-Func", false, ""},
		{"LogEndDur()-Log.", true, "Prefix:EndDur:"},
		{"LogEndDur()-Func", true, ""},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			soc := readStdoutCapLine(buf)
			if v.dur {
				want := v.want + fullFuncName + " Dur:"
				if !strings.HasPrefix(soc, want) {
					t.Errorf("%s: \n     got:%s \nwantPfix:%s ", v.name, soc, want)
				}
				lenw := len(want)
				var dur, dur2 time.Duration
				var err error
				dur, err = time.ParseDuration(soc[lenw:])
				if err != nil {
					t.Errorf("%s: DurBad: got:%q", v.name, soc[lenw:])
				} else {
					maxDur := "1s"
					if dur2, err = time.ParseDuration(maxDur); err != nil {
						panic(err)
					}
					if dur > dur2 {
						t.Errorf("%s: DurTooLong: got:%q max:%s", v.name, soc[lenw:], maxDur)
					}
				}

			} else {
				want := v.want + fullFuncName
				if soc != want {
					t.Errorf("%s: \n got:%s \nwant:%s ", v.name, soc, want)
				}
			}
		})
	}
}
