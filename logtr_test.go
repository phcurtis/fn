// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn_test

import (
	"bytes"
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
	const (
		LTF = iota
		LTMF
		LCTFYes
		LCTFNo
		LCTMFYes
		LCTMFNo
	)

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
		ftype   int
		regexp  bool
		arg1    string
		arg2    string
		lflags  int
		trflags int
		wantb   string
		wante   string
	}{
		{"t1", LTF, false, "", "", fn.LflagsOff, fn.TrFlagsOff,
			"LogFN: BegTrace:" + fullFN,
			"LogFN: EndTrace:" + fullFN},

		{"t2", LTMF, false, "msg1", "msg2", fn.LflagsOff, fn.TrFlagsOff,
			"LogFN: BegTrMsg:" + fullFN + " msg1",
			"LogFN: EndTrMsg:" + fullFN + " msg2"},

		{"t3", LTF, true, "", "", log.Lshortfile, fn.TrFlagsOff,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtr + fullFN + "[ ]*$",
			"LogFN: " + rexpfn1 + `:\d{1,}` + reEndtr + fullFN + "[ ]*$"},

		{"t4", LTF, true, "", "", log.Lshortfile, fn.Trfnbase | fn.Trnodur,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtr + baseFN + "[ ]*$",
			"LogFN: " + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtr + baseFN + "[ ]*$"},

		{"t5", LTF, true, "", "", log.Lshortfile, fn.Trfnbase | fn.Trnodur | fn.Trfbegrefincfile,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtr + baseFN + "[ ]*$",
			"LogFN: " + rexpfn1 + `:\d{1,}<` + rexpfn1 + `:\d{1,}>` + reEndtr + baseFN + "[ ]*$"},

		{"t6", LTF, true, "", "", log.Llongfile, fn.Trfilenogps,
			"LogFN: " + rexpfn2 + `:\d{1,}` + reBegtr + fullFN + "[ ]*$",
			"LogFN: " + rexpfn2 + `:\d{1,}<:\d{1,}>` + reEndtr + fullFN + ` .*Dur:\d{1,}[^ ]*$`},

		{"t7", LTMF, true, "msg1", "msg2", log.Lshortfile, fn.TrFlagsOff,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtrm + fullFN + " msg1$",
			"LogFN: " + rexpfn1 + `:\d{1,}` + reEndtrm + fullFN + " msg2$"},

		{"t8", LTMF, true, "msg1", "msg2", log.Ldate | log.Lshortfile, fn.TrFlagsOff,
			"LogFN: " + `.*` + reDate + ` .*` + rexpfn1 + `:\d{1,}` + reBegtrm + fullFN + " msg1$",
			"LogFN: " + `.*` + reDate + ` .*` + rexpfn1 + `:\d{1,}` + reEndtrm + fullFN + " msg2$"},

		{"t9", LTMF, true, "msg1", "msg2", log.Ldate | log.Lshortfile, fn.Trfilenogps,
			"LogFN: " + `.*` + reDate + ` .*` + rexpfn1 + `:\d{1,}` + reBegtrm + fullFN + " msg1$",
			"LogFN: " + `.*` + reDate + ` .*` + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + fullFN + reMsg2dur + `$`},

		{"tA", LTMF, true, "msg1", "msg2", log.Ldate | log.Ltime | log.Lshortfile, fn.Trfilenogps,
			"LogFN: " + `.*` + reDateTime + ` .*` + rexpfn1 + `:\d{1,}` + reBegtrm + fullFN + " msg1$",
			"LogFN: " + `.*` + reDateTime + ` .*` + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + fullFN + reMsg2dur + `$`},

		{"tB", LTMF, true, "msg1", "msg2", log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile, fn.TrFlagsDef,
			"LogFN: " + `.*` + reDateTimeMicro + ` .*` + rexpfn1 + `:\d{1,}` + reBegtrm + baseFN + " msg1$",
			"LogFN: " + `.*` + reDateTimeMicro + ` .*` + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + baseFN + reMsg2dur + `$`},

		{"tC", LTMF, true, "msg1", "msg2", log.Lshortfile, fn.TrFlagsDef,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtrm + baseFN + " msg1$",
			"LogFN: " + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + baseFN + reMsg2dur + `$`},

		{"tD", LTMF, true, "msg1", "msg2", log.Lshortfile, fn.TrFlagsDef | fn.Trbegtime,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtrm + baseFN + " msg1 .*Time:" + reDate + ` .*` + reTime + "$",
			"LogFN: " + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + baseFN + reMsg2dur + `$`},

		{"tE", LTMF, true, "msg1", "msg2", log.Lshortfile, fn.TrFlagsDef | fn.Trendtime,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtrm + baseFN + " msg1$",
			"LogFN: " + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + baseFN + reMsg2dur + ` * Time:` + reDateTime + "$"},

		{"tF", LTMF, true, "msg1", "msg2", log.Lshortfile, fn.TrFlagsDef | fn.Trbegtime | fn.Trmicroseconds,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtrm + baseFN + " msg1 .*Time:" + reDateTimeMicro + "$",
			"LogFN: " + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + baseFN + reMsg2dur + `$`},

		{"tG", LTMF, true, "msg1", "msg2", log.Lshortfile, fn.TrFlagsDef | fn.Trmicroboth,
			"LogFN: " + rexpfn1 + `:\d{1,}` + reBegtrm + baseFN + " msg1 .*Time:" + reDateTime + `\.\d{6}$`,
			"LogFN: " + rexpfn1 + `:\d{1,}<:\d{1,}>` + reEndtrm + baseFN + reMsg2dur + ` * Time:` + reDateTimeMicro + "$"},
		{"tH", LTF, false, "", "", fn.LflagsOff, fn.Trlogignore,
			"",
			""},
		{"tI", LTMF, false, "", "", fn.LflagsOff, fn.Trlogignore,
			"",
			""},

		{"tJ", LCTFYes, true, "", "", fn.LflagsOff, fn.TrFlagsDef,
			"LogFN:" + reBegtr + baseFN,
			"LogFN:" + reEndtr + baseFN},

		{"tK", LCTFNo, true, "", "", fn.LflagsOff, fn.TrFlagsDef,
			"",
			""},

		{"tL", LCTFNo, true, "", "", fn.LflagsOff, fn.Trlogignore,
			"",
			""},

		{"tM", LCTMFYes, true, "", "", fn.LflagsOff, fn.TrFlagsDef,
			"LogFN:" + reBegtrm + baseFN,
			"LogFN:" + reEndtrm + baseFN},

		{"tN", LCTMFNo, false, "", "", fn.LflagsOff, fn.TrFlagsDef,
			"",
			""},

		{"tO", LCTMFNo, false, "", "", fn.LflagsOff, fn.Trlogignore,
			"",
			""},

		{"tP", LCTMFYes, false, "", "", fn.LflagsOff, fn.Trlogignore,
			"",
			""},

		{"tQ", LCTFYes, false, "", "", fn.LflagsOff, fn.Trlogignore,
			"",
			""},
	}

	var gotb, gote string
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			fn.SetPkgCfgDef(false) // make sure we are in a known state
			fn.LogSetFlags(v.lflags)
			fn.LogSetTraceFlags(v.trflags)
			switch v.ftype {
			case LTF:
				f2 := fn.LogTrace()
				gotb = readStdoutCapLine(buf)
				f2()
				gote = readStdoutCapLine(buf)
			case LCTFYes:
				f2 := fn.LogCondTrace(true)
				gotb = readStdoutCapLine(buf)
				f2()
				gote = readStdoutCapLine(buf)
			case LCTFNo:
				f2 := fn.LogCondTrace(false)
				gotb = readStdoutCapLine(buf)
				f2()
				gote = readStdoutCapLine(buf)
			case LTMF:
				f2 := fn.LogTraceMsgs(v.arg1)
				gotb = readStdoutCapLine(buf)
				f2(v.arg2)
				gote = readStdoutCapLine(buf)
			case LCTMFYes:
				f2 := fn.LogCondTraceMsgs(true, v.arg1)
				gotb = readStdoutCapLine(buf)
				f2(v.arg2)
				gote = readStdoutCapLine(buf)
			case LCTMFNo:
				f2 := fn.LogCondTraceMsgs(false, v.arg1)
				gotb = readStdoutCapLine(buf)
				f2(v.arg2)
				gote = readStdoutCapLine(buf)
			default:
				panic("ftype not recognized")
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
