# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a stack-based calculator REPL written in Go. It combines RPN (Reverse Polish Notation) operators with expression evaluation using the `expr-lang/expr` library. The calculator maintains a stack that displays in the prompt, allowing both traditional expression input (e.g., "2 + 3") and RPN operations (e.g., push values then use "+").

## Build System

The project uses GNU Make with parallel builds (`MAKEFLAGS += -j`). All build artifacts go to `build/` directory.

### Common Commands

- `make` or `make all` - Build binary, run tests, linting, security checks, and vulnerability scans
- `make deep` - Run all checks plus additional static analysis (nilaway, errcheck)
- `make watch` - Auto-rebuild on file changes using `entr`
- `make clean` - Remove all build artifacts
- `./build/c` - Run the calculator REPL

### Build Targets

The Makefile runs these checks automatically:
- `go test` with `-failfast -parallel=8 -count=2 -shuffle=on`
- `golangci-lint run` - Standard linting
- `gosec` - Security scanning
- `govulncheck` - Vulnerability detection
- `gofumpt` - Code formatting (applied automatically during build)
- `scc` - Source code statistics (displayed at end of build)
- Optional deep analysis: `nilaway` (nil safety), `errcheck` (error handling)

All checks output to `build/*.out` files and only display on failure.

## Architecture

### Stack-Based Calculator Core

The calculator (main.go) operates as a REPL with two input modes:

1. **Operator Mode**: Recognizes operators from the `operators` map, executes them on the stack
2. **Expression Mode**: Falls back to `expr.Eval()` for any non-operator input, evaluates expressions and pushes numeric results to stack

### Key Components

- `pop2()`: Helper that pops two operands from stack for binary operations
- `Operator` type: Function signature for stack operations
- `operators` map: Currently only implements "+" operator (note: has bug where it modifies local copy of stack)
- Stack visualization: Prompt updates to show `[ n1 n2 n3 ]>` format after each operation
- History: Persists to `.history` file via readline

### Current Issues

The "+" operator in main.go:25-35 has a bug: it passes a copy of the stack to the operator function, so `pop2(&stack)` modifies the local copy, not the main stack. The operation computes correctly but doesn't actually consume operands from the stack.

## Development Notes

- The project uses Go 1.25.2 and requires `github.com/chzyer/readline` and `github.com/expr-lang/expr`
- The build system automatically runs `go mod tidy` and `go generate` before compilation
- Code formatting with `gofumpt` is applied during every build
- The dependency tracking system (`code-deps.mk`) ensures rebuilds when Go files are added/removed
