Welcome to kubestone's documentation!
=====================================

Kubestone is a `Kubernetes Operator`_ that can run benchmarks on your Kubernetes
cluster. Following the Operator Pattern, benchmarks are defined as Custom
Resources.

Benchmarks
----------
Kubestone currently supports the following benchmarks:

* Disk
   - fio
* Network
   - iperf3

.. _Kubernetes Operator: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/

.. toctree::
   :maxdepth: 2
   :caption: Contents:

   install
   quickstart
   benchmarks/index

Indices and tables
==================

* :ref:`search`
