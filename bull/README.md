use `eval $(minikube docker-env)` to use docker inside minikube
use `docker build -t server:1.0.0 .` to build new server docker image. Remember increase version for each critical change.

`k apply -f k8s/server-depl.yaml`