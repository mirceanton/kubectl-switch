#!/usr/bin/env python3

import os
import shutil
import sys
import yaml

CONFIGS_DIR = os.path.expanduser(
    os.getenv("KSWITCHER_CONFIGS_DIR", "~/.kube/configs")
)
CONFIG_FILE = os.path.expanduser("~/.kube/config")

def list_contexts():
    contexts = {}

    for filename in os.listdir(CONFIGS_DIR):
        file_path = os.path.join(CONFIGS_DIR, filename)
        if os.path.isfile(file_path):
            with open(file_path, "r") as f:
                config = yaml.safe_load(f)
                if "contexts" in config and len(config["contexts"]) > 0:
                    context_name = config["contexts"][0]["name"]
                    contexts[context_name] = file_path

    return contexts

def choose_context(contexts):
    # Show all of the available contexts
    print("Available contexts:")
    for i, context in enumerate(contexts.keys()):
        print(f"{i+1}. {context}")

    # Prompt the user for a choice
    while True:
        choice = input("Choose a context (enter the number): ")
        if choice.isdigit() and 1 <= int(choice) <= len(contexts.keys()):
            return list(contexts.keys())[int(choice) - 1]


def use_context(context, file):
    shutil.copy(file, CONFIG_FILE)
    print(f"Config file for context '{context}' has been moved to '{CONFIG_FILE}'.")

if __name__ == "__main__":
    context_list = list_contexts()
    if not context_list:
        print("No contexts found in the 'configs' directory.")
        exit(1)

    if len(sys.argv) > 1:
        context = sys.argv[1]
    else:
        context = choose_context(context_list)

    use_context(context, context_list[context])
