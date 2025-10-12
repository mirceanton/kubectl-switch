#!/bin/bash
set -e

# Create config directory if it doesn't exist
export KUBECONFIG="./test/config"
export KUBECONFIG_DIR="./test/configs/"
mkdir -p $KUBECONFIG_DIR


setup() {
    echo "Performing setup..."

    # Build kubectl-switch
    echo "Building kubectl-switch..."
    go build -o kubectl-switch

    # Set up Kubernetes clusters
    echo "Setting up test clusters..."
    KUBECONFIG=$KUBECONFIG_DIR/cluster1.yaml minikube start -p test-cluster-1 &
    KUBECONFIG=$KUBECONFIG_DIR/cluster2.yaml minikube start -p test-cluster-2 &
    wait

    echo "Setup completed."
}

cleanup() {
    echo "Performing cleanup..."

    # Stop and delete minikube clusters
    echo "Stopping and deleting test-cluster-1..."
    KUBECONFIG=$KUBECONFIG_DIR/cluster1.yaml minikube stop -p test-cluster-1 || true
    KUBECONFIG=$KUBECONFIG_DIR/cluster1.yaml minikube delete -p test-cluster-1 || true

    echo "Stopping and deleting test-cluster-2..."
    KUBECONFIG=$KUBECONFIG_DIR/cluster2.yaml minikube stop -p test-cluster-2 || true
    KUBECONFIG=$KUBECONFIG_DIR/cluster2.yaml minikube delete -p test-cluster-2 || true

    # Remove config directory
    echo "Removing test directory..."
    rm -rf tests/
    rm kubectl-switch

    echo "Cleanup completed."
}


run_tests() {
    echo "Running tests..."

    # Test cluster switching and namespace operations
    echo "Switching context to test-cluster-1..."
    ./kubectl-switch context test-cluster-1

    echo "Validating cluster switch to test-cluster-1..."
    kubectl get nodes | grep "test-cluster-1" || {
        echo "Error: test-cluster-1 not found in node list!" >&2
        exit 1
    }

    echo "Checking no pods in default namespace..."
    kubectl get pods --namespace=default 2>&1 | grep "No resources found" || {
        echo "Error: Pods found in default namespace!" >&2
        exit 1
    }

    echo "Switching to kube-system namespace..."
    ./kubectl-switch namespace kube-system

    echo "Validating kube-system namespace selection and kube-apiserver running..."
    kubectl get pods --namespace=kube-system | grep "kube-apiserver" || {
        echo "Error: kube-apiserver not found in kube-system!" >&2
        exit 1
    }

    echo "Switching back to default namespace..."
    ./kubectl-switch namespace default

    echo "Checking again that no pods are in default namespace..."
    kubectl get pods --namespace=default 2>&1 | grep "No resources found" || {
        echo "Error: Pods found in default namespace!" >&2
        exit 1
    }

    echo "Switching context to test-cluster-2..."
    ./kubectl-switch context test-cluster-2

    echo "Validating cluster switch to test-cluster-2..."
    kubectl get nodes | grep "test-cluster-2" || {
        echo "Error: test-cluster-2 not found in node list!" >&2
        exit 1
    }

    echo "Switching to previous context..."
    ./kubectl-switch -

    echo "Validating cluster switch to test-cluster-1..."
    kubectl get nodes | grep "test-cluster-1" || {
        echo "Error: test-cluster-1 not found in node list!" >&2
        exit 1
    }
    echo "========================================================================================="
    echo "Tests completed successfully!"
    echo "========================================================================================="
}

# Set up trap to ensure cleanup happens on exit
trap cleanup EXIT

# Main execution
echo "Starting E2E tests..."
setup
run_tests
echo "E2E tests completed successfully!"
