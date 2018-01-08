package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/relnod/calcgo/interpreter"
)

func main() {
	flag.Parse()

	result, errors := interpreter.Interpret(flag.Arg(0))
	if errors != nil {
		fmt.Println("Errors have occured:")
		for _, err := range errors {
			fmt.Println(err)
		}

		os.Exit(1)
	}

	fmt.Println(result)
}
