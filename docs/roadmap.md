# Roadmap

`go-ruby-format/format` is grown **test-first**, each capability differential-tested against MRI
rather than built in isolation. Ruby's format-string engine — the
deterministic, interpreter-independent slice extracted from rbgo's internals — is
**complete**.

| Stage | What | Status |
| --- | --- | --- |
| Conversion directives | Every MRI conversion: integers `d`/`i`/`u`/`b`/`o`/`x`/`X`, floats `f`/`e`/`E`/`g`/`G`/`a`/`A`, `s` string, `p` inspect, `c` character, and the `%%` literal — each matching reference Ruby's rounding, padding and sign rules. | **Done** |
| Flags, width & precision | All flags (`-` left-justify, `+` and space sign, `0` zero-pad, `#` alternate form) with width and precision, including the `*` form that reads width/precision from an argument. | **Done** |
| Named & positional refs | Named references `%<name>s` and `%{name}` resolved against a hash, and absolute `n$` argument indexing — matching MRI's mixing rules and its errors for clashing styles. | **Done** |
| Bignum & hex-float | Arbitrary-precision integer output (`%d` on a Bignum), and IEEE hex-float (`%a` / `%A`) emitted exactly as MRI's `sprintf.c` does, down to the mantissa digits and exponent. | **Done** |
| Error taxonomy | MRI's full set of `ArgumentError` / `KeyError` / `TypeError` cases — too few arguments, unknown directive, mixed numbered/unnumbered, missing named key — surfaced as Go errors. | **Done** |
| Differential oracle & coverage | A wide format-string corpus rendered both here and by the system `ruby`, compared byte-for-byte against MRI; 100% coverage, gofmt + go vet clean, green across all six 64-bit Go arches and three OSes. | **Done** |

## Documented out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **No interpreter.** The library implements the deterministic algorithm; it
  never runs arbitrary Ruby. Anything that needs a live binding or evaluation is
  the consumer's job — that is why `rbgo` binds this module rather than the
  reverse.
- **Reference is reference Ruby (MRI).** Byte-for-byte conformance targets MRI's
  behaviour; differences across MRI releases are matched to the reference used by
  the differential oracle.
- **Standalone & reusable.** The module has no dependency on the Ruby runtime;
  the dependency runs the other way.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
deterministic/interpreter split.
