# kswitcher

The most basic implementation of a `kubeconfig` file manager, written in python.

## Prerequisites

- Python 3

## Installation

### Convenience Script

Download the installation script to a file, inspect it, and then run it:

```bash
curl -fsSL https://raw.githubusercontent.com/mirceanton/kswitcher/v1.0.0/scripts/install.sh -o get-kswitcher.sh
sudo sh ./get-kswitcher.sh
```

### Manual Installation

1. Clone the repository:

    ```shell
    git clone https://github.com/mirceanton/kswitcher.git
    ```

2. Change into the cloned directory:

    ```shell
    cd kswitcher
    ```

3. Install the dependencies:

    ```shell
    pip install -r requirements.txt
    ```

4. Add the python script to `PATH`

    ```shell
    chmod +x kswitcher.py
    sudo mv kswitcher.py /usr/local/bin/kswitcher
    ```

## Usage

To switch between Kubernetes contexts, follow these steps:

1. Configure your Kubernetes context files:

    - By default, the script assumes the context files are located in `~/.kube/configs` directory. You can specify a different directory by setting the environment variable `KSWITCHER_CONFIGS_DIR` to the desired directory path.
    - Each context file should be a valid Kubernetes YAML configuration file, containing a single context definition
    - The script reads the contexts section of each configuration file to determine the available contexts.

2. Run the script:

    ```shell
    python kswitcher.py [context_name]
    ```

    If you provide the `context_name` argument, the script will switch to that context if it exists.

    If no argument is provided, the script will prompt you to choose a context from the available options.

## License

MIT

## Author Information

A script developed by [Mircea-Pavel ANTON](https://www.mirceanton.com).
