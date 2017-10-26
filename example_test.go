// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/phcurtis/fn"
)

func ExampleCur() {
	fmt.Println("My FuncName:", fn.Cur())
	// Output:
	// My FuncName: github.com/phcurtis/fn_test.ExampleCur
}

func ExampleLvl() {
	func() {
		fmt.Println("fn.Lvl(Lme)     [func1]:", fn.Lvl(fn.Lme))
		fmt.Println("fn.Lvl(fn.Lpar) [func1]:", fn.Lvl(fn.Lpar))
		func() {
			fmt.Println(" fn.Cur()           [func1.1]:", fn.Cur())
			fmt.Println(" fn.Lvl(fn.Lpar))   [func1.1]:", fn.Lvl(fn.Lpar))
			fmt.Println(" fn.Lvl(fn.Lgpar))) [func1.1]:", fn.Lvl(fn.Lgpar))
		}()
	}()
	// Output:
	// fn.Lvl(Lme)     [func1]: github.com/phcurtis/fn_test.ExampleLvl.func1
	// fn.Lvl(fn.Lpar) [func1]: github.com/phcurtis/fn_test.ExampleLvl
	//  fn.Cur()           [func1.1]: github.com/phcurtis/fn_test.ExampleLvl.func1.1
	//  fn.Lvl(fn.Lpar))   [func1.1]: github.com/phcurtis/fn_test.ExampleLvl.func1
	//  fn.Lvl(fn.Lgpar))) [func1.1]: github.com/phcurtis/fn_test.ExampleLvl
}

func ExampleLvlBase() {
	func() {
		func() {
			fmt.Println("fn.LvlBase(0)          [func1.1]:", fn.LvlBase(0))
			fmt.Println("fn.LvlBase(fn.Lpar))   [func1.1]:", fn.LvlBase(fn.Lpar))
			fmt.Println("fn.LvlBase(fn.Lgpar))) [func1.1]:", fn.LvlBase(fn.Lgpar))
		}()
	}()
	// Output:
	// fn.LvlBase(0)          [func1.1]: fn_test.ExampleLvlBase.func1.1
	// fn.LvlBase(fn.Lpar))   [func1.1]: fn_test.ExampleLvlBase.func1
	// fn.LvlBase(fn.Lgpar))) [func1.1]: fn_test.ExampleLvlBase
}

func ExampleCurBase() {
	fmt.Println("fn.CurBase():", fn.CurBase())
	// Output:
	// fn.CurBase(): fn_test.ExampleCurBase
}

func ExampleCStk() {
	const baseName = "github.com/phcurtis/fn_test."
	wantPfix := baseName + "ExampleCStk.func1<--" + baseName + "ExampleCStk<--"

	got := func() string {
		return fn.CStk()
	}()
	if strings.HasPrefix(got, wantPfix) {
		fmt.Println(wantPfix)
	} else {
		fmt.Println(got)
	}
	// Output: github.com/phcurtis/fn_test.ExampleCStk.func1<--github.com/phcurtis/fn_test.ExampleCStk<--
}

// Examples of Idiomatic usages of both LogTrace and LogTraceMsgs
func Example_logtrace() {
	// output below: current func name:fn_test.Example_logtrace
	defer fn.LogTrace()() // this is line 77 see output below
	fmt.Println("Hi There Gopher")
	time.Sleep(time.Second)
	defer fn.LogTraceMsgs("message1")("message2") // this is line 80 and (line 83 WAS ending func brace)
	fmt.Println("Goodbye Gopher")

	/* Representative Output follows with lines wrapped and indented:
	   LogFN: 2017/10/21 13:24:10 example_test.go:77      		<wrapped-indented>
	   		BegTrace:fn_test.Example_logtrace
	   Hi There Gopher
	   LogFN: 2017/10/21 13:24:11 example_test.go:80      		<wrapped-indented>
	   		BegTrMsg:fn_test.Example_logtrace message1
	   Goodbye Gopher
	   LogFN: 2017/10/21 13:24:11 example_test.go:83<:80> 		<wrapped-indented>
	   		EndTrMsg:fn_test.Example_logtrace message2 Dur:301µs
	   LogFN: 2017/10/21 13:24:11 example_test.go:83<:77> 		<wrapped-indented>
	   		EndTrace:fn_test.Example_logtrace Dur:1.000603s
	*/
}

func ExampleLvlInfo() {
	var file string
	var line int
	var name string
	file, line, name = fn.LvlInfo(0, fn.IflagsDef) // line 98
	fmt.Printf("fn.LvlInfo(0,fn.IflagsDef): \nfile:%s \nline:%v \nname:%s\n\n", file, line, name)

	file, line, name = fn.LvlInfo(0, fn.Ifilelong|fn.Ifnfull) // line 101
	fmt.Printf("fn.LvlInfo(0,fn.Ifilelong|fn.Ifnfull): \nfile:%s \nline:%v \nname:%s\n\n", file, line, name)

	file, line, name = fn.LvlInfo(0, fn.Ifilelong|fn.Ifilenogps) // line 104
	fmt.Printf("fn.LvlInfo(0,fn.Ifilelong|fn.Ifilenogps): \nfile:%s \nline:%v \nname:%s\n\n", file, line, name)

	fmt.Printf("fn.LvlInfo(0, fn.IflagsDef):%s\n", fn.LvlInfoStr(0, fn.IflagsDef)) // 107

	// Representative Output:
	// fn.LvlInfo(0,fn.IflagsDef):
	// file:github.com/phcurtis/fn/example_test.go
	// line:98
	// name:fn_test.ExampleLvlInfo()
	//
	// fn.LvlInfo(0,fn.Ifilelong|fn.Ifnfull):
	// file:/home/paul/go/src/github.com/phcurtis/fn/example_test.go
	// line:101
	// name:github.com/phcurtis/fn_test.ExampleLvlInfo()
	//
	// fn.LvlInfo(0,fn.Ifilelong|fn.Ifilenogps):
	// file:github.com/phcurtis/fn/example_test.go
	// line:104
	// name:github.com/phcurtis/fn_test.ExampleLvlInfo()
	//
	// fn.LvlInfo(0,fn.IflagsDef):
	// github.com/phcurtis/fn/example_test.go:107:fn_test.ExampleLvlInfo()
	//
}

func Example_lvlinfocmn() {
	fmt.Printf(fn.LvlInfoCmn(0)) // line 112
	// Representative Output:
	// github.com/phcurtis/fn/example_test.go:112:fn_test.ExampleLvlInfoCmn()
}

func Example_lvlinfoshort() { // line 117
	fmt.Printf(fn.LvlInfoShort(0))
	// Representative Output:
	// example_test.go:117:fn_test.ExampleLvlInfoShort()
}
