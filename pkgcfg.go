// Copyright 2017 phcurtis fn Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"io"
)

// PkgCfgStruct - package config structure less log Output.
type PkgCfgStruct struct {
	Logflags      int
	LogPrefix     string
	LogTraceFlags int
	LogAlignFile  int
	LogAlignFunc  int
}

// PkgCfgDef - returns package config defaults and logOutput
func PkgCfgDef() (pkgCfg *PkgCfgStruct, logOutput io.Writer) {
	muLogt.Lock()
	defer muLogt.Unlock()
	p := PkgCfgStruct{
		Logflags:      LflagsDef,
		LogPrefix:     LogPrefixDef,
		LogTraceFlags: TrFlagsDef,
		LogAlignFile:  LogAlignFileDef,
		LogAlignFunc:  LogAlignFuncDef,
	}
	return &p, logOutputDef
}

// SetPkgCfgDef - sets this package configuration to its defaults.
func SetPkgCfgDef(resetLogOutput bool) {
	muLogt.Lock()
	defer muLogt.Unlock()

	logt.SetFlags(LflagsDef)
	logt.SetPrefix(LogPrefixDef)
	logTraceFlags = TrFlagsDef
	logAlignFile = LogAlignFileDef
	logAlignFunc = LogAlignFuncDef
	if resetLogOutput {
		logSetOutput(logOutputDef)
	}
	return
}

// PkgCfg - returns current package config and logOutput
func PkgCfg() (pkgCfg *PkgCfgStruct, logOutput io.Writer) {
	muLogt.Lock()
	defer muLogt.Unlock()
	p := PkgCfgStruct{
		Logflags:      logt.Flags(),
		LogPrefix:     logt.Prefix(),
		LogTraceFlags: logTraceFlags,
		LogAlignFile:  logAlignFile,
		LogAlignFunc:  logAlignFunc,
	}
	return &p, logOutputCur
}

// SetPkgCfg - updates the passed in PkgCfgStruct to applicable vars
// and logOutput if that is not nil.
func SetPkgCfg(p *PkgCfgStruct, logOutput io.Writer) {
	muLogt.Lock()
	defer muLogt.Unlock()

	logt.SetFlags(p.Logflags)
	logt.SetPrefix(p.LogPrefix)
	logTraceFlags = p.LogTraceFlags
	logAlignFile = p.LogAlignFile
	logAlignFunc = p.LogAlignFunc
	if logOutput != nil {
		logSetOutput(logOutput)
	}
	return
}
