# kubectl-switch

A simple tool to manage and switch between multiple Kubernetes configuration files. This tool allows you to dump all your kubeconfig files in a single directory and easily switch between them as well as configure the current namespace.

## Features

- Manage multiple kubeconfig files in a single directory.
- Easily switch between different Kubernetes contexts.
- Easily switch between different namespaces.
- Persistent context/namespace switching across different shells.
- Interactive prompts support fuzzy search.
- Non-destructive operations on your original kubeconfig files. `kubectl-switch` will never edit those. It will only work with copies of them.

## Why `kubectl-switch`?

`kubectl-switch` is an alternative to tools like `kubectx`, `kubens` and `kubie`. It has been created because I feel all of those fall short in certain regards:

- `kubectx` and `kubens` assume all your contexts are defined in a single config file. Yes, there is some hackery you can do to your `KUBECONFIG` environment variable to make it work with multiple files, but it is (in my opinion) a sub-par workflow and I never really liked it
- `kubie` spawns a new shell when you use it to change contexts, making it practically impossible to integrate into scripts or taskfile automation. Also I consider it to be too complicated of a solution for what is fundamentally a really simple problem

What I wanted was something very simple conceptually: I just want to dump all my `kubeconfig` files in a single directory and then have my tool parse them and "physically" move over the config file to `.kube/config` (or whatever is configured in my `KUBECONFIG` env var) such that it is also persistent between different shells. Here is where `kubectl-switch` comes in!

## Installation

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

### Running via Docker

`kubectl-switch` is also available as a Docker container:

```bash
docker pull ghcr.io/mirceanton/kubectl-switch
```

### Install via homebrew

1. Add the tap

    ```bash
    brew tap mirceanton/taps
    ```

2. Install `kubectl-switch`

    ```bash
    brew install kubectl-switch
    ```

### Build from Source

1. Clone the repository:

    ```bash
    git clone https://github.com/mirceanton/kubectl-switch
    cd kubectl-switch
    ```

2. Build the tool:

    If you have [Taskfile](https://taskfile.dev/) installed, run:

    ```bash
    task build
    ```

    Otherwise, simply run the `go build` command:

    ```bash
    go build -o kubectl-switch
    ```

## Usage

1. Place all of your `kubeconfig` files in a single directory. I personally prefer `~/.kube/configs/` but you can do whatever you fancy.

2. Set the environment variable `KUBECONFIG_DIR`

    ```sh
    export KUBECONFIG_DIR="~/.kube/configs/" # <- put here whatever folder you decided on at step 1
    ```

3. Run `kubectl switch ctx` or `kubectl switch context` to interactively select your context from a list

    ```sh
    vscode ➜ /workspaces/kubectl-switch $ kubectl switch context
    ? Choose a context:  [Use arrows to move, type to filter]
      cluster1kubectl-switch
      cluster2
    > cluster3kubectl-switch

    INFO[0102] Switched tkubectl-switchuster3'

    vscode ➜ /workspaces/kubectl-switch $ kubectl get nodes
    NAME       STATUS   ROLES           AGE   VERSION
    cluster3   Ready    control-plane   21h   v1.30.0
    ```

    Alternatively, pass in a context name to the command to non-interactively switch to it:

    ```sh
    vscode ➜ /workspaces/kubectl-switch $ kubectl switch ctx cluster3
    INFO[0000] Switched to context 'cluster3'

    vscode ➜ /workspaces/kubectl-switch $ kubectl get nodes
    NAME       STATUS   ROLES           AGE   VERSION
    cluster3   Ready    control-plane   21h   v1.30.0
    ```

    This will literally detect the file in which that context is defined and copy-paste it to the path defined in your `KUBECONFIG`.

4. Run the `kubectl switch ns` or `kubectl switch namespace` to interactively select your current namespace from a list

    ```sh
    vscode ➜ /workspaces/kubectl-switch $ kubectl switch namespace
    ? Choose a namespace:  [Use arrows to move, type to filter]
    > default
      kube-node-lease
      kube-public
      kube-system

    INFO[0012] Switched to namespace 'kube-system'

    vscode ➜ /workspaces/kubectl-switch $ kubectl get pods
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
    vscode ➜ /workspaces/kubectl-switch $ kubectl get pods
    No resources found in default namespace.
    vscode ➜ /workspaces/kubectl-switch $ kubectl switch ns kube-public
    INFO[0000] Switched to namespace 'kube-public'
    vscode ➜ /workspaces/kubectl-switch $ kubectl get pods
    No resources found in kube-public namespace.
    ```

    This will modify your currently active kubeconfig file (not the original one from `KUBECONFIGS_DIR`, but rather the copy from `KUBECONFIG`) to set the current namespace.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
