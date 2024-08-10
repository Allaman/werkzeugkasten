package tool

import (
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
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
	// Target       string   `yaml:"target"`
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

		if strings.HasPrefix(key, "WK_") {
			trimmedKey := strings.TrimPrefix(key, "WK_")
			splittedKey := strings.Split(trimmedKey, "_")
			if len(splittedKey) != 2 {
				slog.Warn("ignoring environment variable", "var", key)
				continue
			}
			tool := strings.ToLower(splittedKey[0])
			field := strings.ToLower(splittedKey[1])
			if field != "tag" {
				slog.Warn("ignoring malformed environment variable", "var", key)
				continue
			}
			// TODO: overwrite more fields dynamically this way
			if t, ok := tools.Tools[tool]; ok {
				slog.Debug("overwriting tag", "tool", tool, "tag", value)
				// tools.Tools[tool].Tag = value not working because
				// when modifying the fields of the struct obtained from the map, you are modifying a copy of the struct!
				t.Tag = value
				tools.Tools[tool] = t
			}
		}
	}
	return tools, nil
}

func execEget(workingDir string, tool Tool) ([]byte, error) {
	tag := tool.Tag
	name := tool.Identifier
	cmd := exec.Command("./eget", "-q", name, "--to", workingDir)
	if tag != "" {
		cmd = exec.Command("./eget", "-q", "-t", tag, name, "--to", workingDir)
	}
	if len(tool.AssetFilters) > 0 {
		for _, af := range tool.AssetFilters {
			cmd.Args = append(cmd.Args, fmt.Sprintf("--asset=%s", af))
		}
	}
	if tool.File != "" {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--file=\"%s\"", tool.File))
	}
	cmd.Dir = workingDir
	slog.Debug("executing command", "cmd", cmd, "wd", cmd.Dir, "env", cmd.Env)
	out, err := cmd.CombinedOutput()
	return out, err
}

func DownloadToolWithEget(workingdir string, tool Tool) error {
	tool.Identifier = strings.Replace(tool.Identifier, "ARCH", runtime.GOARCH, 1)
	tool.Identifier = strings.Replace(tool.Identifier, "OSNAME", runtime.GOOS, 1)
	tag := "latest"
	if tool.Tag != "" {
		tag = tool.Tag
	}
	slog.Info("downloading tool", "tool", tool.Identifier, "tag", tag)
	out, err := execEget(workingdir, tool)
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
