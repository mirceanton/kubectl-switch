# kube-switcher

`kube-switcher` is a CLI tool for managing and switching Kubernetes contexts. It allows you to interactively or non-interactively switch between different Kubernetes contexts defined in multiple kubeconfig files.

## Why `kube-switcher`?

`kube-switcher` is an alternative to tools like `kubectx` and `kubie`. It has been created because I feel both of those fall short in certain regards:

- `kubectx` assumes all your contexts are defined in a single config file. Yes, there is some hackery you can do to your `KUBECONFIG` environment variable to make it work with multiple files, but it is a sub-par workflow in my opinion and I never really liked it
- `kubie` spawns a new shell when you use it to change contexts, making it practically impossible to integrate into scripts or taskfile automation. Also I consider it to be too complicated of a solution for what is fundamentally a really simple problem

What I wanted was something much simpler conceptually: I just want to dump all my kubeconfig files in a single directory and then, have my tool parse them and "physically" move over the config file to `.kube/config` such that it is also persistent between different shells.

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
export KUBESWITCHER_CONFIG_DIR="~/.kube/configs/"
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
