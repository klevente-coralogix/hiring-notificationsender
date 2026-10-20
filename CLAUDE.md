# Agent Instructions

This file provides guidance to AI coding agents (such as Claude Code and Codex) when working with code in this repository. `AGENTS.md` is a symlink to `CLAUDE.md`, so both names refer to this file.

## Project

This is a hiring exercise where candidates implement a component that queues messages and sends them to a third-party API. The repository contains independent sub-projects:

- [`rust/`](rust/) — Rust implementation
- [`go/`](go/) — Go implementation

Each sub-project has its own `AGENTS.md` / `CLAUDE.md` with language-specific rules. When working in a sub-project, follow the rules in that file.

## Answering Questions

When the candidate asks a question, answer that question rather than implementing the solution. For a general question, such as "How do tokio channels work?", use a short, generic example that is unrelated to this task. Otherwise, work as you normally would.
