// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func strMinWidth(str string, width int) string {
	len1 := len(str)
	if len1 < width {
		str = str + strings.Repeat(" ", width-len1)
	}
	return str
}

// fn log trace but may first need to find appropriate filename and line num
func helplt(lvl int, msg, reffile, reflnum string) (newReffile, newReflnum string) {
	var filenlr string

	// get original [current] log flags
	orgflags := logt.Flags()
	sl := log.Lshortfile | log.Llongfile
	lfn := orgflags & sl

	var ref string

	// if log flags are including filename
	if lfn > 0 {
		_, file, line, _ := runtime.Caller(lvl)
		linenum := fmt.Sprintf(":%d", line)
		if orgflags&log.Lshortfile > 0 {
			file = filepath.Base(file)
		} else {
			// log.Llongfile
			if logTraceFlags&Trfilenogps > 0 {
				if strings.HasPrefix(file, gopathsrc) {
					file = file[len(gopathsrc):]
				}
			}
		}
		newReffile = filepath.Base(file)
		newReflnum = linenum

		// seems if this is true there is a problem elsewhere
		// as in a weird invocation end portion of log trace,
		// possible by passing that portion and invoking in
		// another function which package fn does NOT support.
		if reffile != "" && reffile != newReffile {
			logt.Panic(errors.New("reffile:" + reffile + " != newReffile:" + newReffile))
		}

		if reffile != "" && logTraceFlags&Trfnobegref == 0 {
			if logTraceFlags&Trfbegrefincfile > 0 {
				ref = "<" + reffile + reflnum + ">"
			} else {
				ref = "<" + reflnum + ">"
			}
		}

		filenlr = file + linenum + ref + " "
		filenlr = strMinWidth(filenlr, logAlignFile)

		// set log flags not to include filename
		logt.SetFlags(orgflags &^ sl)
	}
	logt.Printf("%s%s", filenlr, msg)
	if lfn > 0 {
		// restore log flags
		logt.SetFlags(orgflags)
	}
	return newReffile, newReflnum
}

func formatTime(t time.Time, microseconds bool, msg string) string {
	yr, mon, dy := t.Date()
	hr, min, sec := t.Clock()
	var micro string
	if microseconds {
		micro = fmt.Sprintf(".%06d", t.Nanosecond()/1e3)
	}
	return fmt.Sprintf("%s%d/%02d/%02d %02d:%02d:%02d%s", msg, yr, mon, dy, hr, min, sec, micro)
}

func helpltend(lvladj int, trlabel string, start time.Time, begFn, reffile, reflnum, endMsg string) {
	endTime := time.Now()
	endFn := Lvl(Lgpar + lvladj)
	if begFn != endFn {
		if strings.Contains(CStk(), "<--runtime.gopanic") {
			logt.Println("GOPANIC DETECTED --exiting '"+trlabel+"'(helpltend)>CStk:", CStk())
			logt.Println("begFn:"+begFn+" != endFn:"+endFn, " reffile:", reffile, " reflnum", reflnum, "\n\n ")
			return
		}
		// if Idiomatic usage of LogTrace and LogTraceMsgs then should not have a panic.
		err := fmt.Sprintf("begFn != endFn\n begFn:%s\n endFn:%s\n  Cstk:%s \n"+
			"Panic probable cause due to end trace pairing return portion called from different func",
			begFn, endFn, CStk())
		logt.Panic(errors.New(err)) // see todo above
	}
	if endMsg != "" {
		endMsg = " " + endMsg
	}
	var str string

	muLogt.Lock()
	defer muLogt.Unlock()

	if logTraceFlags&Trfnbase > 0 {
		str = filepath.Base(endFn)
	} else {
		str = endFn
	}

	str = strMinWidth(str, logAlignFunc) + endMsg

	if logTraceFlags&Trnodur == 0 {
		str += " Dur:" + endTime.Sub(start).Round(time.Microsecond).String()
	}
	if logTraceFlags&Trendtime > 0 {
		str += formatTime(endTime, logTraceFlags&Trmicroseconds > 0, " Time:")
	}
	helplt(3+lvladj, trlabel+str, reffile, reflnum)
}

func helpltbeg(lvladj int, trlabel string, begMsg string) (begTime time.Time, begFn, reffile, reflnum string) {
	begTime = time.Now()
	begFn = Lvl(Lgpar + lvladj)
	if begMsg != "" {
		begMsg = " " + begMsg
	}
	var str string
	muLogt.Lock()
	defer muLogt.Unlock()

	if logTraceFlags&Trfnbase > 0 {
		str = filepath.Base(begFn)
	} else {
		str = begFn
	}

	str = strMinWidth(str, logAlignFunc) + begMsg

	if logTraceFlags&Trbegtime > 0 {
		str += formatTime(begTime, logTraceFlags&Trmicroseconds > 0, " Time:")
	}
	reffile, reflnum = helplt(3, trlabel+str, "", "")
	return begTime, begFn, reffile, reflnum
}

