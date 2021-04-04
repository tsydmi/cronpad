package main

import (
	"fmt"
	"github.com/ts-dmitry/cronpad/backend/cmd"
	"os"
)

func main() {
	err := cmd.RunApp()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
