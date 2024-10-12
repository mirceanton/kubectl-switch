# kube-switcher

A simple tool to manage and switch between multiple Kubernetes configuration files. This tool allows you to dump all your kubeconfig files in a single directory and easily switch between them as well as configure the current namespace.

Available both as a standalone CLI or a `kubectl` plugin.

## Features

- Manage multiple kubeconfig files in a single directory.
- Easily switch between different Kubernetes contexts.
- Easily switch between different namespaces.
- Persistent context/namespace switching across different shells.
- Interactive prompts support fuzzy search.
- Non-destructive operations on your original kubeconfig files. `kube-switcher` will never edit those. It will only work with copies of them.

## Why `kube-switcher`?

`kube-switcher` is an alternative to tools like `kubectx`, `kubens` and `kubie`. It has been created because I feel all of those fall short in certain regards:

- `kubectx` assumes all your contexts are defined in a single config file. Yes, there is some hackery you can do to your `KUBECONFIG` environment variable to make it work with multiple files, but it is a sub-par workflow in my opinion and I never really liked it
- `kubie` spawns a new shell when you use it to change contexts, making it practically impossible to integrate into scripts or taskfile automation. Also I consider it to be too complicated of a solution for what is fundamentally a really simple problem

What I wanted was something much simpler conceptually: I just want to dump all my `kubeconfig` files in a single directory and then, have my tool parse them and "physically" move over the config file to `.kube/config` (or whatever is configured in my `KUBECONFIG` env var) such that it is also persistent between different shells.

## Installation

> [!TIP]
> When installing `kube-switcher`, you can use it as a `kubectl` plugin by just renaming the binary to `kubectl-switch`.
> This way, you can run it via `kubectl switch ...`

### Download Precompiled Binaries

Precompiled binaries are available for various platforms. You can download the latest release from the [GitHub Releases page](https://github.com/mirceanton/kube-switcher/releases/latest).

1. Download the appropriate binary for your system and extract the archive.
2. Make the extracted binary executable:

    ```bash
    chmod +x kube-switcher
    ```

3. Move the binary to a directory in your PATH:

    ```bash
    mv kube-switcher /usr/local/bin/kube-switcher
    ```

### Running via Docker

`kube-switcher` is also available as a Docker container:

```bash
docker pull ghcr.io/mirceanton/kube-switcher
```

### Install via homebrew

1. Add the tap

    ```bash
    brew tap mirceanton/taps
    ```

2. Install `kube-switcher`

    ```bash
    brew install kube-switcher
    ```

### Build from Source

1. Clone the repository:

    ```bash
    git clone https://github.com/mirceanton/kube-switcher
    cd kube-switcher
    ```

2. Build the tool:

    If you have [Taskfile](https://taskfile.dev/) installed, run:

    ```bash
    task build
    ```

    Otherwise, simply run the `go build` command:

    ```bash
    go build -o talswitcher
    ```

## Usage

1. Place all of your `kubeconfig` files in a single directory. I personally prefer `~/.kube/configs/`

2. Set the environment variable `KUBECONFIG_DIR`

    ```sh
    export KUBECONFIG_DIR="~/.kube/configs/"
    ```

3. Run `kube-swticher context`/`kube-switcher ctx` or `kubectl switch context`/`kubectl switch ctx` to interactively select your context from a list

    ```sh
    vscode ➜ /workspaces/kube-switcher $ kubectl switch context
    ? Choose a context:  [Use arrows to move, type to filter]
      cluster1
      cluster2
    > cluster3

    INFO[0102] Switched to context 'cluster3'

    vscode ➜ /workspaces/kube-switcher $ kubectl get nodes
    NAME       STATUS   ROLES           AGE   VERSION
    cluster3   Ready    control-plane   21h   v1.30.0
    ```

    Alternatively, pass in a context name to the command to non-interactively switch to it:

    ```sh
    vscode ➜ /workspaces/kube-switcher $ kubectl switch ctx cluster3
    INFO[0000] Switched to context 'cluster3'

    vscode ➜ /workspaces/kube-switcher $ kubectl get nodes
    NAME       STATUS   ROLES           AGE   VERSION
    cluster3   Ready    control-plane   21h   v1.30.0
    ```

4. Run the `kube-switcher namespace`/`kube-switcher ns` or `kubectl switch namespace`/`kubectl switch ns` to interactively select your current namespace from a list

    ```sh
    vscode ➜ /workspaces/kube-switcher $ kubectl switch namespace
    ? Choose a namespace:  [Use arrows to move, type to filter]
    > default
      kube-node-lease
      kube-public
      kube-system

    INFO[0012] Switched to namespace 'kube-system'

    vscode ➜ /workspaces/kube-switcher $ kubectl get pods
    NAME                               READY   STATUS    RESTARTS      AGE
    coredns-7db6d8ff4d-nqmlf           1/1     Running   1 (69m ago)   21h
    etcd-cluster1                      1/1     Running   1 (69m ago)   21h
    kube-apiserver-cluster1            1/1     Running   1 (69m ago)   21h
    kube-controller-manager-cluster1   1/1     Running   1 (69m ago)   21h
    kube-proxy-lxnvr                   1/1     Running   1 (69m ago)   21h
    kube-scheduler-cluster1            1/1     Running   1 (69m ago)   21h
    storage-provisioner                1/1     Running   3 (68m ago)   21h
    ```

    Alternatively, pass in a namespace name to the command to non-interactively switch to it:

    ```sh
    vscode ➜ /workspaces/kube-switcher $ kubectl get pods
    No resources found in default namespace.
    vscode ➜ /workspaces/kube-switcher $ kubectl switch ns kube-public
    INFO[0000] Switched to namespace 'kube-public'
    vscode ➜ /workspaces/kube-switcher $ kubectl get pods
    No resources found in kube-public namespace.
    ```

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
