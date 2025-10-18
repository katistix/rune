# rune - simple stack based virtual machine written in Golang

#### goal

this is a learning project. i want to implement a simple stack based virtual machine with its own instruction set. basically trying to understand how computers work at the lowest level by building one from scratch.

#### what it can do right now

- **basic arithmetic**: add, subtract, multiply, divide
- **bitwise operations**: and, or, xor, not, bit shifting
- **logic operations**: equals, not equals, greater than, less than  
- **control flow**: jumps, conditional jumps based on stack values
- **memory operations**: load and store bytes in 64KB of RAM
- **host bindings**: call back into Go functions for I/O and system stuff

the VM uses an 8-bit data stack and can execute programs loaded into its memory. it's got about 20 opcodes so far, which is enough to do some interesting things.

#### how to run it

```bash
go run cmd/cli/main.go
```

right now there's a small test program hardcoded in `main.go` that pushes some values, does math, and prints the result. not fancy but it works.

#### inspiration

heavily inspired by [Uxn by Hundred Rabbits](https://100r.co/site/uxn.html). their philosophy of minimal computing really resonates with me.

#### current limitations

- you have to write bytecode by hand (arrays of hex values)
- no proper debugging tools
- limited to 8-bit values on the stack
- memory operations are pretty basic

#### future goals

right now, the only way to program the VM is by hardcoding the *bytecode* in the `main.go` file. but in the future i will implement:

- **assembler**: so you can write assembly instead of raw bytecode
- **better debugging**: stack inspection, step-through execution
- **more host bindings**: file I/O, networking maybe?
- **small compiler**: for a simple high-level language that compiles to rune bytecode

this is just a side project i work on when i want to understand how things work under the hood. no rush, just learning and having fun with it.

