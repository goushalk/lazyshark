package main

import (
	"flag"
	"os"

	"github.com/charmbracelet/log"

	"github.com/goushalk/lazyshark/internal/tui"
)

func main() {
	filePath := flag.String("f", "", "file path to read")
	flag.Parse()

	if *filePath == "" {
		log.Error("file not provided")
		os.Exit(1)
	}

	// TUI

	if err := tui.StartTUI(*filePath); err != nil {
		log.Error(err)
		os.Exit(1)
	}

}
