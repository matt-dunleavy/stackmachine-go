package compiler

import (
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/matt-dunleavy/stackmachine-go/internal/vm"
)

// Compiler translates source code to machine code
type Compiler struct {
	machine   *vm.VM
	vm        *vm.VM
	forwards  []vm.Label
	errorFunc func(string)
}

// NewCompiler creates a new compiler
func NewCompiler(errorCallback func(string)) *Compiler {
	machine := vm.NewMachine(errorCallback)
	return &Compiler{
		machine:   machine,
		vm:        machine,
		forwards:  make([]vm.Label, 0),
		errorFunc: errorCallback,
	}
}

// Error reports an error
func (c *Compiler) Error(msg string) {
	if c.errorFunc != nil {
		c.errorFunc(msg)
	}
}

// IsLabel checks if a token is a label (ends with colon)
func (c *Compiler) IsLabel(s string) bool {
	return len(s) > 0 && s[len(s)-1] == ':'
}

// IsComment checks if a token is a comment (starts with semicolon)
func (c *Compiler) IsComment(s string) bool {
	return len(s) > 0 && s[0] == ';'
}

// TokenToOp converts a token to an operation
func (c *Compiler) TokenToOp(s string) vm.Op {
	return vm.FromString(s)
}

// IsLiteral checks if a token is a literal value (not an operation or label)
func (c *Compiler) IsLiteral(s string) bool {
	if c.IsLabel(s) {
		return false
	}
	return c.TokenToOp(s) == vm.NOP_END
}

// IsNumber checks if a string is a number
func (c *Compiler) IsNumber(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return len(s) > 0
}

// IsChar checks if a token represents a character literal
func (c *Compiler) IsChar(s string) bool {
	if len(s) == 3 && s[0] == '\'' && s[2] == '\'' && s[1] != '\\' {
		return true
	}
	if len(s) == 4 && s[0] == '\'' && s[3] == '\'' && s[1] == '\\' &&
		(s[2] == 't' || s[2] == 'r' || s[2] == 'n' || s[2] == '0') {
		return true
	}
	return false
}

// ToOrd converts a character literal to its ASCII value
func (c *Compiler) ToOrd(s string) int32 {
	if len(s) == 3 { // 'x'
		return int32(s[1])
	}
	if len(s) == 4 { // '\x'
		switch s[2] {
		case 't':
			return int32('\t')
		case 'r':
			return int32('\r')
		case 'n':
			return int32('\n')
		case '0':
			return int32(0)
		}
	}
	c.Error("Unknown character literal: " + s)
	return 0
}

// IsLabelRef checks if a token is a label reference (starts with &)
func (c *Compiler) IsLabelRef(s string) bool {
	return len(s) > 0 && s[0] == '&'
}

// ToLiteral converts a string to a numeric literal
func (c *Compiler) ToLiteral(s string) int32 {
	if c.IsNumber(s) {
		val, _ := strconv.Atoi(s)
		return int32(val)
	}
	if c.IsChar(s) {
		return c.ToOrd(s)
	}
	return -1
}

// IsHalt checks if a token represents a halt instruction
func (c *Compiler) IsHalt(s string) bool {
	return s == "" || strings.ToUpper(s) == "HALT"
}

// CheckLabelName validates a label name
func (c *Compiler) CheckLabelName(label string) {
	if strings.ToUpper(label) == "HERE" {
		c.Error("Label is reserved: HERE")
	}
}

// CompileLabel compiles a label reference
func (c *Compiler) CompileLabel(label string) {
	address := c.vm.GetLabelAddress(label)

	c.vm.Load(vm.PUSH)

	// If label not found, mark it for update
	if address == -1 {
		c.CheckLabelName(label)
		c.forwards = append(c.forwards, vm.NewLabel(label, c.vm.Pos()))
	}

	c.vm.LoadInt(address)
}

// CompileFunctionCall compiles a function call
func (c *Compiler) CompileFunctionCall(function string) {
	// Return address is here plus four instructions
	c.vm.Load(vm.PUSHIP)
	c.vm.LoadInt(c.vm.Pos() + 4*c.vm.WordSize())

	// Push function destination address -- update it later
	c.vm.Load(vm.PUSH)
	c.forwards = append(c.forwards, vm.NewLabel(function, c.vm.Pos()))
	c.vm.LoadInt(-1) // Just push an arbitrary number

	// Jump to function
	c.vm.Load(vm.JMP)
}

// CompileLiteral compiles a literal value
func (c *Compiler) CompileLiteral(token string) {
	if c.IsLabelRef(token) {
		c.CompileLabel(token[1:])
		return
	}

	literal := c.ToLiteral(token)

	// Literals are pushed onto the stack
	if literal != -1 {
		c.vm.Load(vm.PUSH)
		c.vm.LoadInt(literal)
		return
	}

	// Unknown literals are treated as forward function calls
	c.CompileFunctionCall(token)
}

// ResolveForwards resolves forward references
func (c *Compiler) ResolveForwards() {
	for _, forward := range c.forwards {
		address := c.vm.GetLabelAddress(forward.Name)

		if address == -1 {
			c.Error("Code label not found: " + forward.Name)
		}

		// Update label jump to address
		c.vm.SetMem(forward.Pos, address)
	}
}

// CompileToken compiles a single token
// Returns false when compilation is finished
func (c *Compiler) CompileToken(s string, p *Parser) (bool, error) {
	if s == "" {
		c.vm.LoadHalt()
		c.ResolveForwards()
		return false, nil
	} else if c.IsHalt(s) {
		c.vm.LoadHalt()
	} else if c.IsComment(s) {
		p.SkipLine()
	} else if c.IsLiteral(s) {
		c.CompileLiteral(s)
	} else if c.IsLabel(s) {
		c.vm.AddLabel(s, c.vm.Pos())
	} else {
		op := c.TokenToOp(s)

		if op == vm.NOP_END {
			c.Error("Unknown operation: " + s)
		}

		c.vm.Load(op)
	}

	return true, nil
}

// GetProgram returns the compiled machine
func (c *Compiler) GetProgram() *vm.VM {
	return c.machine
}

// CompileSource compiles source code from a reader
func (c *Compiler) CompileSource(r io.Reader) error {
	parser := NewParser(r)

	for {
		token, err := parser.NextToken()
		if err != nil && err != io.EOF {
			return err
		}

		continueCompiling, err := c.CompileToken(token, parser)
		if err != nil {
			return err
		}

		if !continueCompiling || (err == io.EOF) {
			break
		}
	}

	return nil
}
