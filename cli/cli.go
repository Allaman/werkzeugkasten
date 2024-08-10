package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// will be overwritten in release pipeline
var Version = "dev"

type CliConfig struct {
	Accessible  bool
	Category    string
	Debug       bool
	DownloadDir string
	Tools       bool
	Categories  bool
	ToolList    toolList
}
type toolList []string

func (s *toolList) String() string {
	return strings.Join(*s, ", ")
}

func (s *toolList) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func Cli() CliConfig {
	var cliFlags CliConfig
	var toolList toolList
	helpFlag := flag.Bool("help", false, "Print help message")
	versionFlag := flag.Bool("version", false, "Print version")
	debugFlag := flag.Bool("debug", false, "Enable debug output")
	accessibleFlag := flag.Bool("accessible", false, "Enable accessibility mode for interactive use")
	downloadDirFlag := flag.String("dir", ".", "Where to download the tools")
	listToolsFlag := flag.Bool("tools", false, "Print all available tools")
	listCategoriesFlag := flag.Bool("categories", false, "Print all categories and tool count")
	listByCategoriesFlag := flag.String("category", "", "List tools by category")
	flag.Var(&toolList, "tool", "Specify multiple tools to install programmatically (e.g., -tool kustomize -tool task)")
	flag.Parse()
	if *helpFlag {
		fmt.Println("Usage: werkzeugkasten [flags]")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		os.Exit(0)
	}
	if *versionFlag {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *listToolsFlag {
		cliFlags.Tools = true
	}
	if *listCategoriesFlag {
		cliFlags.Categories = true
	}
	if *debugFlag {
		cliFlags.Debug = true
	}
	if *accessibleFlag {
		cliFlags.Accessible = true
	}
	if *downloadDirFlag != "" {
		cliFlags.DownloadDir = *downloadDirFlag
	}
	if *listByCategoriesFlag != "" {
		cliFlags.Category = *listByCategoriesFlag
	}
	cliFlags.ToolList = []string{}
	if len(toolList) > 0 {
		cliFlags.ToolList = toolList
	}
	return cliFlags
}
