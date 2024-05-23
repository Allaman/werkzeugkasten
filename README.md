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

_Conveniently download your favorite binaries (currently 69 supported)!_

![](./screenshot.png)

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

## How it works

Werkzeugkasten is basically a wrapper around the excellent [eget](https://github.com/zyedidia/eget) that does the heavy lifting and is responsible for downloading the chosen tools. Eget itself is downloaded as binary via `net/http` call and decompression/extracting logic. The awesome [charmbracelet](https://github.com/charmbracelet) tools [huh](https://github.com/charmbracelet/huh), [log](https://github.com/charmbracelet/log), and [lipgloss](https://github.com/charmbracelet/lipgloss) are used for a modern look and feel. By default, the latest release of a tool is downloaded (see [Configuration](#configuration)).

## What Werkzeugkasten is not

Werkzeugkasten is not intended to replace package managers (such as apt, brew, ...) or configuration management tools (such as Ansible, ...). It is also not intended to be used non-interactively.

## Configuration

```
❯ werkzeugkasten -h
Usage of werkzeugkasten:
  -accessible
        Enable accessibility mode
  -debug
        Enable debug output
  -help
        Print help message
  -version
        Print version
```

Besides CLI flags, further configuration is possible with environment variables. Since Werkzeugkasten is designed to run on minimal systems, I cannot rely on having an editor available for writing configuration files.

Set a tool's version/tag explicitly:

```sh
export WK_<TOO_NAME>_<TAG>>=1.33.7
export WK_KUSTOMIZE_TAG=v5.3.0`
```

Set a GitHub token to get more than the 60 API calls per hour limit:

```sh
export EGET_GITHUB_TOKEN=<token>
```
