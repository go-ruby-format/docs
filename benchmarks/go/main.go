// SPDX-License-Identifier: BSD-3-Clause
//
// Library-level driver for github.com/go-ruby-format/format. Times the pure-Go
// Kernel#sprintf/format engine through its Go API (format.Sprintf) on the same
// format strings and arguments the Ruby side (ruby/format.rb) drives through the
// reference runtimes' native sprintf. Emits RESULT lines (consumed by run.sh)
// and ORACLE lines (the exact formatted output, diffed byte-identical to MRI by
// the verify step; run.sh ignores non-RESULT lines).
package main

import (
	"fmt"

	"github.com/go-ruby-format/format"
)

func oracle(label, s string) { fmt.Printf("ORACLE\t%s\t%s\n", label, s) }

func mustFmt(f string, args ...any) string {
	s, err := format.Sprintf(f, args...)
	if err != nil {
		panic(err)
	}
	return s
}

// One representative combined format string exercising every family: integer
// (d/x/o/b), float (f/e/g and the flagged, width+precision +08.3f), string
// (s and left-justified -20s), and a literal %%.
const mixedFmt = "%d %x %o %b | %f %e %g %+08.3f | [%s] [%-20s] 100%%"

// Fixed, deterministic arguments — identical to ruby/format.rb.
func mixedArgs() []any {
	return []any{
		1234567,          // %d
		3735928559,       // %x  -> deadbeef
		342391,           // %o
		5461,             // %b  -> 1010101010101
		3.14159265358979, // %f
		6.022e23,         // %e
		0.000123456,      // %g
		-2.71828,         // %+08.3f
		"hello",          // %s
		"world",          // %-20s
	}
}

func main() {
	ma := mixedArgs()

	// Strong oracle: format is exact, so the whole output string must match MRI
	// sprintf byte-for-byte. Printed before any timing.
	oracle("mixed", mustFmt(mixedFmt, ma...))
	oracle("int-d", mustFmt("%d", 1234567))
	oracle("float-f", mustFmt("%+08.3f", -2.71828))
	oracle("hex-x", mustFmt("%08x", 3735928559))
	oracle("str-s", mustFmt("%-20s", "hello"))

	bench("mixed", 2000, func() { sink, _ = format.Sprintf(mixedFmt, ma...) })
	bench("int-d", 5000, func() { sink, _ = format.Sprintf("%d", 1234567) })
	bench("float-f", 5000, func() { sink, _ = format.Sprintf("%+08.3f", -2.71828) })
	bench("hex-x", 5000, func() { sink, _ = format.Sprintf("%08x", 3735928559) })
	bench("str-s", 5000, func() { sink, _ = format.Sprintf("%-20s", "hello") })
}
