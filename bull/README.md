use `eval $(minikube docker-env)` to use docker inside minikube
use `docker build -t server:1.0.0 .` to build new server docker image. Remember increase version for each critical change.

`k apply -f k8s/server-depl.yaml`

localhost:9090 prometheus `k -n monitoring port-forwa
rd svc/prometheus-operated 9090`
localhost:3000 grafana admin/devops123 `k -n grafana port-forward svc/grafana 3000`

