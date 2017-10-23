// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func Test_unexportFuncs(t *testing.T) {
	tests := []struct {
		name  string
		nform nameform
		want  string
	}{
		{"lvlll-fullfn", nfull, "github.com/phcurtis/fn.Test_unexportFuncs.func1"},
		{"lvlll-basefn", nbase, "fn.Test_unexportFuncs.func1"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := lvlll(0, test.nform)
			if got != test.want {
				t.Errorf("\n got:%s \nwant:%s", got, test.want)
			}
		})
	}
}

func Test_forcepanichelplt(t *testing.T) {
	if !testing.Verbose() {
		LogSetOutput(ioutil.Discard) // to hide trace output
	} else {
		LogSetOutput(os.Stderr) // so will see trace output and LogTrace stuff
	}
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard) // toss log.Panic output
	}

	defer SetPkgCfgDef(true)       // restore defaults at end of this func
	defer log.SetOutput(os.Stderr) // restore log output

	f := func() func() {
		begTime, begFn, reffile, reflnum := helpltbeg(0, LbegTraceLab, "")
		return func() {
			helpltend(0, LendTraceLab, begTime, begFn, "hack"+reffile, reflnum, "")
		}
	}

	defer func() {
		var p interface{}
		p = recover()
		if testing.Verbose() {
			log.Printf("panicErr:%v\n", p)
		}
		if p == nil {
			t.Errorf("should have paniced ... due to hacking reffile ")
		}
		if testing.Verbose() {
			log.Println("Recovered from panic")
		}
	}()

	f()() // make panic happen
}
