# event-driven-go

## ELK
### Elasticsearch
Check status elasticsearch

`
curl http://localhost:9200 -u elastic:MyPw123"=
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

# Reference

- [zipkin-kafka-collector](https://github.com/openzipkin/zipkin/blob/master/zipkin-server/README.md#kafka-collector)

- [go-zipkin](https://medium.com/oracledevs/setup-a-distributed-tracing-infrastructure-with-zipkin-kafka-and-cassandra-d0a68fb3eee6)