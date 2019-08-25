# Fio - Flexible I/O tester

!!! quote
    fio is a tool that will spawn a number of threads or processes doing a particular type of I/O action as specified by the user.  The typical use of  fio  is  to  write  a  job  file matching the I/O load one wants to simulate. 

With the [fio](https://fio.readthedocs.io/en/latest/fio_doc.html) benchmark you can measure the I/O performance of the disks used in your Kubernetes cluster. Kubestone generates a Kubernetes Job from each fio CR that will run a single pod with the defined fio job. 






