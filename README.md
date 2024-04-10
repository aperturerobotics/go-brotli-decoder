# Brotli Decoder for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/aperturerobotics/go-brotli-decoder.svg)](https://pkg.go.dev/github.com/aperturerobotics/go-brotli-decoder)
[![Go Report Card Widget]][Go Report Card]

[Go Report Card Widget]: https://goreportcard.com/badge/github.com/aperturerobotics/go-brotli-decoder
[Go Report Card]: https://goreportcard.com/report/github.com/aperturerobotics/go-brotli-decoder

## Introduction

This package is a brotli decompressor implemented in pure Go.

This package is a brotli compressor and decompressor implemented in Go.

It was translated from the reference implementation (https://github.com/google/brotli)
with the `c2go** tool at https://github.com/andybalholm/c2go.

## Upstream

This package is a fork of the [upstream project] to create a more minimal
package with just the decoder and not the encoder.

This is a significantly lighter package in terms of binary size.

It was created by deleting the writer types and then repeatedly removing all
unused symbols (detected with the gounused linter).

If you need the brotli compressor, see the [upstream project].

[upstream project]: https://github.com/andybalholm/brotli
