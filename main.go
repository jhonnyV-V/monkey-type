package main

import (
	"fmt"
	"mokey-type/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s!, This Is MonkeyType\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
