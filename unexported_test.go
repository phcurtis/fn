// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
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
