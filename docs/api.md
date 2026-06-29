# Usage & API

The public API lives at the module root (`github.com/go-ruby-format/format`). It is **Ruby-shaped but Go-idiomatic**: `Sprintf` mirrors `Kernel#sprintf`, while the surface follows Go conventions — an explicit `error`, value types, no global state.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-format/format`, bound into
    `rbgo` as a native module; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-format/format
```

## Worked example

```go
s, _ := format.Sprintf("%05.2f%%", 3.14159)   // "03.14%"
s, _  = format.Sprintf("%1$s %1$s", "hi")      // "hi hi"
s, _  = format.Sprintf("%#x", 255)              // "0xff"
s, _  = format.Sprintf("%b", bigInt)            // arbitrary precision
s, _  = format.Format("%<who>s!", nil,
        format.Named(map[string]format.Value{"who": "world"}))  // "world!"
```

## Shape

```go
// Sprintf renders an MRI format string against positional args,
// matching Kernel#sprintf / format / String#% byte-for-byte.
func Sprintf(format string, args ...Value) (string, error)

// Format is the lower-level entry: positional args plus a named-arg
// table for %<name>s / %{name} references.
func Format(format string, args []Value, named *NamedArgs) (string, error)
```

## MRI conformance

Correctness is defined by reference Ruby. A **differential oracle** runs a wide
corpus through both the system `ruby` and this library and compares the results
**byte-for-byte** — not approximated from memory. The oracle tests skip
themselves where `ruby` is not on `PATH` (e.g. the qemu arch lanes), so the
cross-arch builds still validate the library.

## Relationship to Ruby

`go-ruby-format/format` is **standalone and reusable**, and is the backend bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo` as a
native module — the same way [go-ruby-regexp](https://github.com/go-ruby-regexp)
and [go-ruby-erb](https://github.com/go-ruby-erb) are bound. The dependency runs
the other way: this library has no dependency on the Ruby runtime.
