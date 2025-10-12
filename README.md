# kubectl-switch

`kubectl-switch` is a command-line tool for managing and switching between multiple Kubernetes configuration files located in the same directory. It simplifies the process of selecting a Kubernetes context from multiple kubeconfig files and updating the active configuration or namespace.

Just dump all your `kubeconfigs` into a single dir and let `kubectl-switch` manage them for you!

## Features

- **Multiple kubeconfig files**: Manage multiple kubeconfig files in a single directory without merging them
- **Context & namespace switching**: Switch between contexts and namespaces from multiple config files
- **Interactive & non-interactive modes**: Select from a list or specify directly as an argument (with tab completion support!)

## Why `kubectl-switch`?

`kubectl-switch` is an alternative to tools like `kubectx`, `kubens` and `kubie`. It has been created because I feel all of those fall short in certain regards:

- `kubectx` and `kubens` assume all your contexts are defined in a single config file. Yes, there is some hackery you can do to your `KUBECONFIG` environment variable to make it work with multiple files, but it is (in my opinion) a sub-par workflow and I never really liked it
- `kubie` spawns a new shell when you use it to change contexts, making it practically impossible to integrate into scripts or taskfile automation. Also I consider it to be too complicated of a solution for what is fundamentally a really simple problem

What I wanted was something very simple conceptually: I just want to dump all my `kubeconfig` files in a single directory and then have my tool parse them and "physically" move over the config file to `.kube/config` (or whatever is configured in my `KUBECONFIG` env var) such that it is also persistent between different shells. Here is where `kubectl-switch` comes in!

## Installation

### Install via go install

Install the latest stable version directly using Go:

```bash
# Install the latest version
go install github.com/mirceanton/kubectl-switch/v2@latest

# Or install a specific version
go install github.com/mirceanton/kubectl-switch/v2@v2.2.6
```

> **Note**: Versions prior to v2.2.6 cannot be installed via `go install` due to module path issues. If you need an older version, please use one of the alternative installation methods below.

### Download Precompiled Binaries

Precompiled binaries are available for various platforms. You can download the latest release from the [GitHub Releases page](https://github.com/mirceanton/kubectl-switch/releases/latest).

1. Download the appropriate binary for your system and extract the archive.
2. Make the extracted binary executable:

   ```bash
   chmod +x kubectl-switch
   ```

3. Move the binary to a directory in your PATH:

   ```bash
   mv kubectl-switch /usr/local/bin/kubectl-switch
   ```

### Install via Homebrew

1. Add the tap

   ```bash
   brew tap mirceanton/taps
   ```

2. Install `kubectl-switch`

   ```bash
   brew install kubectl-switch
   ```

### Running via Docker

`kubectl-switch` is also available as a Docker container:

```bash
docker pull ghcr.io/mirceanton/kubectl-switch
```

To use it with your local kubeconfig files:

```bash
docker run -v ~/.kube:/root/.kube -v /path/to/kubeconfigs:/kubeconfigs \
    -e KUBECONFIG_DIR=/kubeconfigs \
    ghcr.io/mirceanton/kubectl-switch context
```

### Build from Source

1. Clone the repository:

   ```bash
   git clone https://github.com/mirceanton/kubectl-switch
   cd kubectl-switch
   ```

2. Build the tool:

   ```bash
   go build -o kubectl-switch
   ```

3. Move the binary to your PATH:

   ```bash
   mv kubectl-switch /usr/local/bin/kubectl-switch
   ```

## Usage

### Context Command

The `context` (or `ctx`) subcommand is used to switch between Kubernetes contexts (think of `kubectx`):

```bash
# Interactive mode - select context from a list
kubectl-switch context

# Switch to a specific context
kubectl-switch ctx my-context
```

### Namespace Command

The `namespace` (or `ns`) subcommand is used to switch the current namespace (think of `kubens`):

```bash
# Interactive mode - select namespace from a list
kubectl-switch namespace

# Switch to a specific namespace
kubectl-switch ns kube-system
```

### Quickly Switch to Previous Configuration

Switch back to the previous configuration:

```bash
kubectl-switch -
```

### Usage with kubectl plugin

When installed as a kubectl plugin, you can use it directly with the `kubectl` command:

```sh
# These commands are equivalent
kubectl-switch ctx
kubectl switch ctx
```

## Configuration

`kubectl-switch` uses Viper for configuration management, allowing you to configure the tool via command-line flags or environment variables, with flags taking precedence.

### Configuration Options

| Option               | Flag               | Environment Variable | Default            | Description                                                       |
| -------------------- | ------------------ | -------------------- | ------------------ | ----------------------------------------------------------------- |
| Kubeconfig Directory | `--kubeconfig-dir` | `KUBECONFIG_DIR`     | `~/.kube/configs/` | Directory containing your kubeconfig files                        |
| Kubeconfig           | `--kubeconfig`     | `KUBECONFIG`         | `~/.kube/config`   | Path to the currently active kubeconfig file.                     |
| Log Level            | `--log-level`      | `LOG_LEVEL`          | `info`             | Logging verbosity (trace, debug, info, warn, error, fatal, panic) |
| Log Format           | `--log-format`     | `LOG_FORMAT`         | `text`             | Log output format (text, json)                                    |

## Shell Completion

The `completion` subcommand generates shell completion scripts:

```bash
# Generate completions for bash
kubectl-switch completion bash > /etc/bash_completion.d/kubectl-switch

# Generate completions for zsh
kubectl-switch completion zsh > ~/.zsh/completion/_kubectl-switch

# Generate completions for fish
kubectl-switch completion fish > ~/.config/fish/completions/kubectl-switch.fish

# Generate completions for powershell
kubectl-switch completion powershell > ~/kubectl-switch.ps1
```

## Contributing

Contributions are welcome! Please fork the repository, make your changes, and submit a pull request.

## License

This project is licensed under the MIT License. See the [`LICENSE`](./LICENSE) file for details.
