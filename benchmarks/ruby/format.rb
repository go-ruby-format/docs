# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
#
# Reference-runtime side of the go-ruby-format library benchmark. Drives the
# native Kernel#sprintf of MRI / MRI+YJIT / JRuby / TruffleRuby on exactly the
# format strings and arguments the Go driver (go/main.go) drives through the
# pure-Go format.Sprintf. Emits RESULT lines (consumed by run.sh) and ORACLE
# lines (the exact output, diffed byte-identical to MRI by the verify step).
require_relative "_harness"

# One representative combined format string exercising every family: integer
# (d/x/o/b), float (f/e/g and the flagged, width+precision +08.3f), string
# (s and left-justified -20s), and a literal %%.
MIXED_FMT = "%d %x %o %b | %f %e %g %+08.3f | [%s] [%-20s] 100%%"

# Fixed, deterministic arguments — identical to go/main.go.
MIXED_ARGS = [
  1234567,          # %d
  3735928559,       # %x  -> deadbeef
  342391,           # %o
  5461,             # %b  -> 1010101010101
  3.14159265358979, # %f
  6.022e23,         # %e
  0.000123456,      # %g
  -2.71828,         # %+08.3f
  "hello",          # %s
  "world"           # %-20s
].freeze

def oracle(label, s)
  printf("ORACLE\t%s\t%s\n", label, s)
end

# Strong oracle: printed before any timing.
oracle("mixed",   sprintf(MIXED_FMT, *MIXED_ARGS))
oracle("int-d",   sprintf("%d", 1234567))
oracle("float-f", sprintf("%+08.3f", -2.71828))
oracle("hex-x",   sprintf("%08x", 3735928559))
oracle("str-s",   sprintf("%-20s", "hello"))

bench("mixed",   2000) { sprintf(MIXED_FMT, *MIXED_ARGS) }
bench("int-d",   5000) { sprintf("%d", 1234567) }
bench("float-f", 5000) { sprintf("%+08.3f", -2.71828) }
bench("hex-x",   5000) { sprintf("%08x", 3735928559) }
bench("str-s",   5000) { sprintf("%-20s", "hello") }
