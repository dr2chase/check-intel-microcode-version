// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type pinfo struct {
	processor, stepping, family, model, steppingID, platformID string
	oldVersion, newVersion                                     int
	products                                                   string
}

// {"Processor Model","Stepping","Family Code","Model Number","Stepping Id","Platform Id","Old Version","New Version","Products"},

var processors = []pinfo{
	{"SKL-U/Y", "D0", "6", "4e", "3", "c0", 0xd4, 0xd6, "Core Gen6 Mobile"},
	{"SKL-U23e", "K1", "6", "4e", "3", "c0", 0xd4, 0xd6, "Core Gen6 Mobile"},
	{"SKL-H/S/E3", "N0/R0/S0", "6", "5e", "3", "36", 0xd4, 0xd6, "Core Gen6"},
	{"AML-Y22", "H0", "6", "8e", "9", "10", 0xc6, 0xca, "Core Gen8 Mobile"},
	{"KBL-U/Y", "H0", "6", "8e", "9", "c0", 0xc6, 0xca, "Core Gen7 Mobile"},
	{"KBL-U23e", "J1", "6", "8e", "9", "c0", 0xc6, 0xca, "Core Gen7 Mobile"},
	{"CFL-U43e", "D0", "6", "8e", "a", "c0", 0xc6, 0xca, "Core Gen8 Mobile"},
	{"KBL-R U", "Y0", "6", "8e", "a", "c0", 0xc6, 0xca, "Core Gen8 Mobile"},
	{"WHL-U", "W0", "6", "8e", "b", "d0", 0xc6, 0xca, "Core Gen8 Mobile"},
	{"AML-Y42", "V0", "6", "8e", "c", "94", 0xc6, 0xca, "Core Gen10 Mobile"},
	{"WHL-U", "V0", "6", "8e", "c", "94", 0xc6, 0xca, "Core Gen8 Mobile"},
	{"CML-U42", "V0", "6", "8e", "c", "94", 0xc6, 0xca, "Core Gen10 Mobile"},
	{"KBL-G/H/S/X/E3", "B0", "6", "9e", "9", "2a", 0xc6, 0xca, "Core Gen7 Desktop, Mobile, Xeon E3 v6"},
	{"CFL-H/S/E3", "U0", "6", "9e", "a", "22", 0xc6, 0xca, "Core Gen8 Desktop, Mobile, Xeon E"},
	{"CFL-S", "B0", "6", "9e", "b", "02", 0xc6, 0xca, "Core Gen8"},
	{"CFL-S", "P0", "6", "9e", "c", "22", 0xc6, 0xca, "Core Gen9 Desktop"},
	{"CFL-H/S/E3", "R0", "6", "9e", "d", "22", 0xc6, 0xca, "Core Gen9 Desktop, Mobile, Xeon E"},
	{"CML-U62", "A0", "6", "a6", "0", "80", 0xc6, 0xca, "Core Gen10 Mobile"}}

var getHWInfo func(what string) string

func darwin(what string) string {
	cmd := exec.Command("sysctl", what)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return string(out)
}

var cachedLscpu []string

func linux(what string) string {
	if len(cachedLscpu) == 0 {
		cmd := exec.Command("egrep", `(stepping|model|microcode|cpu family)\W*:`, "/proc/cpuinfo")
		out, err := cmd.CombinedOutput()
		if err != nil {
			panic(err)
		}
		for _, b := range bytes.Split(out, []byte("\n")) {
			cachedLscpu = append(cachedLscpu, string(b))
		}
	}
	for _, s := range cachedLscpu {
		if strings.HasPrefix(s, what) {
			return s
		}
	}
	return ""
}

func hwString(what string) string {
	s := getHWInfo(what)
	if s == "" {
		return s
	}
	s = strings.ReplaceAll(s, what, "")
	s = strings.TrimSpace(s)
	s = s[1:] // Lose the colon
	return strings.TrimSpace(s)
}

func hwInt(what string) int {
	s := hwString(what)
	if s == "" {
		return 0
	}
	if strings.HasPrefix(s, "0x") {
		i, err := strconv.ParseInt(s[2:], 16, 32)
		if err != nil {
			panic(err)
		}
		return int(i)
	}
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(i)
}

func hwHexString(what string) string {
	i := hwInt(what)
	return strconv.FormatInt(int64(i), 16)
}

func thisPinfo() (p *pinfo, ucode int, match func(have, reference *pinfo) bool) {
	p = &pinfo{}
	if runtime.GOARCH != "amd64" {
		panic(runtime.GOARCH + " is not supported, want amd64")
	}
	if runtime.GOOS == "darwin" {
		getHWInfo = darwin
		p.steppingID = hwHexString("machdep.cpu.stepping")
		p.model = hwHexString("machdep.cpu.model")
		p.family = hwHexString("machdep.cpu.family")
		ucode = hwInt("machdep.cpu.microcode_version")
	} else if runtime.GOOS == "linux" {
		getHWInfo = linux
		p.steppingID = hwHexString("stepping")
		p.model = hwHexString("model")
		p.family = hwHexString("cpu family")
		ucode = hwInt("microcode")
	} else {
		panic(runtime.GOOS + " is not a saupported, maybe you could fix that")
	}
	p.oldVersion = ucode
	p.newVersion = ucode
	match = func(have, reference *pinfo) bool {
		return have.steppingID == reference.steppingID && have.model == reference.model && have.family == reference.family
	}
	return
}

func main() {
	p, ucode, match := thisPinfo()
	fmt.Printf("This processor is %+v\n", *p)
	matched := false
	for _, q := range processors {
		if match(p, &q) {
			fmt.Printf("This computer matches %+v, this microcode is %d\n", q, ucode)
			matched = true
		}
	}
	if !matched {
		fmt.Printf("No processors on the microcode update list matched\n")
	}
}
