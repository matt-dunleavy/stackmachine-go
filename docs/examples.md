# Example Programs

The Stack Machine includes several example programs in the `bin/testdata` directory that demonstrate various features and programming techniques.

## Programs Overview

| Program         | Description                                       |
|-----------------|---------------------------------------------------|
| yo.src          | Simple "Yo" output program                        |
| hello.src       | Classic "Hello, World!" program                   |
| forward-goto.src| Demonstrates forward jumps and labels             |
| func.src        | Shows function call implementation                |
| core-test.src   | Tests core VM operations                          |
| core.src        | More extensive core functionality tests           |
| fib.src         | Fibonacci sequence calculator                     |

## Example Walkthrough

### yo.src

A minimal program that prints "Yo" followed by a newline.

```
'Y OUT
'o OUT
'\n OUT
HALT
```

### hello.src

The classic "Hello, World!" program, demonstrating character output.

```
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

### forward-goto.src

Demonstrates forward jumps and conditional branches with label references.

```
  PUSH 5      ; Counter
loop:
  DUP         ; Duplicate counter for output
  OUTNUM      ; Print counter
  PUSH ' OUT  ; Print space
  PUSH 1      ; Value to decrement by
  SUB         ; Subtract 1 from counter
  DUP         ; Duplicate for conditional
  PUSH &done  ; Push address for conditional jump
  SWAP        ; Swap for JZ
  JZ          ; Jump if zero
  PUSH &loop  ; Loop address for jump
  JMP         ; Jump back to loop
done:
  DROP        ; Clean up stack
  HALT
```

Output: `5 4 3 2 1 `

### func.src

Demonstrates function call implementation using the instruction pointer stack.

```
; Main program
  'M OUT     ; Print "Main"
  'a OUT
  'i OUT
  'n OUT
  ':' OUT
  ' ' OUT

  PUSHIP &return  ; Push return address
  PUSH &function  ; Push function address
  JMP             ; Call function

return:
  '\n OUT    ; Print newline after return
  HALT       ; End program

; Function
function:
  'F OUT     ; Print "Function"
  'u OUT
  'n OUT
  'c OUT
  't OUT
  'i OUT
  'o OUT
  'n OUT
  POPIP      ; Return to caller
```

Output: `Main: Function`

### core-test.src

A simple test of core arithmetic and stack operations.

```
  PUSH 2      ; Push first operand
  PUSH 3      ; Push second operand
  ADD         ; Add: 2+3=5
  OUTNUM      ; Print result: 5
  PUSH ' OUT  ; Print space
  
  PUSH 10     ; Push first operand
  PUSH 4      ; Push second operand
  SUB         ; Subtract: 10-4=6
  OUTNUM      ; Print result: 6
  PUSH '\n OUT
  HALT
```

Output: `5 6`

### core.src

More extensive tests of the VM's core functionality, including arithmetic, logic operations, and stack manipulation.

This program tests:
- Basic arithmetic (ADD, SUB)
- Logical operations (AND, OR, XOR, NOT)
- Stack manipulation (DUP, SWAP, ROL3)
- Memory operations (LOAD, STOR)
- Conditional jumps

### fib.src

A more complex program that calculates and prints Fibonacci numbers, demonstrating recursion, function calls, and parameter passing.

```
; Main program - calculate and print the first 10 Fibonacci numbers
  PUSH 0      ; First Fibonacci number: F(0) = 0
  OUTNUM
  PUSH ' OUT
  
  PUSH 1      ; Second Fibonacci number: F(1) = 1
  OUTNUM
  PUSH ' OUT
  
  PUSH 2      ; Start counter at 2
  PUSH 10     ; End at 10
  
fib_loop:
  DUP         ; Duplicate counter for Fibonacci calculation
  PUSHIP &after_fib  ; Return address after fib call
  PUSH &fib   ; Function address
  JMP         ; Call fib function
  
after_fib:
  OUTNUM      ; Print result
  PUSH ' OUT  ; Print space
  
  SWAP        ; Get counter and limit
  DUP         ; Duplicate counter
  PUSH 1      ; Increment value
  ADD         ; Increment counter
  SWAP        ; Swap counter and limit
  DUP2        ; Duplicate both values
  SUB         ; counter - limit
  PUSH &done  ; Target for jump if done
  SWAP        ; Swap for JZ
  JZ          ; Jump if counter >= limit
  PUSH &fib_loop  ; Target for jump back to loop
  JMP         ; Continue loop
  
done:
  DROP        ; Clean up stack
  DROP
  HALT        ; End program

; Fibonacci function: F(n) = F(n-1) + F(n-2) for n >= 2
fib:
  DUP         ; Duplicate argument
  PUSH 2      ; Compare with 2
  SUB         ; n-2
  PUSH &recurse  ; Address for recursion
  SWAP        ; Swap for JNZ
  JNZ         ; If n-2 != 0, need recursion
  
  ; Base cases: F(0) = 0, F(1) = 1
  PUSH 0      ; Value for F(0)
  SWAP        ; Swap with n
  PUSH 1      ; For comparison
  SUB         ; n-1
  PUSH &case_one  ; Jump target for case 1
  SWAP        ; Swap for JZ
  JZ          ; If n-1 == 0, it's case 1
  DROP        ; It's case 0, drop n-1
  POPIP       ; Return F(0) = 0
  
case_one:
  DROP        ; Drop n-1
  PUSH 1      ; Return F(1) = 1
  POPIP       ; Return
  
recurse:
  ; Calculate F(n-1)
  DUP         ; Duplicate n
  PUSH 1      ; Value to subtract
  SUB         ; n-1
  PUSHIP &after_n_1  ; Return address
  PUSH &fib   ; Function address
  JMP         ; Call fib(n-1)
  
after_n_1:
  ; Save F(n-1) and calculate F(n-2)
  SWAP        ; Swap F(n-1) with n
  PUSH 2      ; Value to subtract
  SUB         ; n-2
  PUSHIP &after_n_2  ; Return address
  PUSH &fib   ; Function address
  JMP         ; Call fib(n-2)
  
after_n_2:
  ; Add F(n-1) + F(n-2)
  ADD         ; F(n) = F(n-1) + F(n-2)
  POPIP       ; Return result
```

Output: `0 1 1 2 3 5 8 13 21 34`

## Running the Examples

You can run these examples using the interpret command:

```bash
stackmachine interpret bin/testdata/hello.src
```

Or compile and run them:

```bash
stackmachine compile bin/testdata/fib.src
stackmachine run fib.bin
```

## Learning from Examples

These examples demonstrate several important programming techniques:

1. **Basic I/O**: Using `OUT` and `OUTNUM` for output, `IN` for input
2. **Control Flow**: Using `JMP`, `JZ`, and `JNZ` for loops and conditionals
3. **Function Calls**: Using `PUSHIP` and `POPIP` for subroutines
4. **Stack Manipulation**: Managing the stack with `DUP`, `SWAP`, `ROL3`
5. **Memory Usage**: Using `LOAD` and `STOR` for memory operations
6. **Algorithm Implementation**: Implementing common algorithms like Fibonacci

By studying these examples, you can learn how to implement more complex programs using the Stack Machine's simple instruction set. 