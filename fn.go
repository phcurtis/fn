// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fn - includes APIs relating to function names (fn).
// Such as returning a given func name relative to its position on the
// call stack. Other APIs include returning all the func names on the
// call stack, and trace logging the entry and exiting of a func including
// its time duration.
package fn

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// Version of package fn
const Version = 0.101

// Level genealogy values for exported Lvl functions
const (
	Lme     = 0       // me
	Lpar    = Lme + 1 // parent
	Lgpar   = Lme + 2 // grandparent
	Lggpar  = Lme + 3 // great-grandparent
	Lgggpar = Lme + 4 // great-great-grandparent
)

// nameform - contains form of func name to return
type nameform uint8

// list of forms of a func name to return
const (
	nfull nameform = 0 // full name form
	nbase nameform = 1 // filepath.Base form
)

const cStkEndPfix = "<EndOfCallStack:lvlll-lvl="

// low level func getting a given 'lvl' func name
func lvlll(lvl int, nform nameform) string {
	const baselvl = 2
	pc := make([]uintptr, 10)
	runtime.Callers(baselvl+lvl, pc)
	name := runtime.FuncForPC(pc[0]).Name()
	if nform == nbase {
		name = filepath.Base(name)
	}
	if name == "" {
		name = fmt.Sprintf(cStkEndPfix+"%d>", lvl)
	}
	return name
}

// Lvl - returns the func name relative to levels back on
// caller stack it was invoked from. Use lvl=Lpar for parent func,
// lvl=Lgpar or lvl=2 for GrandParent and so on.
func Lvl(lvl int) string {
	return lvlll(lvl+Lpar, nfull)
}

// LvlBase - returns the filepath.Base form of func name relative to
// levels back on caller stack it was invoked from.
func LvlBase(lvl int) string {
	return lvlll(lvl+Lpar, nbase)
}

// Cur - returns the current func name relative to where it was invoked from.
func Cur() string {
	return lvlll(Lpar, nfull)
}

// CurBase - returns the filepath.Base form of func name relative to
// where it it was invoked from.
func CurBase() string {
	return lvlll(Lpar, nbase)
}

// LvlCStkMax -- max Level call stack depth that LvlCStk will search too.
const LvlCStkMax = 500

// LvlCStk returns func names in call stack for a given level relative
// to were it was invoked from; Typically one should use CStk instead.
// Use lvl=Lpar for parent func, lvl=LgPar for GrandParent and so on
func LvlCStk(lvl int) string {
	var name, sep string
	for i := lvl; i <= LvlCStkMax; i++ {
		cname := Lvl(i + Lpar)
		if strings.HasPrefix(cname, cStkEndPfix) {
			break
		}
		name += sep + cname
		sep = "<--" // do not change - testing is dependent on this
	}
	return name
}

// CStk - returns func names in call stack relative to where it was invoked from.
func CStk() string {
	return LvlCStk(Lpar)
}
