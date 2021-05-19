package cmd

import (
	"fmt"
	"os"
)

func CmdPrintln(a ...interface{}) (int, error) {
	return fmt.Println(a...)
}

func CmdPrintErrorln(a ...interface{}) (int, error) {
	return fmt.Fprintln(os.Stderr, a...)
}

func CmdPrettyPrintln(a ...interface{}) (int, error) {
	return fmt.Fprintln(os.Stdout, a...)
}
