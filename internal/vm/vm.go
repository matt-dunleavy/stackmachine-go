package vm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ErrorCallback is a function type for error handling
type ErrorCallback func(msg string)

// VM represents a stack machine
type VM struct {
	stack     []int32       // Data stack
	stackIP   []int32       // Instruction pointer stack
	labels    []Label       // Code labels
	memSize   int           // Memory size in words
	memory    []int32       // VM memory
	ip        int32         // Instruction pointer
	in        *bufio.Reader // Input stream
	out       io.Writer     // Output stream
	running   bool          // VM running state
	errorFunc ErrorCallback // Error callback function
}

// NewMachine creates a new machine instance with default settings
func NewMachine(errorCallback ErrorCallback) *VM {
	return NewMachineWithSize(1000*1024, os.Stdout, os.Stdin, errorCallback)
}

// NewMachineWithSize creates a new machine instance with specified memory size and I/O
func NewMachineWithSize(memorySize int, out io.Writer, in io.Reader, errorCallback ErrorCallback) *VM {
	m := &VM{
		stack:     make([]int32, 0),
		stackIP:   make([]int32, 0),
		labels:    make([]Label, 0),
		memSize:   memorySize,
		memory:    make([]int32, memorySize),
		ip:        0,
		in:        bufio.NewReader(in),
		out:       out,
		running:   true,
		errorFunc: errorCallback,
	}
	m.Reset()
	return m
}

// Clone creates a copy of the machine
func (m *VM) Clone(errorCallback ErrorCallback) *VM {
	clone := &VM{
		stack:     make([]int32, len(m.stack)),
		stackIP:   make([]int32, len(m.stackIP)),
		labels:    make([]Label, len(m.labels)),
		memSize:   m.memSize,
		memory:    make([]int32, m.memSize),
		ip:        m.ip,
		in:        m.in,
		out:       m.out,
		running:   m.running,
		errorFunc: errorCallback,
	}

	copy(clone.stack, m.stack)
	copy(clone.stackIP, m.stackIP)
	copy(clone.labels, m.labels)
	copy(clone.memory, m.memory)

	return clone
}

// Reset initializes the machine state
func (m *VM) Reset() {
	// Clear memory with NOP instructions
	for i := range m.memory {
		m.memory[i] = int32(NOP)
	}
	m.stack = m.stack[:0] // Clear stack
	m.ip = 0
}

// Error reports an error via the error callback
func (m *VM) Error(msg string) {
	if m.errorFunc != nil {
		m.errorFunc(msg)
	}
}

// Push pushes a value onto the data stack
func (m *VM) Push(n int32) {
	m.stack = append(m.stack, n)
}

// PushIP pushes a value onto the instruction pointer stack
func (m *VM) PushIP(n int32) {
	m.stackIP = append(m.stackIP, n)
}

// PopIP pops a value from the instruction pointer stack
func (m *VM) PopIP() int32 {
	if len(m.stackIP) == 0 {
		m.Error("POP empty IP stack")
		return 0
	}
	n := m.stackIP[len(m.stackIP)-1]
	m.stackIP = m.stackIP[:len(m.stackIP)-1]
	return n
}

// Pop pops a value from the data stack
func (m *VM) Pop() int32 {
	if len(m.stack) == 0 {
		m.Error("POP empty stack")
		return 0
	}
	n := m.stack[len(m.stack)-1]
	m.stack = m.stack[:len(m.stack)-1]
	return n
}

// CheckBounds checks if an address is within memory bounds
func (m *VM) CheckBounds(n int32, msg string) bool {
	if n < 0 || int(n) >= m.memSize {
		m.Error(msg)
		return false
	}
	return true
}

// Next advances the instruction pointer to the next word
func (m *VM) Next() {
	m.ip += 4 // Assuming 4-byte words
	if m.ip < 0 {
		m.Error("IP < 0")
	}
	if int(m.ip) >= m.memSize {
		m.ip = 0 // Wrap around
	}
}

// Prev moves the instruction pointer to the previous word
func (m *VM) Prev() {
	if m.ip == 0 {
		m.Error("prev() reached zero")
		return
	}
	m.ip -= 4 // Assuming 4-byte words
}

// Load loads an opcode into memory and advances IP
func (m *VM) Load(op Op) {
	m.memory[m.ip] = int32(op)
	m.Next()
}

