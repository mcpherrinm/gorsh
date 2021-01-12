package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mcpherrinm/gorsh/pkg/shell"
)

func main() {
	err := shell.New(os.Stdin, os.Stdout).Interact(context.Background())
	if err != nil {
		fmt.Printf("%v", err)
	}
}