// LogTrace - log the begin tracing portion of the current function name
// [adjusting output according to the configuration settings such as Trace Flags,
// stdlib log, etc. at the time of its execution] and return a pairing end func
// that must be called within the same func.
//	Idiomatic usage at func start: defer fn.LogTrace()()
//  NOTE that the pairing return function uses the configurations that are at
//  its time of execution which [you] may have changed since the begin portion.
//  Also see LogCondTrace.
func LogTrace() func() {
	muLogt.Lock()
	if logTraceFlags&Trlogignore > 0 {
		muLogt.Unlock()
		return func() {}
	}
	muLogt.Unlock()

	begTime, begFn, reffile, reflnum := helpltbeg(0, LbegTraceLab, "")
	return func() {
		helpltend(0, LendTraceLab, begTime, begFn, reffile, reflnum, "")
	}
}

// LogCondTrace - conditional version of LogTrace.
// cond - if true call LogTrace.
func LogCondTrace(cond bool) func() {
	if !cond {
		return func() {}
	}

	muLogt.Lock()
	if logTraceFlags&Trlogignore > 0 {
		muLogt.Unlock()
		return func() {}
	}
	muLogt.Unlock()

	begTime, begFn, reffile, reflnum := helpltbeg(0, LbegTraceLab, "")
	return func() {
		helpltend(0, LendTraceLab, begTime, begFn, reffile, reflnum, "")
	}
}

// LogCondMsg - logs message if cond true.
func LogCondMsg(cond bool, msg string) {
	if !cond {
		return
	}

	muLogt.Lock()
	// may not want this check
	if logTraceFlags&Trlogignore > 0 {
		muLogt.Unlock()
		return
	}
	muLogt.Unlock()

	helpltbeg(0, "Msg:", msg)
}

// LogCondTraceMsgs - conditional version of LogTraceMsgs.
//	cond - if true call LogTraceMsgs.
func LogCondTraceMsgs(cond bool, begMsg string) func(endMsg string) {
	if !cond {
		return func(string) {}
	}

	muLogt.Lock()
	if logTraceFlags&Trlogignore > 0 {
		muLogt.Unlock()
		return func(string) {}
	}
	muLogt.Unlock()

	begTime, begFn, reffile, reflnum := helpltbeg(0, LbegTraceMsgsLab, begMsg)
	return func(endMsg string) {
		helpltend(0, LendTraceMsgsLab, begTime, begFn, reffile, reflnum, endMsg)
	}
}

// LogTraceMsgs - log the begin tracing portion of the current function name
// [adjusting output according to the configuration settings such as Trace Flags,
// stdlib log, etc. at the time of its execution] and return a pairing end func
// that must be called within the same func.
//	begMsg - printed after funcname when the LogTraceMsgs is invoked
//	endMsg - printed after funcname when the LogTraceMsgs returned func is invoked.
//	Idiomatic usage at func start: defer fn.LogTraceMsgs("begMsg")("endMsg")
//  NOTE that the pairing return function uses the configurations that are at
//  its time of execution which [you] may have changed since the begin portion.
//  Also see LogTraceMsgp, LogCondTraceMsgs, LogCondTraceMsgp.
func LogTraceMsgs(begMsg string) func(endMsg string) {
	muLogt.Lock()
	if logTraceFlags&Trlogignore > 0 {
		muLogt.Unlock()
		return func(string) {}
	}
	muLogt.Unlock()

	begTime, begFn, reffile, reflnum := helpltbeg(0, LbegTraceMsgsLab, begMsg)
	return func(endMsg string) {
		helpltend(0, LendTraceMsgsLab, begTime, begFn, reffile, reflnum, endMsg)
	}
}

// LogTraceMsgp - same as LogTraceMsgs however endMsg is a pointer to a string.
func LogTraceMsgp(begMsg string) func(endMsg *string) {
	muLogt.Lock()
	if logTraceFlags&Trlogignore > 0 {
		muLogt.Unlock()
		return func(*string) {}
	}
	muLogt.Unlock()

	begTime, begFn, reffile, reflnum := helpltbeg(0, LbegTraceMsgpLab, begMsg)
	return func(endMsg *string) {
		helpltend(0, LendTraceMsgpLab, begTime, begFn, reffile, reflnum, *endMsg)
	}
}

// LogCondTraceMsgp - same as LogCondTraceMsgs however endMsg is a pointer to a string.
func LogCondTraceMsgp(cond bool, begMsg string) func(endMsg *string) {
	if !cond {
		return func(*string) {}
	}

	muLogt.Lock()
	if logTraceFlags&Trlogignore > 0 {
		muLogt.Unlock()
		return func(*string) {}
	}
	muLogt.Unlock()

	begTime, begFn, reffile, reflnum := helpltbeg(0, LbegTraceMsgpLab, begMsg)
	return func(endMsg *string) {
		helpltend(0, LendTraceMsgpLab, begTime, begFn, reffile, reflnum, *endMsg)
	}
}
