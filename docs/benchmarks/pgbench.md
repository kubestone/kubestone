title: Kubestone - pgbench: Performance benchmark for PostgreSQL

# pgbench - Performance benchmark for PostgreSQL

!!! quote
    pgbench is a simple program for running benchmark tests on PostgreSQL. It runs the same sequence of SQL commands over and over, possibly in multiple concurrent database sessions, and then calculates the average transaction rate (transactions per second). By default, pgbench tests a scenario that is loosely based on TPC-B, involving five SELECT, UPDATE, and INSERT commands per transaction. However, it is easy to test other cases by writing your own transaction script files.

With the [pgbench](https://www.postgresql.org/docs/11/pgbench.html) benchmark, you can benchmark your PostgreSQL database, which can both run in the same Kubernetes cluster as kubestone, or anywhere else, as long as it's reachable.



## Mode of operation

In the pgbench CR, you need to specify the connection details to your PostgreSQL database in the `postgres` section, containing the `host`, `port`, `database`, and your `username`/`password`.

Kubestone then generates a single Kubernetes job from the CR. The pod behind the job will have an init container that runs the `pgbench -i` initialization command, and a main container that will run the actual benchmark. For these containers, you can use any options described in the [official pgbench documentation](https://www.postgresql.org/docs/11/pgbench.html) with `InitArgs` and `Args`, respectively.



## Example configuration

You can find [configuration example](https://github.com/xridge/kubestone/blob/master/config/samples/perf_v1alpha1_pgbench.yaml) in the GitHub repository.



## Sample benchmark
To run the example CR, you need the corresponding PostgreSQL database. You can create it in the same Kubernetes cluster with the following command:
```bash
kubectl create --namespace kubestone -f https://raw.githubusercontent.com/xridge/kubestone/master/tests/e2e/conf/postgres.yaml
```
Naturally, you can deploy postgres to a separate namespace, too, but then don't forget to update the `postgres.host` in the CR accordingly.

Now, you can run the sample benchmark:
```bash
kubectl create --namespace kubestone -f https://raw.githubusercontent.com/xridge/kubestone/master/config/samples/perf_v1alpha1_pgbench.yaml
```


Please refer to the [quickstart guide](../quickstart.md) for details on generic principles and setup of Kubestone.




## pgbench configuration

The complete documentation of the pgbench CR can be found in the [API Docs](../apidocs.md#perf.kubestone.xridge.io/v1alpha1.PgbenchSpec).



## Docker Image

[Docker Image for pgbench](https://hub.docker.com/r/xridge/pgbench) is provided via [xridge's pgbench-docker repository](https://github.com/xridge/pgbench-docker).



## Legal

pgbench is part of PostgreSQL, which is licensed under the [PostgreSQL License](https://opensource.org/licenses/postgresql), a liberal Open Source license, similar to the BSD or MIT licenses.