// LoadInt loads a 32-bit integer into memory and advances IP
func (m *VM) LoadInt(n int32) {
	m.memory[m.ip] = n
	m.Next()
}

// Run executes the program from a given address
func (m *VM) Run(startAddr int32) int {
	m.ip = startAddr
	m.running = true

	for m.running {
		m.Exec(Op(m.memory[m.ip]))
	}

	return 0 // Exit code
}

// Size returns the size of the program in memory
func (m *VM) Size() int32 {
	// Find the end of the program by scanning backwards until a non-NOP is found
	for i := m.memSize - 1; i >= 0; i-- {
		if m.memory[i] != int32(NOP) {
			return int32(i + 4) // Add word size
		}
	}
	return 0
}

// Cur returns the current memory content at the instruction pointer
func (m *VM) Cur() int32 {
	return m.memory[m.ip]
}

// Pos returns the current instruction pointer position
func (m *VM) Pos() int32 {
	return m.ip
}

// AddLabel adds a label to the program
func (m *VM) AddLabel(name string, pos int32) {
	// Remove trailing colon if present
	name = strings.TrimSuffix(name, ":")
	m.labels = append(m.labels, NewLabel(name, pos))
}

// GetLabelAddress returns the address of a label by name
func (m *VM) GetLabelAddress(name string) int32 {
	upperName := strings.ToUpper(name)

	// Special label "HERE" returns current position
	if upperName == "HERE" {
		return m.ip
	}

	for _, label := range m.labels {
		if strings.ToUpper(label.Name) == upperName {
			return label.Pos
		}
	}

	return -1 // Not found
}

// IsRunning returns the running state of the machine
func (m *VM) IsRunning() bool {
	return m.running
}

// SetOutput sets the output writer
func (m *VM) SetOutput(out io.Writer) {
	m.out = out
}

// SetInput sets the input reader
func (m *VM) SetInput(in io.Reader) {
	m.in = bufio.NewReader(in)
}

// SetMem sets a memory value at a specific address
func (m *VM) SetMem(addr int32, val int32) {
	if m.CheckBounds(addr, "set_mem out of bounds") {
		m.memory[addr] = val
	}
}

// GetMem gets a memory value from a specific address
func (m *VM) GetMem(addr int32) int32 {
	if m.CheckBounds(addr, "get_mem out of bounds") {
		return m.memory[addr]
	}
	return 0
}

// WordSize returns the size of a word in bytes
func (m *VM) WordSize() int32 {
	return 4 // 32-bit words
}

// LoadHalt loads a halt instruction sequence
func (m *VM) LoadHalt() {
	m.Load(PUSH)
	m.LoadInt(m.ip + 4) // Next instruction
	m.Load(JMP)
}

// LoadImage loads a program image from a reader
func (m *VM) LoadImage(r io.Reader) error {
	m.Reset()
	buf := make([]byte, 4)

	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if n != 4 {
			return fmt.Errorf("incomplete read: got %d bytes, expected 4", n)
		}

		op := int32(buf[0]) | int32(buf[1])<<8 | int32(buf[2])<<16 | int32(buf[3])<<24
		m.Load(Op(op))
	}

	m.ip = 0
	return nil
}

// SaveImage saves the program to a writer
func (m *VM) SaveImage(w io.Writer) error {
	size := m.Size()
	buf := make([]byte, 4)

	for i := int32(0); i < size; i += 4 {
		val := m.memory[i]
		buf[0] = byte(val)
		buf[1] = byte(val >> 8)
		buf[2] = byte(val >> 16)
		buf[3] = byte(val >> 24)

		if _, err := w.Write(buf); err != nil {
			return err
		}
	}

	return nil
}

// Exec executes a single instruction
func (m *VM) Exec(op Op) {
	switch op {
	case NOP:
		m.InstrNOP()
	case ADD:
		m.InstrAdd()
	case SUB:
		m.InstrSub()
	case AND:
		m.InstrAnd()
	case OR:
		m.InstrOr()
	case XOR:
		m.InstrXor()
	case NOT:
		m.InstrNot()
	case COMPL:
		m.InstrCompl()
	case IN:
		m.InstrIn()
	case OUT:
		m.InstrOut()
	case LOAD:
		m.InstrLoad()
	case STOR:
		m.InstrStor()
	case PUSH:
		m.InstrPush()
	case DROP:
		m.InstrDrop()
	case PUSHIP:
		m.InstrPushIP()
	case POPIP:
		m.InstrPopIP()
	case DROPIP:
		m.InstrDropIP()
	case JZ:
		m.InstrJZ()
	case JMP:
		m.InstrJmp()
	case JNZ:
		m.InstrJNZ()
	case DUP:
		m.InstrDup()
	case SWAP:
		m.InstrSwap()
	case ROL3:
		m.InstrRol3()
	case OUTNUM:
		m.InstrOutNum()
	default:
		m.Error(fmt.Sprintf("Unknown instruction: %d", op))
	}
}

