[![Go](https://github.com/auntan/network_exporter/actions/workflows/go.yml/badge.svg)](https://github.com/auntan/network_exporter/actions/workflows/go.yml)

# Network exporter
Exporter that provides network conditions as metrics for prometheus

## Ping
Measure round trip time and packet loss using ICMP pings

### Describing deployments
Deployment is something that is deployed on one or many hosts
```yaml
deployments:
  # deployment with id:'database' is deployed on 2 nodes
  database:
    hosts:
      node1:
        address: 192.168.0.1
      node2:
        address: 192.168.0.2
  # 'svc1' on single node
  svc1:
    hosts:
      node3:
        address: 192.168.0.3
  # 'svc2' on single node too
  svc2:
    hosts:
      node4:
        address: 192.168.0.4
```

### Describing probes
Probe is measuring network conditions from source to target deployments
```yaml
probes:
  # svc1 has 1 dependency on database,
  # so we create the probe with id:'svc1' that would ping from svc1 to database deployments hosts
  # finally it resolves to following pings 
  # 192.168.0.3 -> 192.168.0.1 and 192.168.0.2
  svc1: 
    source: svc1 
    targets: [database] 

  # svc 2 has dependencies on database and svc1,
  # so pings would be:
  # 192.168.0.4 -> 192.168.0.1, 192.168.0.2 and 192.168.0.3
  svc2:
    source: svc2
    targets: [database, svc1]
```

## Provided metrics
`rtt_seconds`: Round trip time from source to target

`sent_packets_total`: Total sent packets from source to target

`recv_packets_total`: Total received packets

### Calculating packet loss
Packet loss percent can be calculated as `((sent_packets_total - recv_packets_total) / sent_packets_total) * 100`
