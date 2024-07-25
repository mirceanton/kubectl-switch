# kube-switcher

`kube-switcher` is a CLI tool for managing and switching Kubernetes contexts. It allows you to interactively or non-interactively switch between different Kubernetes contexts defined in multiple kubeconfig files.

## Features

- Interactive Context Selection: Use a terminal-based UI to select and switch contexts.
- Non-Interactive Mode: Switch directly to a specified context by name, suitable for scripting.
- Config File Management: Handles multiple kubeconfig files and manages context switching seamlessly.

## Installation

### Prerequisites

- Go 1.22.5 or higher
- `kubectl` (for context verification)

### Building from Source

- Clone the repository:

```sh
git clone https://github.com/your-username/kube-switcher.git
cd kube-switcher
```

- Build the binary:

```sh
go build -o kube-switcher
```

- (Optional) Move the binary to a directory in your `$PATH`:

```sh
sudo mv kube-switcher /usr/local/bin/
```

### Installation via Release

- Download the latest release from the releases page.
- Extract the archive and move the binary to a directory in your `$PATH`.

## Usage

### Configuration

Before using `kube-switcher`, set the `KUBESWITCHER_CONFIG_DIR` environment variable to the directory containing your kubeconfig files:

```sh
export KUBESWITCHER_CONFIG_DIR="/path/to/your/kubeconfig/files"
```

### Interactive Mode

To interactively select and switch to a Kubernetes context, run:

```sh
kube-switcher switch
```

You will be prompted to choose a context from the list of available contexts.

### Non-Interactive Mode

To switch to a specific context by name without interaction, use the `--context` flag:

```sh
kube-switcher switch --context <context_name>
```

Replace <context_name> with the name of the context you want to switch to.

### Examples

#### Interactive Selection:

```sh
kube-switcher switch
```

Output:

```sh
Use the arrow keys to navigate: ↓ ↑ → ←
? Select Kubernetes Context:
  ▸ my-cluster
    another-cluster
```

#### Non-Interactive Selection:

```sh
kube-switcher switch --context my-cluster
```

Output:

```sh
Switched to context 'my-cluster' from file '/home/user/.kube/configs/my-cluster.yaml'
```

## Error Handling

- If multiple contexts with the same name are found in different files, kube-switcher will report an error and exit.
- If the `KUBESWITCHER_CONFIG_DIR` environment variable is not set or is incorrect, kube-switcher will provide an appropriate error message.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.

- Fork the repository,
- Create a new branch for your changes,
- Make your changes and commit them,
- Submit a pull request describing your changes.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
