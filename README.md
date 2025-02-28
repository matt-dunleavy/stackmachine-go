# Stack Machine Go

A stack-based virtual machine and compiler written in Go.

[![License: ISC](https://img.shields.io/badge/License-ISC-blue.svg)](https://opensource.org/licenses/ISC)

## Overview

Stack Machine Go is a lightweight, stack-based virtual machine that executes simple programs in a Forth/PostScript-like language. It consists of:

- A virtual machine (VM) that executes bytecode
- A compiler that translates assembly-like source code to bytecode
- An interpreter that compiles and runs code on-the-fly
- A disassembler that converts bytecode back to human-readable form

This project is a Go port of the original Stack Machine created by Christian Stigen Larsen in 2010-2011.

## Features

- Simple but powerful instruction set
- Stack-based execution model
- Memory manipulation
- Branching and function calls
- Input/output operations
- Label-based addressing

## Installation

### Prerequisites

- Go 1.18 or higher

### Building from Source

```bash
git clone https://github.com/matt-dunleavy/stackmachine-go.git
cd go-vm
go build -o stackmachine
```

Alternatively, use the provided Makefile:

```bash
make build
```

## Usage

Stack Machine Go has several commands:

### Compile Source to Bytecode

```bash
smg compile hello.src
# Outputs hello.bin
```

### Run Compiled Bytecode

```bash
smg run hello.bin
```

### Interpret Source Code Directly

```bash
smg interpret hello.src
```

### Disassemble Bytecode

```bash
smg disassemble hello.bin
```

### Using Standard Input/Output

All commands accept input from stdin if no file is specified:

```bash
cat hello.src | smg interpret
```

## Source Code Examples

### Hello World

```
; Hello World program
'H OUT
'e OUT
'l OUT
'l OUT
'o OUT
', OUT
' OUT
'W OUT
'o OUT
'r OUT
'l OUT
'd OUT
'! OUT
'\n OUT
HALT
```

### Simple Function Call

```
; Main program
'H OUT
'i OUT
'! OUT
' OUT
PUSHIP 24  ; Push return address
32         ; Call function at address 32
JMP        ; Jump to function
'\n OUT    ; Print newline after return
HALT       ; End program

; Function that prints "there"
't OUT     ; Print "there"
'h OUT
'e OUT
'r OUT
'e OUT
POPIP      ; Return to caller
```

## Architecture

The Stack Machine is based on a simple stack-based architecture:

- Operations operate on a data stack
- A separate instruction pointer (IP) stack enables function calls
- Memory is organized as a flat array of 32-bit words
- Instructions are encoded as 32-bit words

See the [docs](./docs/) directory for detailed documentation on the architecture, instruction set, and compiler.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the ISC License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Original Stack Machine by Christian Stigen Larsen
