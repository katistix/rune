package main

import (
	"fmt"

	"github.com/katistix/rune/internal/bindings"
	"github.com/katistix/rune/internal/vmcore"
)

// 8

func main() {
	fmt.Printf("[rune] hello cli\n\n")

	testBytecode := []byte{

		vmcore.OP_PUSH,
		'a',
		vmcore.OP_PUSH,
		5,

		vmcore.OP_ADD,

		vmcore.OP_HOSTCALL,
		0x01,

		vmcore.OP_HALT,
	}

	vm := vmcore.NewRuneVM(testBytecode)
	// vm.LoadSpec(bindings.CoreSpec)
	vm.LoadSpec(bindings.StandardSpec)

	vm.Run()
}
