kubectl create secret generic prometheus-secret-file --from-file=./vol/prometheus/prometheus.yml
kubectl create configmap prometheus-config --from-file=./vol/prometheus/prometheus.yml
