# Virtual Machine Architecture

## Overview

The Stack Machine is a simple stack-based virtual machine that executes bytecode instructions. Each instruction operates on a stack of 32-bit values. The machine has a flat memory model, a data stack, and an instruction pointer stack for handling function calls.

## Components

### Memory

The VM has a linear memory space organized as an array of 32-bit words. By default, the memory size is 1024 KB (1,048,576 bytes or 262,144 words).

Memory is used to store:
- Program code (instructions)
- Data
- The runtime stack

### Data Stack

The main stack is used for data manipulation. Most operations pop values from the stack, perform calculations, and push results back.

Example operations:
- `ADD`: Pops two values, adds them, and pushes the result
- `PUSH`: Pushes a literal value onto the stack
- `DUP`: Duplicates the top stack value

### Instruction Pointer

The instruction pointer (IP) keeps track of the currently executing instruction. It automatically advances after each instruction unless modified by a jump operation.

### Instruction Pointer Stack

A separate stack is used to store return addresses for function calls. This enables subroutine functionality:
- `PUSHIP`: Pushes the current IP onto the IP stack
- `POPIP`: Pops an address from the IP stack and jumps to it
- `DROPIP`: Pops and discards an address from the IP stack

### I/O System

The VM includes simple I/O operations:
- `IN`: Reads one byte from stdin and pushes it onto the stack
- `OUT`: Pops a value from the stack and writes it to stdout as a byte
- `OUTNUM`: Pops a value and prints it as a number

## Execution Model

1. The VM loads bytecode into memory starting at address 0
2. Execution begins with the instruction pointer at 0
3. Each instruction is fetched, decoded, and executed
4. The instruction pointer advances to the next instruction (typically +4 bytes)
5. Jumps and function calls can change the instruction pointer
6. Execution continues until a halt instruction is encountered

A program halts when it executes a jump to the current address:
```
PUSH <current_address>
JMP
```

## Word Size

The VM uses 32-bit (4-byte) words for:
- Memory addressing
- Stack values
- Instructions

## Error Handling

The VM implements simple error handling for common issues:
- Stack underflow
- Memory access outside bounds
- Unknown instructions

When an error occurs, the VM calls an error callback function that can be provided during initialization.

## Threading Model

The VM is single-threaded. It processes one instruction at a time and does not provide native concurrency features.

## Memory Management

Memory is statically allocated at VM creation time. The VM does not implement automatic memory management (garbage collection), leaving memory management to the programmer.

## VM Lifecycle

1. **Creation**: A new VM is created with `NewMachine()` or `NewMachineWithSize()`
2. **Loading**: Program code is loaded with `LoadImage()`
3. **Execution**: The program is executed with `Run()`
4. **Reset**: The VM can be reset with `Reset()`
5. **Termination**: When a halt instruction is executed, the VM stops running

## Clone and State Management

The VM state can be cloned using the `Clone()` method, which creates an independent copy of the VM with the same memory, stack contents, and instruction pointer. 