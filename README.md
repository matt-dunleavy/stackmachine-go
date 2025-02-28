# Stack Machine Go

This is a Go port of the Stack Machine project originally created by Christian Stigen Larsen in 2010. The original project is available at [csl.sublevel3.org](http://csl.sublevel3.org).

## Overview

Stack Machine Go is a simple stack-based virtual machine that can run Forth/PostScript-like programs. It includes:

- A virtual machine for executing low-level stack-based instructions
- An assembler supporting a Forth-like language
- An interpreter for running programs on-the-fly
- A disassembler for inspecting compiled bytecode

## Installation

With Go installed, you can install the stack machine with:

```bash
go get github.com/yourusername/stack-machine-go
```

Or clone the repository and build from source:

```bash
git clone https://github.com/yourusername/stack-machine-go.git
cd stack-machine-go
go build
```

## Usage

The stack machine has four main commands:

- `interpret` (alias: `sm`): Compile and run source on-the-fly
- `compile` (alias: `smc`): Compile source to bytecode
- `run` (alias: `smr`): Run compiled bytecode
- `disassemble` (alias: `smd`): Disassemble bytecode

### Examples

Interpret a file:
```bash
stackmachine interpret examples/hello.src
# or
stackmachine sm examples/hello.src
```

Compile a file:
```bash
stackmachine compile examples/hello.src
# or
stackmachine smc examples/hello.src
```

Run a compiled file:
```bash
stackmachine run hello.bin
# or
stackmachine smr hello.bin
```

Disassemble a compiled file:
```bash
stackmachine disassemble hello.bin
# or
stackmachine smd hello.bin
```

## The Stack Machine Language

The stack machine language is similar to Forth and PostScript, using reverse Polish notation (RPN). Here's a simple "Hello, World!" example:

```
; Hello, World!
main:
   72 out          ; "H"
  101 out          ; "e"
  108 dup out out  ; "ll"
  111 out          ; "o"
   44 out          ; ","
   32 out          ; " "
   87 out          ; "W"
  111 out          ; "o"
  114 out          ; "r"
  108 out          ; "l"
  100 out          ; "d"
   33 out          ; "!"
  '\n' out         ; newline
  halt
```

### Key Concepts

1. **Labels**: Identifiers ending with a colon
   ```
   label:      ; define a label
   &label      ; put ADDRESS of label on stack
   &label LOAD ; put VALUE of label on stack
   label       ; EXECUTE code at label position
   ```

2. **Operations**: Stack-based operations like ADD, SUB, AND, etc.

3. **Character literals**: 'a', '\n', etc.

4. **Function calls**: Labels can be used as function calls with automatic return via POPIP

For more detail, see the [original documentation](http://csl.sublevel3.org/stack-machine/).

## License

This project is placed in the public domain, following the original work by Christian Stigen Larsen.