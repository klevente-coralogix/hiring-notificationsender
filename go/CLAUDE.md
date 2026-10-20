# Agent Instructions

This file provides guidance to AI coding agents (such as Claude Code and Codex) when working with code in this repository. `AGENTS.md` is a symlink to `CLAUDE.md`, so both names refer to this file.

## Project

This is a hiring exercise where candidates implement, in Go, a component that queues messages and sends them to a third-party API.

## Rules

- **No external dependencies.** Use only the Go standard library. Do not add modules to `go.mod` unless the user explicitly requests a specific one.
- **Do not modify `thirdparty/thirdparty.go`** unless the user explicitly asks. It represents an external API and its mock.
- **Do not write tests unless explicitly asked.** Never add `_test.go` files or test functions without a direct request.
