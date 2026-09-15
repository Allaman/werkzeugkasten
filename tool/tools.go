package tool

import (
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/goccy/go-yaml"
)

//go:embed tools.yaml
var toolsYAML []byte

type ToolData struct {
	Tools map[string]Tool `yaml:"tools"`
}

type Tool struct {
	Identifier   string   `yaml:"identifier"`
	Tag          string   `yaml:"tag"`
	Categories   []string `yaml:"categories"`
	Description  string   `yaml:"description"`
	AssetFilters []string `yaml:"asset_filters"`
	File         string   `yaml:"file"`
}

func CreateToolData() (ToolData, error) {
	var tools ToolData
	err := yaml.Unmarshal(toolsYAML, &tools)
	if err != nil {
		return ToolData{}, err
	}

	// Overwrite tags based with ENV variables
	// WK_TOOL_NAME_TAG, e.g. WK_KUSTOMIZE_TAG=v5.3.0
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key := pair[0]
		value := pair[1]

		if !strings.HasPrefix(key, "WK_") || !strings.HasSuffix(key, "_TAG") {
			continue
		}
		// the tool name itself may contain underscores, e.g. WK_PROCESS_COMPOSE_TAG
		tool := strings.ToLower(strings.TrimSuffix(strings.TrimPrefix(key, "WK_"), "_TAG"))
		t, ok := tools.Tools[tool]
		if !ok {
			slog.Warn("ignoring environment variable for unknown tool", "var", key, "tool", tool)
			continue
		}
		slog.Debug("overwriting tag", "tool", tool, "tag", value)
		// don't modify a copy of the struct with tools.Tools[tool].Tag
		t.Tag = value
		tools.Tools[tool] = t
	}
	return tools, nil
}

func execEget(workingDir string, tool Tool, system string) ([]byte, error) {
	egetBin := EgetPath()
	tag := tool.Tag
	name := tool.Identifier
	// eget runs with workingDir as its cwd, so --to must be absolute
	// to not resolve relative to workingDir a second time.
	workingDir, err := filepath.Abs(workingDir)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(workingDir); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("download directory %s does not exist, create it or pass an existing one with -dir", workingDir)
	}
	cmd := exec.Command(egetBin, "-q", name, "--to", workingDir)
	if tag != "" {
		cmd = exec.Command(egetBin, "-q", "-t", tag, name, "--to", workingDir)
	}
	if system != "" {
		cmd.Args = append(cmd.Args, "-s", system)
	}
	if len(tool.AssetFilters) > 0 {
		for _, af := range tool.AssetFilters {
			cmd.Args = append(cmd.Args, fmt.Sprintf("--asset=%s", af))
		}
	}
	if tool.File != "" {
		cmd.Args = append(cmd.Args, "--file", tool.File)
	}
	cmd.Dir = workingDir
	slog.Debug("executing command", "cmd", cmd, "wd", cmd.Dir, "env", cmd.Env, "args", cmd.Args)
	out, err := cmd.CombinedOutput()
	return out, err
}

// splitSystem splits an eget system string ("linux/amd64") into os and arch.
// An empty system resolves to the system werkzeugkasten runs on.
func splitSystem(system string) (string, string) {
	if system == "" {
		return runtime.GOOS, runtime.GOARCH
	}
	operatingSystem, arch, _ := strings.Cut(system, "/")
	return operatingSystem, arch
}

func DownloadToolWithEget(workingdir string, tool Tool, system string) error {
	operatingSystem, arch := splitSystem(system)
	// special handling for identifiers containing ARCH, e.g. helm
	tool.Identifier = strings.Replace(tool.Identifier, "ARCH", arch, 1)
	tool.Identifier = strings.Replace(tool.Identifier, "OSNAME", operatingSystem, 1)
	tag := "latest"
	if tool.Tag != "" {
		tag = tool.Tag
	}
	slog.Debug("downloading tool", "tool", tool.Identifier, "tag", tag, "system", system)
	out, err := execEget(workingdir, tool, system)
	if err != nil {
		slog.Debug("could not download tool", "tool", tool.Identifier, "error", err, "out", string(out))
		return err
	}
	return nil
}

func SortTools(tools ToolData) []string {
	sortedTools := make([]string, 0, len(tools.Tools))
	for k := range tools.Tools {
		sortedTools = append(sortedTools, k)
	}
	slices.Sort(sortedTools)
	return sortedTools
}

func GetCategories(tools ToolData) map[string]int {
	categories := make(map[string]int, 0)
	for _, t := range tools.Tools {
		for _, c := range t.Categories {
			categories[c] = categories[c] + 1
		}
	}
	return categories
}

func PrintCategories(categories map[string]int) {
	sortedCategories := make([]string, 0, len(categories))
	for k := range categories {
		sortedCategories = append(sortedCategories, k)
	}
	slices.Sort(sortedCategories)
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Name\tCount")
	for _, c := range sortedCategories {
		fmt.Fprintf(w, "%s\t%d\n", c, categories[c])
	}
	w.Flush()
}

func GetToolsByCategory(category string, tools ToolData) ToolData {
	var toolsFound ToolData
	toolsFound.Tools = make(map[string]Tool, 0)
	lowerCategory := strings.ToLower(category)
	for k, t := range tools.Tools {
		lowerCategories := make([]string, len(t.Categories))
		for i, v := range t.Categories {
			lowerCategories[i] = strings.ToLower(v)
		}
		if slices.Contains(lowerCategories, lowerCategory) {
			toolsFound.Tools[k] = t
		}
	}
	return toolsFound
}

func PrintTools(tools ToolData) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Key\tURL\tDescription")
	sortedTools := SortTools(tools)
	for _, tool := range sortedTools {
		identifier := tools.Tools[tool].Identifier
		url := fmt.Sprintf("https://github.com/%s", identifier)
		// handle packages that are not installed from GitHub
		if strings.HasPrefix(identifier, "https") {
			url = tools.Tools[tool].Identifier
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", tool, url, tools.Tools[tool].Description)
	}
	w.Flush()
}
