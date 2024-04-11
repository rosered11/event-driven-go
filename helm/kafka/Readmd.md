# Kafka

## Manual

1. Create k8s by helm

`helm install my-kafka oci://registry-1.docker.io/bitnamicharts/kafka ./mychart`

- Install with values.yaml

`helm install -f kafka/values.yaml my-kafka bitnami/kafka --version 28.0.3`

- Generate yaml
`helm template -f kafka/values.yaml kafka-broker bitnami/kafka --version 28.0.3 -n {namespace} > temp.yaml`

- Remove
`helm uninstall my-kafka`