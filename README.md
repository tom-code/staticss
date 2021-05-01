# staticss
k8s cni/ipam plugin with configurable static ips for statefulsets

This cni plugin allows to specify static ip address for namespace / pod-name combination.
This is useful of statefulsets which has pod name predictable. It can be used also for other custom created pods with predictable names.

example how it can be configured for additional network interface in NetworkAttachementDefinition specification:

```
apiVersion: "k8s.cni.cncf.io/v1"
kind: NetworkAttachmentDefinition
metadata:
  name: example
spec:
  config: '{
        "cniVersion": "0.3.1",
                "name": "bridge",
                "type": "bridge",
                "bridge": "b1",
                "ipMasq": false,
                "ipam": {
                        "type": "staticss",
                        "allocations": [
                                          {"namespace": "example", "pod": "ex-0", "address": "10.0.2.40/24"},
                                          {"namespace": "example", "pod": "ex-1", "address": "10.0.2.41/24"},
                                          {"namespace": "example", "pod": "ex-2", "address": "10.0.2.42/24"},
                                          {"namespace": "example", "pod": "ex-3", "address": "10.0.2.43/24"},
                                          {"namespace": "example", "pod": "ex-4", "address": "10.0.2.44/24"}
                                       ]
                }
        }'
``` 
