# Dapani

Dapani is a tool that allows you to analyze the costliest workload links in your cluster. It relies on Kubernetes/Istio and Prometheus to gather
data, and uses publicly-available cloud egress rates to estimate the overall egress costs of your services.

## Usage

To use this on your kubernetes cluster, make sure you have a kubeconfig in your home directory, and make sure Istio is installed on your cluster, with the prometheus addon enabled.

### Creating `destination_pod`

First, you must create the `destination_pod` metric for Dapani to read from.

Add the following to all of your deployments:

```yaml
spec:
  template:
    metadata:
      annotations:
        sidecar.istio.io/extraStatTags: destination_pod
```

Add the following to your Istio Operator:

```yaml
spec:
  values:
    telemetry:
      v2:
        prometheus:
          configOverride:
            inboundSidecar:
              metrics:
                - name: request_bytes
                  dimensions:
                    destination_pod: upstream_peer.name
            outboundSidecar:
              metrics:
                - name: request_bytes
                  dimensions:
                    destination_pod: upstream_peer.name
            gateway:
              metrics:
                - name: request_bytes
                  dimensions:
                    destination_pod: upstream_peer.name
```


### Running

To Build `dapani`:

```
go install
```

Run:

```
dapani analyze
```

This assumes your cluster is on GCP. To change this to the two options of AWS and Azure, run as follows:
```
dapani analyze --cloud aws
```
To point dapani to your own pricing sheet, run as follows:
```
dapani analyze --pricePath <path to .json>
```
To only use data from a specific time range, run as follows:
```
dapani analyze --queryBefore 10h
```
This will only use call data from 10 hours ago and previous.
