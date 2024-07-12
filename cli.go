package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type cliConfig struct {
	accessible  bool
	debug       bool
	downloadDir string
	list        bool
	toolList    toolList
}
type toolList []string

func (s *toolList) String() string {
	return strings.Join(*s, ", ")
}

func (s *toolList) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func cli() cliConfig {
	var cliFlags cliConfig
	var toolList toolList
	helpFlag := flag.Bool("help", false, "Print help message")
	versionFlag := flag.Bool("version", false, "Print version")
	debugFlag := flag.Bool("debug", false, "Enable debug output")
	accessibleFlag := flag.Bool("accessible", false, "Enable accessibility mode for interactive use")
	downloadDirFlag := flag.String("installDir", ".", "Where to install the tools")
	listFlag := flag.Bool("list", false, "Print all available tools")
	flag.Var(&toolList, "tool", "Specify multiple tools to install programmatically (e.g., -tool kustomize -tool task)")
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
	if *listFlag {
		cliFlags.list = true
	}
	if *debugFlag {
		cliFlags.debug = true
	}
	if *accessibleFlag {
		cliFlags.accessible = true
	}
	if *downloadDirFlag != "" {
		cliFlags.downloadDir = *downloadDirFlag
	}
	cliFlags.toolList = []string{}
	if len(toolList) > 0 {
		cliFlags.toolList = toolList
	}
	return cliFlags
}
