<!-- SPDX-License-Identifier: BSD-3-Clause -->
# `go-ruby-format` library-level benchmark harness

Reproducible, cross-runtime benchmark of the **pure-Go `go-ruby-format` library**
(`github.com/go-ruby-format/format`, the engine behind Ruby's `Kernel#sprintf`,
`Kernel#format`, and `String#%`) against the reference Ruby runtimes (MRI,
MRI + YJIT, JRuby, TruffleRuby). It measures the **library primitive** through
its Go API (`format.Sprintf`), isolated from the rbgo interpreter, so the numbers
answer: *is the pure-Go implementation as fast as the reference runtime's own
`sprintf`?*

## Layout

- `go/`             — self-contained Go driver; `go.mod` pins the published
  library by pseudo-version (not a `replace`). The built `bench` binary is
  git-ignored.
- `ruby/format.rb`  — the equivalent workload; `ruby/_harness.rb` is the shared
  timer.
- `run.sh`          — runs every available runtime and prints one Markdown table
  per sub-benchmark (ns/op + ratio vs MRI).

## Run

```sh
bash benchmarks/run.sh
```

Environment knobs: `OUTER` (timed passes, default 25), `WARM` (untimed warm-up
passes, default 3), and `RUBY`/`JRUBY`/`TRUFFLERUBY` to select runtime binaries.

## Method

Each process runs `WARM` untimed passes (to let the JVM/GraalVM JITs warm up),
then `OUTER` timed passes of a fixed inner loop, timed with a monotonic clock;
the **best** pass is reported as **ns/op**. Interpreter start-up is outside the
timed region. The Go driver and the Ruby script format the **same format strings
with the same arguments**. Because formatting is exact, the whole output string
is a strong oracle: each side prints `ORACLE` lines (the exact formatted output)
that are checked **byte-identical to MRI `sprintf`** before any timing (`run.sh`
consumes only `RESULT` lines, so `ORACLE` lines are inert there). Results are
published, dated, in `../docs/performance.md`.

## Workload

- **`mixed`** — one representative combined format string exercising every
  family in a single call:
  `"%d %x %o %b | %f %e %g %+08.3f | [%s] [%-20s] 100%%"` — integer (`d`/`x`/`o`/`b`),
  float (`f`/`e`/`g` and the flagged, width+precision `+08.3f`), string (`s` and
  left-justified `-20s`), and a literal `%%`.
- **`int-d`**, **`float-f`**, **`hex-x`**, **`str-s`** — hot single-directive
  cases (`%d`, `%+08.3f`, `%08x`, `%-20s`) for the per-conversion cost.
