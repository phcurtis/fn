// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"testing"

	"github.com/phcurtis/fn"
)

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
	funcname := "Test_logfuncs"
	anonfuncname := "Test_logfuncs.func1"
	baseFN := pkgName + "." + anonfuncname
	fullFN := baseName + anonfuncname
	rexpfn1 := `logtr_test\.go`
	rexpfn2 := `github.com/phcurtis/fn/logtr_test.go`

	// prelims
	// set fn.Log outputs to be sent to bytes.Buffer (in memory)
	buf := bytes.NewBufferString("")
	fn.LogSetOutput(buf)
	ltf := fn.LogTrace
	ltmf := fn.LogTraceMsgs

	reDate := `\d{4}/\d\d/\d\d`
	reTime := `\d\d:\d\d:\d\d`
	reDateTime := reDate + ` .*` + reTime
	reDateTimeMicro := reDateTime + `\.\d{6}`
	reBegtr := `[ ]* BegTrace:`
	reBegtrm := `[ ]* BegTrMsg:`
	reEndtr := `[ ]* EndTrace:`
	reEndtrm := `[ ]* EndTrMsg:`
	reMsg2dur := ` msg2 Dur:\d{1,}[^ ]*`
	tests := []struct {
		name    string
		fnt     interface{}
		regexp  bool
		arg1    string
		arg2    string
		lflags  int
		trflags int
		wantb   string
		wante   string
	}{
		{"t1", ltf, false, "", "", fn.LflagsOff, fn.TrFlagsOff,
			"LogFN: BegTrace:" + fullFN,
			"LogFN: EndTrace:" + fullFN},

		{"t2", ltmf, false, "msg1", "msg2", fn.LflagsOff, fn.TrFlagsOff,
			"LogFN: BegTrMsg:" + fullFN + " msg1",
			"LogFN: EndTrMsg:" + fullFN + " msg2"},

		{"t3", ltf, true, "", "", log.Lshortfile, fn.TrFlagsOff,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtr + fullFN + "[ ]*$",
			"LogFN: " + rexpfn1 + `:\d{1,}` + reEndtr + fullFN + "[ ]*$"},

		{"t4", ltf, true, "", "", log.Lshortfile, fn.Trfnbase | fn.Trnodur,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtr + baseFN + "[ ]*$",
			"LogFN: " + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtr + baseFN + "[ ]*$"},

		{"t5", ltf, true, "", "", log.Lshortfile, fn.Trfnbase | fn.Trnodur | fn.Trfbegrefincfile,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtr + baseFN + "[ ]*$",
			"LogFN: " + rexpfn1 + `:\d{1,}<` + rexpfn1 + `:\d{1,}>` + reEndtr + baseFN + "[ ]*$"},

		{"t6", ltf, true, "", "", log.Llongfile, fn.Trfilenogps,
			"LogFN: " + rexpfn2 + `:\d{1,}` + reBegtr + fullFN + "[ ]*$",
			"LogFN: " + rexpfn2 + `:\d{1,}<:\d{1,}>` + reEndtr + fullFN + ` .*Dur:\d{1,}[^ ]*$`},

		{"t7", ltmf, true, "msg1", "msg2", log.Lshortfile, fn.TrFlagsOff,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtrm + fullFN + " msg1$",
			"LogFN: " + rexpfn1 + `:\d{1,}` + reEndtrm + fullFN + " msg2$"},

		{"t8", ltmf, true, "msg1", "msg2", log.Ldate | log.Lshortfile, fn.TrFlagsOff,
			"LogFN: " + `.*` + reDate + ` .*` + rexpfn1 + `:\d{1,}` + reBegtrm + fullFN + " msg1$",
			"LogFN: " + `.*` + reDate + ` .*` + rexpfn1 + `:\d{1,}` + reEndtrm + fullFN + " msg2$"},

		{"t9", ltmf, true, "msg1", "msg2", log.Ldate | log.Lshortfile, fn.TrFlagsDef,
			"LogFN: " + `.*` + reDate + ` .*` + rexpfn1 + `:\d{1,}` + reBegtrm + fullFN + " msg1$",
			"LogFN: " + `.*` + reDate + ` .*` + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + fullFN + reMsg2dur + `$`},

		{"tA", ltmf, true, "msg1", "msg2", log.Ldate | log.Ltime | log.Lshortfile, fn.TrFlagsDef,
			"LogFN: " + `.*` + reDateTime + ` .*` + rexpfn1 + `:\d{1,}` + reBegtrm + fullFN + " msg1$",
			"LogFN: " + `.*` + reDateTime + ` .*` + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + fullFN + reMsg2dur + `$`},

		{"tB", ltmf, true, "msg1", "msg2", log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile, fn.TrFlagsDef,
			"LogFN: " + `.*` + reDateTimeMicro + ` .*` + rexpfn1 + `:\d{1,}` + reBegtrm + fullFN + " msg1$",
			"LogFN: " + `.*` + reDateTimeMicro + ` .*` + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + fullFN + reMsg2dur + `$`},

		{"tC", ltmf, true, "msg1", "msg2", log.Lshortfile, fn.TrFlagsDef,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtrm + fullFN + " msg1$",
			"LogFN: " + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + fullFN + reMsg2dur + `$`},

		{"tD", ltmf, true, "msg1", "msg2", log.Lshortfile, fn.TrFlagsDef | fn.Trbegtime,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtrm + fullFN + " msg1 .*Time:" + reDate + ` .*` + reTime + "$",
			"LogFN: " + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + fullFN + reMsg2dur + `$`},

		{"tE", ltmf, true, "msg1", "msg2", log.Lshortfile, fn.TrFlagsDef | fn.Trendtime,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtrm + fullFN + " msg1$",
			"LogFN: " + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + fullFN + reMsg2dur + ` * Time:` + reDateTime + "$"},

		{"tF", ltmf, true, "msg1", "msg2", log.Lshortfile, fn.TrFlagsDef | fn.Trbegtime | fn.Trmicroseconds,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtrm + fullFN + " msg1 .*Time:" + reDateTimeMicro + "$",
			"LogFN: " + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + fullFN + reMsg2dur + `$`},

		{"tG", ltmf, true, "msg1", "msg2", log.Lshortfile, fn.TrFlagsDef | fn.Trmicroboth,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtrm + fullFN + " msg1 .*Time:" + reDateTime + `\.\d{6}$`,
			"LogFN: " + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + fullFN + reMsg2dur + ` * Time:` + reDateTimeMicro + "$"},

		{"tH", ltf, false, "", "", fn.LflagsOff, fn.Trlogoff,
			"",
			""},
		{"tI", ltmf, false, "", "", fn.LflagsOff, fn.Trlogoff,
			"",
			""},
	}

	var gotb, gote string
	for i, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			fn.SetPkgCfgDef(false) // make sure we are in a known state
			fn.LogSetFlags(v.lflags)
			fn.LogSetTraceFlags(v.trflags)
			switch f := v.fnt.(type) {
			case (func() func()):
				f2 := f()
				gotb = readStdoutCapLine(buf)
				f2()
				gote = readStdoutCapLine(buf)
			case (func(string) func(string)):
				f2 := f(v.arg1)
				gotb = readStdoutCapLine(buf)
				f2(v.arg2)
				gote = readStdoutCapLine(buf)
			default:
				t.Fatalf("item in list is unsupport func type, i:=%d\n", i)
				fmt.Printf("i:%d v:%#v\n", i, v)
			}
			if testing.Verbose() {
				log.SetFlags(0)
				log.Printf("%[1]s:%[2]s  gotb:%[3]s \n%[1]s:%[2]s wantb:%[4]s \n", funcname, v.name, gotb, v.wantb)
				log.Printf("%[1]s:%[2]s  gote:%[3]s \n%[1]s:%[2]s wante:%[4]s \n", funcname, v.name, gote, v.wante)
			}
			if v.regexp {
				re := regexp.MustCompile(v.wantb)
				if !re.MatchString(gotb) {
					t.Errorf("%s:b \n got:%s \nwant:%s \n", v.name, gotb, v.wantb)
				}
				re = regexp.MustCompile(v.wante)
				if !re.MatchString(gote) {
					t.Errorf("%s:e \n got:%s \nwant:%s \n", v.name, gote, v.wante)
				}
			} else {
				if gotb != v.wantb {
					t.Errorf("%s \n got:%s \nwant:%s \n", v.name, gotb, v.wantb)
				}
				if gote != v.wante {
					t.Errorf("%s \n got:%s \nwant:%s \n", v.name, gote, v.wante)
				}
			}
		})
	}
}
