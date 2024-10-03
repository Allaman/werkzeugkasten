<h1 align="center">Werkzeugkasten 🧰</h1>

<div align="center">
  <p>
    <img src="https://github.com/Allaman/werkzeugkasten/actions/workflows/release.yaml/badge.svg" alt="Release"/>
    <img src="https://img.shields.io/github/repo-size/Allaman/werkzeugkasten" alt="size"/>
    <img src="https://img.shields.io/github/issues/Allaman/werkzeugkasten" alt="issues"/>
    <img src="https://img.shields.io/github/last-commit/Allaman/werkzeugkasten" alt="last commit"/>
    <img src="https://img.shields.io/github/license/Allaman/werkzeugkasten" alt="license"/>
    <img src="https://img.shields.io/github/v/release/Allaman/werkzeugkasten?sort=semver" alt="last release"/>
  </p>
</div>

_Conveniently download your favorite binaries (currently 111 supported)!_

[![screenshot](https://s1.gifyu.com/images/SBpu4.png)](https://gifyu.com/image/SBpu4)

<details>
<summary>Open a tool's README within werkzeugkasten</summary>

[![readme.png](https://s11.gifyu.com/images/SBpu5.png)](https://gifyu.com/image/SBpu5)

</details>

<details>
<summary>List categories and tool count</summary>

`werkzeugkasten -categories`

[![categories.png](https://s1.gifyu.com/images/SBpuN.png)](https://gifyu.com/image/SBpuN)

</details>

<details>
<summary>List tools in category "Text"</summary>

`werkzeugkasten -category text`

[![category.png](https://s1.gifyu.com/images/SBpu9.png)](https://gifyu.com/image/SBpu9)

</details>

<details>
<summary>Install tools in non-interactive mode</summary>

`werkzeugkasten -dir /tmp -tool age -tool kustomize`

[![install.png](https://s1.gifyu.com/images/SBpuv.png)](https://gifyu.com/image/SBpuv)

</details>

From time to time, I need to connect to containers and VMs to troubleshoot them. These systems typically only have the necessary tools for their specific purpose and nothing else. Additionally, there is no root account available, so installing tools through a package manager is not an option. Furthermore, some tools are either not available as a package or the packaged version is outdated.

This is where Werkzeugkasten comes in. You simply need to download the werkzeugkasten binary onto your system, and from that point on, there are no additional requirements, particularly the need for root permissions.

## Get Werkzeugkasten

Unfortunately, a tool to download the werkzeugkasten binary is required. It is possible to download files via bash and `/dev/tcp` **only**, but I couldn't figure out how to handle the redirect from Github when accessing a release URL.

with curl

```sh
VERSION=$(curl -s https://api.github.com/repos/allaman/werkzeugkasten/releases/latest | grep tag_name | cut -d '"' -f 4)
curl -sLo werkzeugkasten https://github.com/Allaman/werkzeugkasten/releases/download/${VERSION}/werkzeugkasten_${VERSION}_$(uname -s)_$(uname -m)
```

with wget

```sh
VERSION=$(wget -qO - https://api.github.com/repos/allaman/werkzeugkasten/releases/latest | grep tag_name | cut -d '"' -f 4)
wget -qO werkzeugkasten https://github.com/Allaman/werkzeugkasten/releases/download/${VERSION}/werkzeugkasten_${VERSION}_$(uname -s)_$(uname -m)
```

```sh
chmod +x werkzeugkasten
./werkzeugkasten
```

You could also integrate werkzeugkasten in your golden (Docker) image. ⚠️ Keep possible security implications in mind.

## How it works

Werkzeugkasten is basically a wrapper around the excellent [eget](https://github.com/zyedidia/eget) that does the heavy lifting and is responsible for downloading the chosen tools. Eget itself is downloaded as binary via `net/http` call and decompression/extraction logic.

The awesome [charmbracelet](https://github.com/charmbracelet) tools [bubbletea](https://github.com/charmbracelet/bubbletea), [glamour](https://github.com/charmbracelet/glamour), and [lipgloss](https://github.com/charmbracelet/lipgloss) are used for a modern look and feel. By default, the latest release of a tool is downloaded (see [Configuration](#configuration)).

## What Werkzeugkasten is not

Werkzeugkasten is not intended to replace package managers (such as apt, brew, ...) or configuration management tools (such as Ansible, ...).

## Usage

```
❯ werkzeugkasten -help
Usage: werkzeugkasten [flags]
Flags:
  -categories
        Print all categories and tool count
  -category string
        List tools by category
  -debug
        Enable debug output
  -dir string
        Where to download the tools (default ".")
  -help
        Print help message
  -tool value
        Specify multiple tools to install programmatically (e.g., -tool kustomize -tool task)
  -tools
        Print all available tools
  -update
        Self-update
  -version
        Print version
```

Werkzeugkasten supports an **interactive** mode and a **non-interactive** mode.

- `werkzeugkasten` will start in interactive mode where you select your tools you want to install from a searchable list.

- `werkzeugkasten -tool age -tool kustomize` will download age and kustomize.

- `werkzeugkasten -tools` will print all available tools.

- `werkzeugkasten -categories` will print all available categories.

- `werkzeugkasten -category network` will print all available tools in the "network" category.

## Configuration

Besides CLI flags, further configuration is possible with environment variables.

Set a tool's version/tag explicitly:

```sh
export WK_<TOOL_NAME>_<TAG>=1.33.7
export WK_KUSTOMIZE_TAG=v5.3.0`
```

Set a GitHub token to get more than the 60 API calls per hour limit:

```sh
export EGET_GITHUB_TOKEN=<token>
```
