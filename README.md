# repo-lookup

Small cli tool to look up and clone git repositories

![Alt Text](demo.gif)

## Installation

Prerequisites:

- Go
- \[Optional] [github-cli](https://github.com/cli/cli/tree/trunk) if you want to use the -f flag

_Unix_

```bash
go build -o ~/.local/bin/rlu
```

_Windows_

I don't know where you put binaries on Windows. I'm sure you can figure it out.

## Usage

```bash
rlu [-f] <repo-name>
```

The -f flag is for forking the repo instead of just cloning.
