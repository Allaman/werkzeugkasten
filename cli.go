package main

import (
	"flag"
	"fmt"
	"os"
)

type cliConfig struct {
	debug      bool
	accessible bool
}

func cli() cliConfig {
	var cliFlags cliConfig
	helpFlag := flag.Bool("help", false, "Print help message")
	versionFlag := flag.Bool("version", false, "Print version")
	debugFlag := flag.Bool("debug", false, "Enable debug output")
	accessibleFlag := flag.Bool("accessible", false, "Enable accessibility mode")
	flag.Parse()
	if *helpFlag {
		fmt.Println("Usage: werkzeugkasten [flags]")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		os.Exit(0)
	}
	if *versionFlag {
		logger.Print(version)
		os.Exit(0)
	}
	if *debugFlag {
		cliFlags.debug = true
	}
	if *accessibleFlag {
		cliFlags.accessible = true
	}
	return cliFlags
}
