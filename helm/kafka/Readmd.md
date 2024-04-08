# Kafka

## Manual

1. Create k8s by helm

`helm install my-kafka oci://registry-1.docker.io/bitnamicharts/kafka ./mychart`

- Install with values.yaml

`helm install -f kafka/values.yaml my-kafka bitnami/kafka --version 28.0.3`

- Remove
`helm uninstall my-kafka`