// Instruction implementations

func (m *VM) InstrNOP() {
	m.Next()
}

func (m *VM) InstrAdd() {
	b := m.Pop()
	a := m.Pop()
	m.Push(a + b)
	m.Next()
}

func (m *VM) InstrSub() {
	b := m.Pop()
	a := m.Pop()
	m.Push(b - a) // Note: Order matches C++ implementation
	m.Next()
}

func (m *VM) InstrAnd() {
	b := m.Pop()
	a := m.Pop()
	m.Push(a & b)
	m.Next()
}

func (m *VM) InstrOr() {
	b := m.Pop()
	a := m.Pop()
	m.Push(a | b)
	m.Next()
}

func (m *VM) InstrXor() {
	b := m.Pop()
	a := m.Pop()
	m.Push(a ^ b)
	m.Next()
}

func (m *VM) InstrNot() {
	a := m.Pop()
	if a == 0 {
		m.Push(1)
	} else {
		m.Push(0)
	}
	m.Next()
}

func (m *VM) InstrCompl() {
	m.Push(^m.Pop())
	m.Next()
}

func (m *VM) InstrIn() {
	b, err := m.in.ReadByte()
	if err != nil {
		m.Push(0) // EOF or error
	} else {
		m.Push(int32(b))
	}
	m.Next()
}

func (m *VM) InstrOut() {
	b := byte(m.Pop())
	m.out.Write([]byte{b})
	m.Next()
}

func (m *VM) InstrOutNum() {
	fmt.Fprintf(m.out, "%d", m.Pop())
	m.Next()
}

func (m *VM) InstrLoad() {
	addr := m.Pop()
	if m.CheckBounds(addr, "LOAD") {
		m.Push(m.memory[addr])
	}
	m.Next()
}

func (m *VM) InstrStor() {
	addr := m.Pop()
	val := m.Pop()
	if m.CheckBounds(addr, "STOR") {
		m.memory[addr] = val
	}
	m.Next()
}

func (m *VM) InstrJmp() {
	addr := m.Pop()
	if !m.CheckBounds(addr, "JMP") {
		m.Next()
		return
	}

	// Check if halting (jumping to current address)
	if addr == m.ip {
		m.running = false
	} else {
		m.ip = addr
	}
}

func (m *VM) InstrJZ() {
	pred := m.Pop()
	addr := m.Pop()

	if pred != 0 {
		m.Next()
	} else {
		if m.CheckBounds(addr, "JZ") {
			m.ip = addr
		} else {
			m.Next()
		}
	}
}

func (m *VM) InstrDrop() {
	m.Pop()
	m.Next()
}

func (m *VM) InstrPopIP() {
	addr := m.PopIP()
	if m.CheckBounds(addr, "POPIP") {
		m.ip = addr
	} else {
		m.Next()
	}
}

func (m *VM) InstrDropIP() {
	m.PopIP()
	m.Next()
}

func (m *VM) InstrJNZ() {
	pred := m.Pop()
	addr := m.Pop()

	if pred == 0 {
		m.Next()
	} else {
		if m.CheckBounds(addr, "JNZ") {
			m.ip = addr
		} else {
			m.Next()
		}
	}
}

func (m *VM) InstrPush() {
	m.Next()
	m.Push(m.memory[m.ip])
	m.Next()
}

func (m *VM) InstrPushIP() {
	m.Next()
	m.PushIP(m.memory[m.ip])
	m.Next()
}

func (m *VM) InstrDup() {
	a := m.Pop()
	m.Push(a)
	m.Push(a)
	m.Next()
}

func (m *VM) InstrSwap() {
	b := m.Pop()
	a := m.Pop()
	m.Push(b)
	m.Push(a)
	m.Next()
}

func (m *VM) InstrRol3() {
	c := m.Pop()
	b := m.Pop()
	a := m.Pop()
	m.Push(b)
	m.Push(c)
	m.Push(a)
	m.Next()
}
