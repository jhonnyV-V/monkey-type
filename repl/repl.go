package repl

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"mokey-type/compiler"
	"mokey-type/lexer"
	"mokey-type/object"
	"mokey-type/parser"
	"mokey-type/vm"

	"github.com/chzyer/readline"
)

const (
	GREEN  = "\033[32m"
	BLUE   = "\033[36m"
	YELLOW = "\033[33m"
	RESET  = "\033[0m"
	PROMPT = GREEN + ">> " + RESET
)

const MONKEY_FACE = YELLOW + `
     w  c(..)o   (
      \__(-)    __)
          /\   (
         /(_)___)
         w /|
          | \
          m  m
` + RESET

func Start(out io.Writer) {
	vim := flag.Bool("vim", false, "activates vim mode")
	flag.Parse()

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	for {

		rl, err := readline.New(PROMPT)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		rl.SetVimMode(*vim)

		input, err := rl.Readline()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		line := strings.TrimSuffix(input, "\n")
		if line == "exit" {
			os.Exit(0)
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(p.Errors())
			continue
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err = comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "!Woops compiling bytecode failed\n error:\n \t%s\n", err)
			continue
		}

		machine := vm.NewWithGlobalsStore(comp.Bytecode(), globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "!Woops executing bytecode failed\n error:\n \t%s\n", err)
			continue
		}
		stackTop := machine.LastPopedStackElement()
		fmt.Fprintln(os.Stderr, stackTop.Inspect())
	}
}

func printParserErrors(errors []string) {
	io.WriteString(os.Stderr, MONKEY_FACE)
	io.WriteString(os.Stderr, "Woops! We ran into some monkey business here!\n")
	io.WriteString(os.Stderr, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(os.Stderr, "\t"+msg+"\n")
	}
}
