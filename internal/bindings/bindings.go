package bindings

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/katistix/rune/internal/vmcore"
)

var CoreSpec = map[byte]vmcore.Binding{
	0x01: func(vm *vmcore.RuneVM) {
		value := vm.Pop()
		fmt.Printf("%c", value) // print top of stack as char
	},
	0x02: func(vm *vmcore.RuneVM) {
		value := vm.Pop()
		fmt.Println(value) // print top of stack as int
	},
}

var StandardSpec = map[byte]vmcore.Binding{
	// ---- CONSOLE IO ----
	0x01: func(vm *vmcore.RuneVM) { // print char
		v := vm.Pop()
		fmt.Printf("%c", v)
	},
	0x02: func(vm *vmcore.RuneVM) { // print int
		v := vm.Pop()
		fmt.Println(v)
	},
	0x03: func(vm *vmcore.RuneVM) { // read char
		reader := bufio.NewReader(os.Stdin)
		c, _, _ := reader.ReadRune()
		vm.Push(byte(c))
	},

	// ---- RANDOMNESS / TIME ----
	0x10: func(vm *vmcore.RuneVM) { // push random byte
		vm.Push(byte(rand.Intn(256)))
	},
	0x11: func(vm *vmcore.RuneVM) { // push seconds since epoch low byte
		vm.Push(byte(time.Now().Unix() & 0xFF))
	},
	0x12: func(vm *vmcore.RuneVM) { // sleep ms
		ms := vm.Pop()
		time.Sleep(time.Duration(ms) * time.Millisecond)
	},

	// ---- SOUND ----
	0x20: func(vm *vmcore.RuneVM) { // beep
		switch runtime.GOOS {
		case "linux", "darwin":
			// uses 'printf "\a"' in the shell
			exec.Command("sh", "-c", "printf '\\a'").Run()
		case "windows":
			// uses powershell to play a beep
			exec.Command("powershell", "-c", "[console]::beep(750,300)").Run()
		default:
			fmt.Println("Beep not supported on this OS")
		}
	},
}
