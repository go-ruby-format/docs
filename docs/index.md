# go-ruby-format documentation

**Ruby's `Kernel#sprintf` / `format` / `String#%` in pure Go — MRI-compatible, no cgo.**

`go-ruby-format/format` is a faithful, pure-Go (zero cgo) reimplementation of Ruby's format-string engine,
matching reference Ruby (MRI) byte-for-byte. The module path is
`github.com/go-ruby-format/format`.

It was **extracted from rbgo's prelude/internals into a reusable standalone
library**: the module is standalone and importable by any Go program, and it is
the backend bound into [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby)
by `rbgo` as a native module — just like
[go-ruby-regexp](https://github.com/go-ruby-regexp) and
[go-ruby-erb](https://github.com/go-ruby-erb). The dependency runs the other
way: this library has **no dependency on the Ruby runtime**.

!!! success "Status: engine complete — MRI byte-exact"
    Faithful port of MRI's `sprintf.c`: **all conversions** (`d i f e g s p x o b c %`), **all flags** (`- + space 0 #`), **width and precision** including the `*` argument form, **named references** (`%<name>s` / `%{name}`), **`n$` indexing**, **Bignum** and **hex-float** (`%a`). Validated by a **differential oracle** against the system `ruby` — rendered output compared byte-for-byte — at 100% coverage, `gofmt` + `go vet` clean, CI green across the six 64-bit Go targets and three OSes.

## Quick taste

```go
s, _ := format.Sprintf("%05.2f%%", 3.14159)   // "03.14%"
s, _  = format.Sprintf("%1$s %1$s", "hi")      // "hi hi"
s, _  = format.Sprintf("%#x", 255)              // "0xff"
s, _  = format.Sprintf("%b", bigInt)            // arbitrary precision
s, _  = format.Format("%<who>s!", nil,
        format.Named(map[string]format.Value{"who": "world"}))  // "world!"
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`format`](https://github.com/go-ruby-format/format) | the library — Ruby's format-string engine in pure Go |
| [`docs`](https://github.com/go-ruby-format/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-format.github.io`](https://github.com/go-ruby-format/go-ruby-format.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-format/brand) | logo and brand assets |

## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- **MRI byte-exact.** Output matches reference Ruby exactly, not approximately,
  validated by a differential oracle against the `ruby` binary.
- **Standalone & reusable.** Extracted from rbgo's internals; no dependency on
  the Ruby runtime — the dependency runs the other way.
- **100% test coverage** is the target, enforced as a CI gate, across 6 arches
  and 3 OSes.

## Where to go next

- [Why pure Go](why.md) — why this slice of Ruby is deterministic enough to live
  as a standalone, interpreter-independent Go library.
- [Usage & API](api.md) — the public surface and worked examples.
- [Roadmap](roadmap.md) — what is done and what is downstream by design.

Source lives at [github.com/go-ruby-format/format](https://github.com/go-ruby-format/format).
