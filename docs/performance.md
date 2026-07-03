# Performance

`go-ruby-format/format` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `format`. This
page records a **comparative benchmark** of that module against the reference
Ruby runtimes, part of the ecosystem-wide per-module parity suite.

## What is measured

The **same** Ruby script — `sprintf` with mixed conversions (`%d %s %f %x %o %+d %e`) in a tight loop — is run under every runtime. `rbgo`'s
number reflects **this pure-Go library doing the work**; every other column is
that interpreter's own `format` stdlib. So the comparison is the **Ruby-visible
operation**, apples-to-apples across interpreters. The script prints a
deterministic checksum and its output is checked **byte-identical to MRI**
before timing.

- **Host:** Apple M4 Max, macOS (darwin/arm64). **Method:** best-of-5 wall time
  (best, not mean, to suppress scheduler noise); single-shot processes, no
  warm-up beyond the script's own loop.
- **Runtimes:** `ruby 4.0.5 +PRISM` (MRI, the oracle) and `ruby --yjit`;
  `jruby 10.1.0.0` (OpenJDK 25); `truffleruby 34.0.1` (GraalVM CE Native).
- The benchmark script and harness live in rbgo's repo under
  [`bench/modules/`](https://github.com/go-embedded-ruby/ruby/tree/main/bench/modules)
  (`format.rb` + `run.sh`). Reproduce:
  `RBGO=./rbgo TRUFFLE=truffleruby bash bench/modules/run.sh 5`.

## Result (best of 5, ms)

| Runtime | time | vs MRI |
| --- | ---: | ---: |
| **rbgo** (go-ruby-format) | 640 | 1.94× |
| MRI (ruby 4.0.5) | 330 | 1.00× |
| MRI + YJIT | 310 | 0.94× |
| JRuby 10.1.0.0 | 1490 | 4.52× |
| TruffleRuby 34.0.1 | 480 | 1.45× |

rbgo runs on **go-ruby-format**. The mixed-conversion sprintf loop is ~1.9x MRI / ~2.1x YJIT — within the clean-interpreter band; the per-iteration cost is dominated by interpreter dispatch around the format call, not the formatter itself.

!!! note "Honest framing"
    JRuby and TruffleRuby are timed **cold, single-shot**, so they carry JVM /
    Graal startup on every run — read them as one-shot `ruby file.rb` costs, the
    same way `rbgo` and MRI are measured, not as steady-state JIT numbers. Rows
    that complete in well under ~200 ms carry the most relative noise; treat
    their ratios as order-of-magnitude. These are real measured numbers from the
    2026-06-29 run — nothing is cherry-picked.

## Library-level benchmark (Go API vs runtimes) — 2026-07-03

This section measures the **pure-Go library directly, through its Go API**
(`format.Sprintf`) — not the `rbgo` interpreter path recorded above. It isolates
the formatting primitive from Ruby-interpreter dispatch, answering the parity
question head-on: *is the pure-Go implementation as fast as the reference
runtime's own `sprintf`?* The **same format strings, same arguments, same
iteration counts** run through the Go library and through each reference runtime's
native `Kernel#sprintf`.

Because formatting is **exact**, the whole output string is a strong oracle: the
Go library's output was checked **byte-identical to MRI `sprintf`** for every case
below before any timing (the combined string and all four single-directive
cases). Every MRI conversion is covered by the library — `d i u`, `f`, `e E`,
`g G`, `a A`, `s`, `p`, `x X`, `o`, `b B`, `c`, `%%`, with the `- + space 0 #`
flags, `*`/named width and precision, `%n$` argument references, `%<name>`/`%{name}`
hash references, and Bignum — so there is **no unhandleable op to exclude** from
this workload.

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5.1 — **date 2026-07-03**.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` · MRI + YJIT · JRuby 10.1.0.0
  (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native).
- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.

### Workload

One representative **combined** format string exercising every family in a single
call — `"%d %x %o %b | %f %e %g %+08.3f | [%s] [%-20s] 100%%"` — integer
(`d`/`x`/`o`/`b`), float (`f`/`e`/`g` and the flagged, width+precision `+08.3f`),
string (`s` and left-justified `-20s`), and a literal `%%`; plus four hot
single-directive cases (`%d`, `%+08.3f`, `%08x`, `%-20s`).

#### mixed

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 743.9 | 0.84× |
| MRI | 889.5 | 1.00× |
| MRI + YJIT | 824.5 | 0.93× |
| JRuby | 1757.1 | 1.98× |
| TruffleRuby | 1521.1 | 1.71× |

#### int-d

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 84.2 | 1.06× |
| MRI | 79.8 | 1.00× |
| MRI + YJIT | 48.0 | 0.60× |
| JRuby | 113.3 | 1.42× |
| TruffleRuby | 98.9 | 1.24× |

#### float-f

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 129.3 | 0.70× |
| MRI | 183.8 | 1.00× |
| MRI + YJIT | 131.6 | 0.72× |
| JRuby | 207.1 | 1.13× |
| TruffleRuby | 249.1 | 1.36× |

#### hex-x

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 94.4 | 0.76× |
| MRI | 124.0 | 1.00× |
| MRI + YJIT | 81.0 | 0.65× |
| JRuby | 153.3 | 1.24× |
| TruffleRuby | 133.7 | 1.08× |

#### str-s

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 95.3 | 0.98× |
| MRI | 97.0 | 1.00× |
| MRI + YJIT | 69.0 | 0.71× |
| JRuby | 144.6 | 1.49× |
| TruffleRuby | 51.0 | 0.53× |

**Reading the numbers.** On the **combined** format string — the realistic case,
ten conversions in one call — the pure-Go library is **faster than every runtime
measured**: 0.84× MRI and, notably, **0.93× MRI + YJIT** (0.744 µs vs MRI's
0.890 µs and YJIT's 0.825 µs). Per single directive it is **at or below MRI's C
`sprintf`** across the board — float `%+08.3f` 0.70×, hex `%08x` 0.76×, string
`%-20s` 0.98× — with only integer `%d` marginally above parity at 1.06×. YJIT
wins the isolated single-directive cases (0.60–0.72×), where the whole operation
is a handful of machine instructions and JIT specialisation dominates, but that
lead **inverts on the combined string**, where the per-call dispatch YJIT saves
is amortised over ten conversions and the go-ruby engine's single-pass parse
pulls ahead. JRuby and TruffleRuby trail here because they are timed on a fixed
warm-up budget (see below). These are the honest residual results: the pure-Go
formatter is **already at reference-C parity or better** on real mixed workloads,
so there is no per-op optimisation gap flagged for this module.

!!! note "Reproduce"
    The harness is committed under
    [`benchmarks/`](https://github.com/go-ruby-format/docs/tree/main/benchmarks):
    a self-contained Go driver (`go/`, pins the published library via `go.mod`),
    the equivalent `ruby/format.rb` workload, and `run.sh`. Run
    `bash benchmarks/run.sh`; env `OUTER`/`WARM` tune the pass budget and
    `RUBY`/`JRUBY`/`TRUFFLERUBY` select the runtime binaries. Each side also
    prints `ORACLE` lines (the exact formatted output) so the byte-identical
    check against MRI is reproducible.

!!! warning "Warm-up budget & noise — honest framing"
    Numbers reflect a **fixed warm-process budget** (3 warm-up + 25 timed passes
    in one process). The JVM/GraalVM JITs (JRuby, TruffleRuby) may need a larger
    warm-up to reach steady state, so their columns can **understate** peak
    throughput — visible in TruffleRuby's spread (0.53× on `str-s` yet 1.36× on
    `float-f`) across otherwise similar sub-microsecond loops. These sub-µs rows
    carry the most relative noise; treat their ratios as order-of-magnitude. Every
    number here is a **real measured value** from the dated run above (Apple M4
    Max, `ruby 4.0.5 +PRISM`, `jruby 10.1.0.0`, `truffleruby 34.0.1`) — nothing is
    fabricated, estimated, or cherry-picked. The go-ruby column is the pure-Go
    library; every other column is that interpreter's own `sprintf` doing the
    equivalent work, its output verified byte-identical to MRI before timing.
