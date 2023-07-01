## Kubernetes Job Queue
A test project to try different solutions to build a job queue on top of Kubernetes.


### Setup
Create a cluster
```bash
go install sigs.k8s.io/kind@v0.20.0 && kind create cluster
```
Build the image
```bash
docker build . -t worker:1
```
Directly running the image will start an API server. Additional flags can be passed to start a worker node (`WORKER`) and to exit if the queue is empty (`EXIT`).  
  
Load the image into the cluster
```bash
kind load docker-image worker:1
```

Deploy the definitions from `deployments/` folder. Each file is a different approach.