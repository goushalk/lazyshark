package main

import (
	"flag"
	"os"

	"github.com/charmbracelet/log"

	"packetB/internal/tui"
)

func main() {
	filePath := flag.String("f", "", "file path to read")
	flag.Parse()

	if *filePath == "" {
		log.Error("file not provided")
		os.Exit(1)
	}

	// Analyzer part

	if err := tui.StartTUI(*filePath); err != nil {
		log.Error(err)
		os.Exit(1)
	}

}
