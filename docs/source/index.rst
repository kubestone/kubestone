.. image:: images/kubestone-logo.png

Kubestone is a `Kubernetes Operator`_ that can run benchmarks on your Kubernetes
cluster. Following the Operator Pattern, benchmarks are defined as Custom
Resources.

Benchmarks
----------
Kubestone currently supports the following benchmarks:

* Disk
   - `fio`_
* Network
   - `iperf3`_

.. _Kubernetes Operator: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
.. _fio: benchmarks/disk.html#fio
.. _iperf3: benchmarks/network.html#iperf3

.. toctree::
   :maxdepth: 2
   :caption: Contents:

   install
   quickstart
   benchmarks/index

Indices and tables
==================

* :ref:`search`
