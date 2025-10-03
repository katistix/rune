package vmcore

import (
	"fmt"
)

// Opcodes
const (
	OP_NOP byte = iota
	OP_PUSH
	OP_POP

	// ARITHMETIC
	OP_ADD
	OP_SUB
	OP_DIV
	OP_MUL

	// BITWISE OPERATIONS
	OP_AND
	OP_OR
	OP_XOR
	OP_NOT
	OP_SHL
	OP_SHR

	// LOGIC
	OP_EQ
	OP_NEQ
	OP_GTH
	OP_LTH

	// CONTROL FLOW
	OP_JMP
	OP_JZ
	OP_JNZ
	// OP_CALL
	// OP_RETURN

	// MEMORY OPs
	OP_LOAD
	OP_STORE

	OP_HOSTCALL

	OP_HALT
)

type Binding func(vm *RuneVM)
type RuneVM struct {
	memory   [65536]byte // 64KiB RAM
	pc       uint16
	ds       []byte // data stack (8-bit values)
	running  bool
	bindings map[byte]Binding
}

func NewRuneVM(bytecode []byte) *RuneVM {
	vm := &RuneVM{
		pc:      0x0000,
		ds:      make([]byte, 0, 256),
		running: false,
	}

	// load the program at 0x0000
	copy(vm.memory[0:], bytecode)

	return vm
}

func (vm *RuneVM) Bind(id byte, fn Binding) {
	if vm.bindings == nil {
		vm.bindings = make(map[byte]Binding)
	}
	vm.bindings[id] = fn
}

func (vm *RuneVM) LoadSpec(spec map[byte]Binding) {
	for id, fn := range spec {
		vm.Bind(id, fn)
	}
}

func (vm *RuneVM) Fetch() byte {
	b := vm.memory[vm.pc]
	vm.pc++
	return b
}

func (vm *RuneVM) Fetch16() uint16 {
	lo := uint16(vm.Fetch()) // fetch low byte
	hi := uint16(vm.Fetch()) // fetch high byte
	return lo | (hi << 8)    // combine into 16-bit value
}

func (vm *RuneVM) Push(value byte) {
	vm.ds = append(vm.ds, value)
}

func (vm *RuneVM) Pop() byte {
	if len(vm.ds) == 0 {
		panic("data stack underflow")
	}

	value := vm.ds[len(vm.ds)-1] // get the value at the TOP
	vm.ds = vm.ds[:len(vm.ds)-1] // POP the value
	return value
}

func (vm *RuneVM) Run() {
	vm.running = true // start the execution

	for vm.running {
		op := vm.Fetch()

		switch op {
		case OP_NOP:
			// do nothing
		case OP_PUSH:
			value := vm.Fetch()
			vm.Push(value)
		case OP_POP:
			vm.Pop()
		// ARITHMETIC
		case OP_ADD:
			b := vm.Pop()
			a := vm.Pop()
			vm.Push(a + b)
		case OP_SUB:
			b := vm.Pop()
			a := vm.Pop()
			vm.Push(a - b)
		case OP_MUL:
			b := vm.Pop()
			a := vm.Pop()
			vm.Push(a * b)
		case OP_DIV:
			b := vm.Pop()
			a := vm.Pop()
			if b == 0 {
				panic("division by zero")
			}
			vm.Push(a / b)

		// BITWISE OPERATIONS
		case OP_AND:
			b := vm.Pop()
			a := vm.Pop()
			vm.Push(a & b)
		case OP_OR:
			b := vm.Pop()
			a := vm.Pop()
			vm.Push(a | b)
		case OP_XOR:
			b := vm.Pop()
			a := vm.Pop()
			vm.Push(a ^ b)
		case OP_NOT:
			a := vm.Pop()
			vm.Push(^a)
		case OP_SHL:
			a := vm.Pop()
			vm.Push(a << 1)
		case OP_SHR:
			a := vm.Pop()
			vm.Push(a >> 1)

		// LOGIC
		case OP_EQ:
			b := vm.Pop()
			a := vm.Pop()
			if a == b {
				vm.Push(1)
			} else {
				vm.Push(0)
			}
		case OP_NEQ:
			b := vm.Pop()
			a := vm.Pop()
			if a != b {
				vm.Push(1)
			} else {
				vm.Push(0)
			}
		case OP_GTH:
			b := vm.Pop()
			a := vm.Pop()
			if a > b {
				vm.Push(1)
			} else {
				vm.Push(0)
			}
		case OP_LTH:
			b := vm.Pop()
			a := vm.Pop()
			if a < b {
				vm.Push(1)
			} else {
				vm.Push(0)
			}
			// CONTROL FLOW
		case OP_JMP:
			addr := vm.Fetch16() // 16-bit jump
			vm.pc = addr

		case OP_JZ:
			addr := vm.Fetch16()
			v := vm.Pop()
			if v == 0 {
				vm.pc = addr
			}

		case OP_JNZ:
			addr := vm.Fetch16()
			v := vm.Pop()
			if v != 0 {
				vm.pc = addr
			}

		// MEMORY OPERATIONS
		case OP_LOAD:
			addr := vm.Fetch16()
			vm.Push(vm.memory[addr])

		case OP_STORE:
			addr := vm.Fetch16()
			value := vm.Pop()
			vm.memory[addr] = value

		case OP_HOSTCALL:
			id := vm.Fetch() // get binding ID
			if fn, ok := vm.bindings[id]; ok {
				fn(vm)
			} else {
				panic(fmt.Sprintf("unknown host binding %d", id))
			}
		case OP_HALT:
			vm.running = false
		default:
			panic(fmt.Sprintf("unknown opcode %02X at PC %04X", op, vm.pc-1))
		}
	}
}
