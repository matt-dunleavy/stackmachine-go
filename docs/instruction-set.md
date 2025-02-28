# Instruction Set

The Stack Machine uses a simple instruction set where most operations manipulate values on a data stack. This document outlines all available instructions, their behavior, and usage.

## Notation

In the descriptions below:
- "pop" means to remove the top value from the stack
- "push" means to add a value to the top of the stack
- The stack grows upward: the top element is the most recently pushed

## Core Instructions

| Opcode | Mnemonic | Description |
|--------|----------|-------------|
| 0x0    | NOP      | No operation |
| 0x1    | ADD      | Pop b, pop a, push (a + b) |
| 0x2    | SUB      | Pop a, pop b, push (b - a) |
| 0x3    | AND      | Pop b, pop a, push (a & b) |
| 0x4    | OR       | Pop b, pop a, push (a \| b) |
| 0x5    | XOR      | Pop b, pop a, push (a ^ b) |
| 0x6    | NOT      | Pop a, push (a == 0 ? 1 : 0) |
| 0x17   | COMPL    | Pop a, push (~a) (bitwise complement) |

## Stack Manipulation

| Opcode | Mnemonic | Description |
|--------|----------|-------------|
| 0xD    | PUSH     | Push the next word as a literal value |
| 0xE    | DUP      | Duplicate the top value on the stack |
| 0xF    | SWAP     | Swap the top two values on the stack |
| 0x10   | ROL3     | Rotate top three values: (a b c) → (b c a) |
| 0x13   | DROP     | Remove the top value from the stack |

## Memory Access

| Opcode | Mnemonic | Description |
|--------|----------|-------------|
| 0x9    | LOAD     | Pop address a, push value at memory[a] |
| 0xA    | STOR     | Pop address a, pop value b, store b at memory[a] |

## Control Flow

| Opcode | Mnemonic | Description |
|--------|----------|-------------|
| 0xB    | JMP      | Pop address a, jump to a |
| 0xC    | JZ       | Pop value a, pop address b, jump to b if a == 0 |
| 0x12   | JNZ      | Pop value a, pop address b, jump to b if a != 0 |

## Function Calls

| Opcode | Mnemonic | Description |
|--------|----------|-------------|
| 0x14   | PUSHIP   | Push value a onto the IP stack |
| 0x15   | POPIP    | Pop a value from IP stack and jump to it |
| 0x16   | DROPIP   | Pop and discard a value from the IP stack |

## Input/Output

| Opcode | Mnemonic | Description |
|--------|----------|-------------|
| 0x7    | IN       | Read one byte from stdin, push as integer |
| 0x8    | OUT      | Pop a value, write to stdout as a byte |
| 0x11   | OUTNUM   | Pop a value, write to stdout as a number |

## Instruction Encoding

Each instruction is encoded as a 32-bit word. Instructions with immediate values (like PUSH) use the next 32-bit word as the operand.

## Examples

### Addition Example

```
PUSH 5    ; Push value 5 onto stack
PUSH 3    ; Push value 3 onto stack
ADD       ; Pop 3, pop 5, push 8
```

Stack after execution: [8]

### Function Call Example

```
PUSHIP 28  ; Push return address (28) onto IP stack
PUSH 40    ; Push function address
JMP        ; Jump to function
; Execution continues here after function returns

; Function at address 40:
; ... function code ...
POPIP      ; Return to caller
```

### Loop Example

```
PUSH 0     ; Initialize counter
PUSH 10    ; Push loop limit
; Loop start at address 8:
SWAP       ; Swap counter and limit
DUP        ; Duplicate counter
OUTNUM     ; Print counter
PUSH 1     ; Push increment value
ADD        ; Increment counter
SWAP       ; Swap counter and limit
DUP        ; Duplicate counter
PUSH 8     ; Push loop start address
SWAP       ; Swap
JNZ        ; Jump to loop start if counter != 0
DROP       ; Clean up stack
DROP
```

## Halt Instruction

The VM doesn't have a dedicated HALT instruction. Instead, a program halts by jumping to the current address:

```
PUSH <current_address>
JMP
```

This pattern is recognized by the VM as a termination request. 