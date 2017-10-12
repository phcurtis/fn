package fn_test

import (
	"fmt"
	"strings"

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

func ExampleLogBeg() {
	defer fn.LogEnd(fn.LogBeg())
	fmt.Println("testpoint-1")
	// Representative output of above code follows:
	// LogFN: 2017/10/10 09:54:02 Beg:github.com/phcurtis/fn_test.ExampleLogBeg
	// testpoint-1
	// LogFN: 2017/10/10 09:54:02 End:github.com/phcurtis/fn_test.ExampleLogBeg
}

func ExampleLogBegDur() {
	defer fn.LogEndDur(fn.LogBegDur())
	fmt.Println("testpointDur")
	// Representative output of above code follows:
	// LogFN: 2017/10/10 12:51:16 BegDur:github.com/phcurtis/fn_test.ExampleLogBegDur
	// testpointDur
	// LogFN: 2017/10/10 12:51:16 EndDur:github.com/phcurtis/fn_test.ExampleLogBegDur Dur:4.323µs
}
