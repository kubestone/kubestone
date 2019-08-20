# Disk

## Fio
With the [fio](https://fio.readthedocs.io/en/latest/fio_doc.html) benchmark you
can measure the I/O performance of the disks used in your Kubernetes cluster.

Kubestone generates a Kubernetes Job from each fio CR that will run a single pod
with the defined fio job.

<!-- TODO: Add usage, examples, CR description -->
