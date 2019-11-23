title: Kubestone - Development guide

# Development guide



## Overview

Kubestone is implemented as a [Kubernetes Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) in [Go language](https://golang.org) with the help of [Kubebuilder](https://kubebuilder.io). 

As with all operators, the manager deployment ties together multiple controllers, whose responsibility is to reach the desired state set by Custom Resources. In case of Kubestone the **controllers are implementing the benchmark execution logic**. Therefore, we have one controller per each supported benchmark. 



!!! hint
    If you are not familiar with Kubebuilder, it is advised to read through it's [super awesome documentation](https://book.kubebuilder.io). It contains very detailed information on how the reconcile loop should be created. This information will be essential if you would like to add support for a new benchmark in Kubestone.



## Dev box setup (OS X & Linux)

### Compiler tools

The following software are required to build and run Kubestone:

- [GoLang](https://golang.org/dl/) v1.11 or greater 
- make
- [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/)
- [Kustomize v3.1.0](https://kustomize.io/)

Kubestone uses the Module feature of Go 1.11. You need to make sure that you enable it for your environment:
```bash
$ export GO111MODULE=on
```



### Kubernetes setup

If you already have access to Kubernetes (i.e. you can use `kubectl` to interact with the cluster) then you are good to go. 

If not, you can use [KinD (Kubernetes in Docker)](https://github.com/kubernetes-sigs/kind), [MiniKube](https://github.com/kubernetes/minikube) or [K3S](https://k3s.io/) to run a local 'cluster' in your machine. Please follow the installation steps of your chosen distribution and make sure you have a working `kubectl` before you proceed to the next step.

!!! note
    Our CI system uses KinD (Kubernetes in Docker) for the end-to-end test execution. We found it reliable and fast enough to support our development. If you are planning to test locally we recommend that you install KinD.




### Build

#### Git Repo Clone

You need to clone the repository to your local machine. This can be done using [xridge's repo](https://github.com/xridge/kubestone) or your own fork. If you would like to send a PR it is required that you fork the repository.

```bash
$ git clone https://github.com/xridge/kubestone
```



#### Local build

Now you obtained the sources it is time to build the bits:

```bash
$ cd kubestone
$ make manager
```

The `make` command will download all the necessary dependencies, generate all the required Kubernetes Objects (CRDs, RBAC objects) and build the go files. End result will be a manager binary located at `bin/manager`. 



#### Installing the CRDs

Before we start the local build we need to make sure that the Custom Resource Definitions are created in the cluster. The following command generates the CRDs from the types and loads them to the server.

```bash
$ make install
```

!!! note
    You need to re-run this command every time the API is changed.



#### Starting the service locally

We have now everything in place: the application is built and the corresponding Custom Resource Definitions are installed to the cluster. 

Time to start the manager:

```bash
$ make run
```

If you load any Custom Resource into the cluster you will notice the reconcile loop executing.



#### VS Code support 

Some of us are using [VS Code](https://code.visualstudio.com) for development and debugging. The necessary configuration files are stored in the repository, so you can open the directory and start execution by pressing F5. 



## Adding a new benchmark

If you would like to add a new benchmark, this guide is for you. Please follow on!

### Creating the API

Every new benchmark should have it's own API, which can be created using `kubebuilder`:

```bash
$ kubebuilder create api --group perf --version v1alpha1 --kind MyBenchmark 
```



The command above creates the scaffold for the benchmark's data structures, controller and CRD:

```bash
$ git status
On branch master
Your branch is up to date with 'origin/master'.

Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git checkout -- <file>..." to discard changes in working directory)

	modified:   PROJECT
	modified:   api/v1alpha1/doc.go
	modified:   api/v1alpha1/zz_generated.deepcopy.go
	modified:   config/crd/kustomization.yaml
	modified:   controllers/suite_test.go
	modified:   main.go

Untracked files:
  (use "git add <file>..." to include in what will be committed)

	api/v1alpha1/mybenchmark_types.go
	config/crd/patches/cainjection_in_mybenchmarks.yaml
	config/crd/patches/webhook_in_mybenchmarks.yaml
	config/samples/perf_v1alpha1_mybenchmark.yaml
	controllers/mybenchmark_controller.go

no changes added to commit (use "git add" and/or "git commit -a")
```



### Adding benchmark logic

Every benchmark needs a set of input parameters. In Kubestone, the input parameters are passed via Custom Resources. The Custom Resources are generated from the API package. In case of `MyBenchmark`, the associated file is `api/v1alpha1/mybenchmark_types.go`.



The benchmark logic should be implemented in the reconcile loop, located under `controllers/mybenchmark_controller.go`. For information on how the reconcile loop should be implemented please refer to Kubebuilder's documentation or take a look in one of the already implemented benchmarks for guidance.



### Testing the benchmark

Once you are ready with the implementation of the benchmark you can use the created CRs (`config/samples/perf_v1alpha1_mybenchmark.yaml`) to trigger execution.

The new parameters introduced in `api/v1alpha1/mybenchmark_types.go` should also be reflected in the CR. 



## Tests

Kubestone has both unit and end-to-end tests. 

### Unit tests

Used to validate the basic logic of the functions. Majority of the cases that means every logic which is used to transform a CR to Kubernetes Object should be tested. Unit tests can be triggered with the following command:

```bash
$ make test
```

Unit tests are located in `controller` package.



### End to end tests

Used to validate that the operator deployment and the provided test examples are working. During end-to-end tests the following steps are executed

1. Kubernetes in Docker, kubectl, kustomize is installed to a temp directory.
2. Local, 3 node Kubernetes deployment is started.
3. Docker image is built from the Operator.
4. The provided system tests are executed. 
   This step mostly involves executing at least one from the provided examples.
   Validation (pods are succeeded, objects are created) occurs in this phase.

End to end tests require more CPU and Memory for execution. They can be triggered with:

```bash
$ make e2e-test
```

End to end tests are located in `tests/e2e/` package.



### Linters

During code build and unit test execution a set of linters are executed to make sure that the codebase meets formal standards. If you would like to execute the linting steps the following command can be used:

```bash
$ make fmt vet lint
```



## Contributions

Every contribution is appreciated! 
If you happen to find a bug or added a new benchmark or improved the documentation please post a Pull Request for the project. We will do our best to review it in a timely manner.



If you are proposing a new benchmark, please make sure that:

1. It meets the formal standards (`make fmt vet lint` should pass)
2. Have a reasonable amount of tests:
   1. Unit tests where it makes sense
   2. System tests for at least one of the provided examples
3. Have enough documentation for the users to get started



If you have any questions, got stuck please reach out to us in [Slack](https://kubestone.slack.com).