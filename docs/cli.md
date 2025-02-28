# Command-Line Interface

The Stack Machine provides a set of command-line tools for compiling, running, and analyzing programs.

## Command Overview

| Command     | Description                                      | Alias |
| ----------- | ------------------------------------------------ | ----- |
| compile     | Compile source code to bytecode                  | smc   |
| run         | Execute compiled bytecode                        | smr   |
| interpret   | Compile and execute source code in one step      | sm    |
| disassemble | Convert bytecode back to human-readable assembly | smd   |

## Common Features

All commands share these common behaviors:

- If no filenames are provided, input is read from standard input
- Most commands support the `-h` or `--help` flag for usage information
- Multiple files can be processed in sequence

## Command Details

### compile

Compiles stack machine source code to bytecode.

```bash
smg compile [file...]
```

**Options:**

- None specific to this command

**Examples:**

```bash
# Compile a source file to bytecode
smg compile program.src
# Output: program.bin

# Compile from stdin to out.bin
cat program.src | smg compile
```

### run

Executes compiled bytecode files.

```bash
smg run [file...]
```

**Options:**

- `-h`, `--help`: Show instruction set information

**Examples:**

```bash
# Run a compiled bytecode file
smg run program.bin

# Run from stdin
cat program.bin | smg run
```

### interpret

Compiles and executes source code in a single step.

```bash
smg interpret [file...]
```

**Options:**

- None specific to this command

**Examples:**

```bash
# Interpret a source file
smg interpret program.src

# Interpret from stdin
cat program.src | smg interpret
```

### disassemble

Converts bytecode back to human-readable assembly code.

```bash
smg disassemble [file...]
```

**Options:**

- None specific to this command

**Examples:**

```bash
# Disassemble a bytecode file
smg disassemble program.bin

# Disassemble from stdin
cat program.bin | smg disassemble
```

## Input/Output Behavior

### File Extensions

While not strictly required, the following conventions are recommended:

- `.src`: Source code files
- `.bin`: Compiled bytecode files

### Default Output Files

When compiling, if the output filename is not specified:

- For a source file named `program.src`, the output will be `program.bin`
- For input from stdin, the output will be `out.bin`

### Standard Input/Output

All commands can read from standard input and write to standard output:

```bash
# Pipeline example
smg compile program.src | smg disassemble

# Using - as a filename to indicate stdin
smg compile - < program.src | smg run -
```

## Error Handling

When errors occur during compilation or execution:

- Error messages are printed to stderr
- The program exits with a non-zero status code
- Error messages include the filename and error description

## Environment

The Stack Machine commands do not depend on environment variables or configuration files. All behavior is controlled through command-line arguments.
