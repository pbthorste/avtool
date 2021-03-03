# avtool

[![License: MIT](https://img.shields.io/badge/License-GPL_v3-brightgreen.svg)](https://github.com/clok/avtool/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/clok/avtool)](https://goreportcard.com/report/clok/avtool)
[![Coverage Status](https://coveralls.io/repos/github/clok/avtool/badge.svg?branch=main)](https://coveralls.io/github/clok/avtool?branch=main)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/clok/avtool/v2?tab=overview)

> NOTE: Original code written by [@pbthorste](https://github.com/pbthorste) for [https://github.com/pbthorste/avtool](https://github.com/pbthorste/avtool)
>
> HUGE SHOUT OUT to [@pbthorste](https://github.com/pbthorste)

This module provides a reimplementation of `ansible-vault` encrypt and decprypt functionality in Go.

## Why the fork?

As of writing the mainline has been stale for ~4 years.

I have found this code to be highly useful and important for writing other `ansible-vault` related tools. I wanted to modernize the work done previously to support `go.mod` while also updating the interface as an importable module for other code.

## CLI Tool

Please see [gwvault](https://github.com/GoodwayGroup/gwvault) for a purpose built `ansible-vault` binary written in go.

It leverages the work done by [@pbthorste](https://github.com/pbthorste) for [https://github.com/pbthorste/avtool](https://github.com/pbthorste/avtool) while further fleshing out the CLI tool to be more in line with the original `ansible-vault` CLI tool.

## Thanks and Attribution

Original code written by [@pbthorste](https://github.com/pbthorste)
