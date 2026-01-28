# Installation

## Go Install (Recommended)

The simplest way to install render is using `go install`:

```bash
go install github.com/wernerstrydom/render/cmd/render@latest
```

This requires Go 1.24 or later. The binary will be installed to `$GOPATH/bin` (or `$HOME/go/bin` if `GOPATH` is not set).

## Build from Source

Clone the repository and build:

```bash
git clone https://github.com/wernerstrydom/render.git
cd render
go build -o bin/render ./cmd/render
```

Move the binary to a directory in your PATH:

```bash
# Linux/macOS
sudo mv bin/render /usr/local/bin/

# Or add to your personal bin directory
mv bin/render ~/bin/
```

## Using Make

The project includes a Makefile for common operations:

```bash
git clone https://github.com/wernerstrydom/render.git
cd render
make build
```

The binary will be created at `bin/render`.

## Verify Installation

```bash
render --help
```

You should see the help text with available options and examples.

## Shell Completion

render supports shell completion through Cobra. Generate completion scripts:

```bash
# Bash
render completion bash > /etc/bash_completion.d/render

# Zsh
render completion zsh > "${fpath[1]}/_render"

# Fish
render completion fish > ~/.config/fish/completions/render.fish

# PowerShell
render completion powershell > render.ps1
```

## Man Pages

Generate and install man pages:

```bash
# Generate
render gen man ./man

# View locally
man ./man/render.1

# Install system-wide (Linux/macOS)
sudo mkdir -p /usr/local/share/man/man1
sudo cp ./man/render.1 /usr/local/share/man/man1/
sudo mandb  # Update man database (Linux)
```
