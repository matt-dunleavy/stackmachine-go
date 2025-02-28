# Compiler Documentation

## Overview

The Stack Machine compiler translates human-readable assembly code into bytecode that can be executed by the virtual machine. The compiler processes source code through several phases: lexical analysis (tokenization), parsing, and code generation.

## Source Language

The Stack Machine assembly language is line-oriented and token-based. Each instruction or directive is separated by whitespace. The language has the following elements:

### Instructions

Instructions are written as uppercase mnemonics:

```
PUSH 10
ADD
DUP
OUTNUM
```

Instructions correspond directly to VM opcodes and are case-insensitive (e.g., `PUSH`, `Push`, and `push` are equivalent).

### Literals

The compiler supports several types of literals:

1. **Numeric literals**: Decimal integers
   ```
   PUSH 42
   ```

2. **Character literals**: Characters enclosed in single quotes
   ```
   PUSH 'A'    ; Pushes the ASCII value of 'A' (65)
   ```

3. **Escape sequences**: Special characters in character literals
   ```
   PUSH '\n'   ; Pushes the value for a newline character (10)
   PUSH '\t'   ; Pushes the value for a tab character (9)
   PUSH '\r'   ; Pushes the value for a carriage return (13)
   PUSH '\0'   ; Pushes the value for a null character (0)
   ```

### Labels

Labels define named positions in the code and can be referenced by other instructions:

```
start:     ; Define a label named "start"
  PUSH 1
  OUTNUM
  PUSH &start  ; Reference to the label "start"
  JMP
```

Labels are defined by a name followed by a colon (`:`) and are referenced with an ampersand (`&`) prefix.

### Comments

Comments begin with a semicolon (`;`) and continue to the end of the line:

```
PUSH 10   ; This is a comment
ADD       ; Add the top two values on the stack
```

### Function Calls

Function calls are implemented using labels and the instruction pointer stack:

```
; Main program
PUSHIP &return  ; Push return address
PUSH &function  ; Push function address
JMP             ; Jump to function
return:         ; Execution continues here after return

; Function
function:
  ; ... function body ...
  POPIP         ; Return to caller
```

## Compilation Process

### Tokenization

The compiler breaks down source code into tokens, separating instructions, literals, labels, and comments. Whitespace is used as a delimiter between tokens.

### Parsing

Tokens are categorized as:
- Instructions (matching known opcodes)
- Labels (ending with a colon)
- Label references (starting with &)
- Literals (numeric or character)
- Comments (starting with a semicolon)

### Code Generation

For each token, the compiler:
1. Ignores comments
2. Records label positions
3. Generates bytecode for instructions
4. Handles forward references to labels
5. Manages function call patterns

### Forward References

The compiler supports forward references to labels not yet defined. When a reference to a not-yet-defined label is encountered, the compiler:
1. Generates the correct instruction
2. Records the location of the reference
3. Uses a placeholder value
4. Resolves all forward references at the end of compilation

### Error Handling

During compilation, the compiler reports errors such as:
- Unknown instructions
- Invalid label names
- References to undefined labels
- Malformed character literals

## Special Directives

### HALT

While `HALT` is not a native VM instruction, the compiler recognizes it as a directive to generate a halt sequence:

```
HALT
```

Compiles to:
```
PUSH <current_address + 4>
JMP
```

## Byte Code Format

The compiler generates a binary file consisting of 32-bit words. Each word is stored in little-endian format:

1. Instructions are encoded as a single 32-bit word
2. Instructions with immediate values (like PUSH) use two words:
   - The instruction opcode
   - The immediate value

## Usage Examples

### Basic Compilation

```
PUSH 65    ; ASCII 'A'
OUT        ; Print 'A'
PUSH 10    ; ASCII newline
OUT        ; Print newline
HALT       ; End program
```

### Loop with Labels

```
  PUSH 10      ; Loop counter
start:
  DUP          ; Duplicate counter
  OUTNUM       ; Print counter
  PUSH 1       ; Decrement value
  SUB          ; Decrement counter
  DUP          ; Duplicate for check
  PUSH &end    ; End address
  SWAP         ; Swap for JZ
  JZ           ; Jump if zero
  PUSH &start  ; Loop address
  JMP          ; Jump to start
end:
  DROP         ; Clean up stack
  HALT         ; End program
```

### Function Definition and Call

```
; Main program
  PUSH 5       ; Argument for factorial
  PUSHIP &ret  ; Push return address
  PUSH &fact   ; Push function address
  JMP          ; Call factorial
ret:
  OUTNUM       ; Print result
  HALT         ; End program

; Factorial function (n! = n * (n-1)!)
fact:
  DUP          ; Duplicate argument
  PUSH 1       ; Compare with 1
  SUB          ; n-1
  PUSH &recur  ; Recursion branch address
  SWAP         ; Swap for JNZ
  JNZ          ; If n-1 != 0, recurse
  DROP         ; Drop n-1, keep n for base case
  POPIP        ; Return to caller
recur:
  PUSHIP &mul  ; Push return address for recursion
  SWAP         ; Swap return addr and n-1
  PUSH &fact   ; Address of factorial
  JMP          ; Recursive call
mul:
  MUL          ; Multiply n * (n-1)!
  POPIP        ; Return to caller
``` 