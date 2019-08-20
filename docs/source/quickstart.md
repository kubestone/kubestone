# Quickstart

### Run a benchmark

Create a dedicated namespace for benchmarking

```bash
$ kubectl create namespace kubestone
namespace/kubestone created
```

Start sample fio benchmark by creating CR

```bash
$ kustomize build config/samples/fio/overlays/builtin_jobs | kubectl create -n kubestone -f -
fio.perf.kubestone.xridge.io/fio-sample created
```

Additional benchmarks are located in the config/samples/ directory.
