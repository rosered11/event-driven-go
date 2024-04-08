# event-driven-go

## ELK
### Elasticsearch
Check status elasticsearch

`
curl http://localhost:9200 -u elastic:MyPw123
`

Genrate token for Kibana

`curl -X POST "http://localhost:9200/_security/service/elastic/kibana/credential/token/token1?pretty" -u elastic:MyPw123
`

## Kibana

Set token from generate token on

```.env
ELASTICSEARCH_SERVICEACCOUNTTOKEN={service-accout-token-elastic}
```

## Logstash

Need to create user in Elasticsearch for Logstash for create data to indict in Elasticsearch

## Zipkin

Schema body data
```json
[
    {
        "traceId": "7ed77808cef3db4e",
        "id": "7ed77808cef3db4e",
        "kind": "CLIENT",
        "name": "http/post",
        "timestamp": 1696765593499252,
        "duration": 430112,
        "localEndpoint": {
            "serviceName": "my_service",
            "port": 8081
        },
        "annotations": [
            {
                "timestamp": 1696765593499319,
                "value": "Connecting"
            },
            {
                "timestamp": 1696765593499498,
                "value": "DNS Start"
            },
            {
                "timestamp": 1696765593501685,
                "value": "DNS Done"
            },
            {
                "timestamp": 1696765593501741,
                "value": "Connect Start"
            },
            {
                "timestamp": 1696765593685768,
                "value": "Connect Done"
            },
            {
                "timestamp": 1696765593685904,
                "value": "Connected"
            },
            {
                "timestamp": 1696765593686118,
                "value": "Wrote Headers"
            },
            {
                "timestamp": 1696765593686139,
                "value": "Wrote Request"
            },
            {
                "timestamp": 1696765593929133,
                "value": "First Response Byte"
            }
        ],
        "tags": {
            "http.method": "POST",
            "http.path": "",
            "http.response.size": "1256",
            "httptrace.connect_done.addr": "[2606:2800:220:1:248:1893:25c8:1946]:80",
            "httptrace.connect_done.network": "tcp",
            "httptrace.connect_start.addr": "[2606:2800:220:1:248:1893:25c8:1946]:80",
            "httptrace.connect_start.network": "tcp",
            "httptrace.dns_done.addrs": "2606:2800:220:1:248:1893:25c8:1946 , 93.184.216.34",
            "httptrace.dns_start.host": "example.com",
            "httptrace.get_connection.host_port": "example.com:80",
            "httptrace.got_connection.reused": "false",
            "httptrace.got_connection.was_idle": "false"
        }
    }
]
```

## Kafka

Generate config in linux

```
echo -e "security.protocol=SASL_PLAINTEXT\nsasl.mechanism=PLAIN\nsasl.jaas.config=org.apache.kafka.common.security.scram.ScramLoginModule required username=\"user1\" password=\"GaferjmEV1%\";" > kafka-config.properties
```

### Create topic

```
kafka-topics.sh --bootstrap-server localhost:9092 --topic first_topic --create --partitions 3 --replication-factor 1 --command-config config.properties
```

### Consumer

consume cli
`kafka-console-consumer.sh --topic zipkin --bootstrap-server localhost:9092 --group zipkin_debug --consumer.config {}`

consumer.config
`

`

### Producer

producer cli
`kafka-console-producer.sh --broker-list 192.168.65.3:31363,192.168.65.3:30307,192.168.65.3:30513 --topic test`

## Grafana

### Dashboard

- [D1](https://grafana.com/grafana/dashboards/18276-kafka-dashboard/)
- [D2](https://grafana.com/grafana/dashboards/721-kafka/)

## K8S

Convert file to secret
`kubectl create secret generic my-secret --from-file=mysecretfile.txt=mysecretfile.txt`

Convert file to configmap
`kubectl create configmap my-configmap --from-file=myconfig.yml`

### Context View
kubectl config view

### Switch Context
kubectl config use-context docker-desktop

### Setup Access Cluster in Azure
Warning or Error:
The azure auth provider is already removed as of today which was earlier deprecated. kubelogin is now the default way. So you might get:

WARNING: the azure auth plugin is deprecated

error: The azure auth plugin has been removed

Fix:
Below steps fixed the issue for me on windows:

Remove all the config folders inside user folder viz., %USERPROFILE% which will be .kube, .azure-kubelogin.

Then download and install the latest releases of kubectl and kubelogin (which is essentially just unzipping to a folder and adding that into user's Path environment variable.)

Perform:

```
az login
az account set --subscription <subscription id>
az aks get-credentials --resource-group <resource group name> --name <AKS cluster name>
kubelogin convert-kubeconfig -l azurecli
Now try: kubectl get pods -A
```

### Command Forward port
kubectl port-forward svc/my-service 6399:6379 -n default

### Docker use Context Azure
docker context create aci my-context

### Describe

Check topic info
`kafka-topics.sh --bootstrap-server localhost:9092 --describe --topic t1`

## Resource

### Kafka

```
- brokers 3
- ram 2 gi
- cpu 500m
```

### Prometheus

Prometheus stores an average of only 1-2 bytes per sample. Thus, to plan the capacity of a Prometheus server, you can use the rough formula:

`needed_disk_space = retention_time_seconds * ingested_samples_per_second * bytes_per_sample`

- keep data 7 day
  7 day = 604800 seconds
- kafka-broker(sample 3)
- scrape_interval 5 seconds

disk = 604800 * (5 * 3) * 2 = 18144000 bytes = 18.144 MB

```
- brokers 1
- ram 512 Mi
- cpu 500m
```
### Zipkin in Cassandra

How to set TTL for logging in Cassandra
```
ALTER TABLE zipkin2.span
WITH  default_time_to_live= 604800 # unit is seconds
```


Note Dev:
Storage 4 Gi
- kafka storage - 2
- prometheus - 1
- Grafana - 1
- Logstash - 1
- Elasticsearch - 1
- Kibana - 1
- Zipkin - 1

# Reference

- [zipkin-kafka-collector](https://github.com/openzipkin/zipkin/blob/master/zipkin-server/README.md#kafka-collector)

- [go-zipkin](https://medium.com/oracledevs/setup-a-distributed-tracing-infrastructure-with-zipkin-kafka-and-cassandra-d0a68fb3eee6)

- [jmx-exporter](https://github.com/prometheus/jmx_exporter/tree/main)