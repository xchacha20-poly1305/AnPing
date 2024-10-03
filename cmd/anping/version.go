package main

import (
	"fmt"
	"os"
)

var Version = "unknown"

func printVersion() {
	fmt.Println("AnPing:", Version)
	os.Exit(0)
}
