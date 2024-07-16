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
	tools       bool
	categories  bool
	theme       string
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
	downloadDirFlag := flag.String("dir", ".", "Where to download the tools")
	listToolsFlag := flag.Bool("tools", false, "Print all available tools")
	listCategoriesFlag := flag.Bool("categories", false, "Print all available categories")
	themeFlag := flag.String("theme", "catppuccin", "Set theme for interactive mode")
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
	if *listToolsFlag {
		cliFlags.tools = true
	}
	if *listCategoriesFlag {
		cliFlags.categories = true
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
	cliFlags.theme = *themeFlag
	return cliFlags
}
