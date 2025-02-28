package vm

import (
	"strings"
)

// Op represents a machine instruction opcode
type Op int32

// Instruction set opcodes
const (
	NOP     Op = iota
	ADD        // pop a, pop b, push a + b
	SUB        // pop a, pop b, push a - b
	AND        // pop a, pop b, push a & b
	OR         // pop a, pop b, push a | b
	XOR        // pop a, pop b, push a ^ b
	NOT        // pop a, push !a
	IN         // read one byte from stdin, push as word on stack
	OUT        // pop one word and write to stream as one byte
	LOAD       // pop a, push word read from address a
	STOR       // pop a, pop b, write b to address a
	JMP        // pop a, goto a
	JZ         // pop a, pop b, if a == 0 goto b
	PUSH       // push next word
	DUP        // duplicate word on stack
	SWAP       // swap top two words on stack
	ROL3       // rotate top three words on stack once left, (a b c) -> (b c a)
	OUTNUM     // pop one word and write to stream as number
	JNZ        // pop a, pop b, if a != 0 goto b
	DROP       // remove top of stack
	PUSHIP     // push a in IP stack
	POPIP      // pop IP stack to current IP, effectively performing a jump
	DROPIP     // pop IP, but do not jump
	COMPL      // pop a, push the complement of a
	NOP_END    // placeholder for end of enum; MUST BE LAST
)

// OpStr contains string representations of opcodes
var OpStr = []string{
	"NOP",
	"ADD",
	"SUB",
	"AND",
	"OR",
	"XOR",
	"NOT",
	"IN",
	"OUT",
	"LOAD",
	"STOR",
	"JMP",
	"JZ",
	"PUSH",
	"DUP",
	"SWAP",
	"ROL3",
	"OUTNUM",
	"JNZ",
	"DROP",
	"PUSHIP",
	"POPIP",
	"DROPIP",
	"COMPL",
	"NOP_END",
}

// String returns the string representation of an opcode
func (op Op) String() string {
	if op >= NOP && op < NOP_END {
		return OpStr[op]
	}
	return "<?>"
}

// FromString converts a string to an opcode
func FromString(s string) Op {
	upper := strings.ToUpper(s)
	for i, str := range OpStr {
		if upper == str {
			return Op(i)
		}
	}
	return NOP_END
}
