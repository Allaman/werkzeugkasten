# Werkzeugkasten

_Conveniently download your favorite binaries!_

From time to time, I need to connect to containers and VMs to troubleshoot them. These systems typically only have the necessary tools for their specific purpose and nothing else. Additionally, there is no root account available, so installing tools through a package manager is not an option. Furthermore, some tools are either not available in a package or the packaged version is outdated.

This is where Werkzeugkasten comes in. You simply need to download the werkzeugkasten binary onto the system, and from that point on, there are no additional requirements, particularly the need for root permissions.

## Get Werkzeugkasten

```sh
curl -sLo werkzeugkasten https://github.com/Allaman/werkzeugkasten/releases/download/0.9.0/werkzeugkasten_0.9.0_darwin_arm64
curl -sLo werkzeugkasten https://github.com/Allaman/werkzeugkasten/releases/download/0.9.0/werkzeugkasten_0.9.0_linux_arm64

curl -sLo werkzeugkasten https://github.com/Allaman/werkzeugkasten/releases/download/0.9.0/werkzeugkasten_0.9.0_darwin_amd64
curl -sLo werkzeugkasten https://github.com/Allaman/werkzeugkasten/releases/download/0.9.0/werkzeugkasten_0.9.0_linux_amd64
```

```sh
wget -qO werkzeugkasten https://github.com/Allaman/werkzeugkasten/releases/download/0.9.0/werkzeugkasten_0.9.0_darwin_arm64
wget -qO werkzeugkasten https://github.com/Allaman/werkzeugkasten/releases/download/0.9.0/werkzeugkasten_0.9.0_linux_arm64

wget -qO werkzeugkasten https://github.com/Allaman/werkzeugkasten/releases/download/0.9.0/werkzeugkasten_0.9.0_darwin_amd64
wget -qO werkzeugkasten https://github.com/Allaman/werkzeugkasten/releases/download/0.9.0/werkzeugkasten_0.9.0_linux_amd64
```

## How it works

Werkzeugkasten is basically a wrapper around the excellent [eget](https://github.com/zyedidia/eget) that does the heavy lifting and is responsible for downloading the chosen tools. Eget itself is downloaded as binary via `net/http` call and decompression/extracting logic. The awesome [charmbracelet](https://github.com/charmbracelet) tools [huh](https://github.com/charmbracelet/huh), [log](https://github.com/charmbracelet/log), and [lipgloss](https://github.com/charmbracelet/lipgloss) are used for a modern look and feel.

## What Werkzeugkasten is not

Werkzeugkasten is not intended to replace package managers (such as apt, brew, ...) or configuration management tools (such as Ansible, ...). It is also not intended to be used non-interactively.

## Configuration

Configuration can only be done using environment variables. Since Werkzeugkasten is designed to run on minimal systems, I cannot rely on having an editor available for writing configuration files.

Enable debug output:

```sh
export WK_DEBUG=true
```

Enable [accessibility](https://github.com/charmbracelet/huh?tab=readme-ov-file#accessibility) specifically designed for screen readers:

```sh
export WK_ACCESSIBLE=true
```

Overwrite tool version/tag defined in [tools.yaml](https://github.com/Allaman/werkzeugkasten/blob/main/tools.yaml):

```sh
export WK_<TOO_NAME>_<TAG>>=1.33.7
export WK_KUSTOMIZE_TAG=v5.3.0`
```

Set a GitHub token to get more than the 60 API calls per hour limit:

```sh
export EGET_GITHUB_TOKEN=<token>
```
