# Stack Machine Go Documentation

Welcome to the Stack Machine Go documentation. This documentation provides comprehensive information about the Stack Machine, a stack-based virtual machine implemented in Go.

## Contents

- [Virtual Machine Architecture](virtual-machine.md) - Details about the VM's architecture, components, and execution model
- [Instruction Set](instruction-set.md) - Complete reference of all supported instructions
- [Compiler](compiler.md) - Information about the compiler architecture and language syntax
- [Command-Line Interface](cli.md) - Guide to using the command-line tools
- [Example Programs](examples.md) - Walkthrough of example programs demonstrating various features

## Quick Links

- [Installation and Setup](../README.md#installation)
- [Basic Usage](../README.md#usage)
- [Source Code Examples](../README.md#source-code-examples)

## Project Overview

Stack Machine Go is a lightweight, stack-based virtual machine that executes simple programs in a Forth/PostScript-like language. It consists of:

- A virtual machine (VM) that executes bytecode
- A compiler that translates assembly-like source code to bytecode
- An interpreter that compiles and runs code on-the-fly
- A disassembler that converts bytecode back to human-readable form

## Component Relationships

1. **Source Code** → **Compiler** → **Bytecode** → **VM** → **Execution**
2. **Bytecode** → **Disassembler** → **Assembly Source**
3. **Source Code** → **Interpreter** → **Execution**

## Development

This project is an educational tool that demonstrates:

- Stack-based architecture principles
- Compiler construction
- Virtual machine implementation
- Assembly language design

Feel free to explore the source code and documentation to understand how each component works and interacts with others. 