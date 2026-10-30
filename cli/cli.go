package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// will be overwritten in release pipeline
var Version = "dev"

type CLIConfig struct {
	Category    string
	Debug       bool
	NoColor     bool
	Update      bool
	DownloadDir string
	Tools       bool
	Categories  bool
	System      string
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

// CLI parses the flags. The returned config is usable even on error so callers can set up logging first.
func CLI() (CLIConfig, error) {
	var cliFlags CLIConfig
	var toolList toolList
	helpFlag := flag.Bool("help", false, "Print help message")
	versionFlag := flag.Bool("version", false, "Print version")
	updateFlag := flag.Bool("update", false, "Self-update")
	debugFlag := flag.Bool("debug", false, "Enable debug output")
	noColorFlag := flag.Bool("no-color", false, "Disable colored log output (also respects NO_COLOR)")
	downloadDirFlag := flag.String("dir", ".", "Where to download the tools")
	listToolsFlag := flag.Bool("tools", false, "Print all available tools")
	listCategoriesFlag := flag.Bool("categories", false, "Print all categories and tool count")
	listByCategoriesFlag := flag.String("category", "", "List tools by category")
	systemFlag := flag.String("system", "", "Target system to download for in the form os/arch (e.g., linux/amd64), defaults to the current system")
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
	if *updateFlag {
		cliFlags.Update = true
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
	if *noColorFlag {
		cliFlags.NoColor = true
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
	if *systemFlag != "" {
		if o, a, ok := strings.Cut(*systemFlag, "/"); !ok || o == "" || a == "" {
			return cliFlags, fmt.Errorf("-system must be in the form os/arch, e.g. linux/amd64, got %q", *systemFlag)
		}
		cliFlags.System = *systemFlag
	}
	return cliFlags, nil
}
