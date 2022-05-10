# staticss
k8s cni/ipam plugin with configurable static ips for statefulsets

This cni plugin allows to specify static ip address for namespace / pod-name combination.
This is useful for statefulsets which has predictable pod name. It can be used also for other custom created pods with predictable names.

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

## compile
To compile you need go compiler.
```
go build github.com/tom-code/staticss
```

## install
To install, copy staticss binary to /opt/cni/bin of all kubernetes hosts which can run related pods.
Don't forget to make sure staticss binary is executable.


## routes
plugin can also install additional routes. see example which adds route to 10.0.9.0/24 via gateway 10.0.0.99:
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
                        "routes": [{"dst": "10.0.9.0/24", "gw": "10.0.0.99"}],
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


## multiple ip address
it is possible to assign multimple addresses to one pod:
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
                        "routes": [{"dst": "10.0.9.0/24", "gw": "10.0.0.99"}],
                        "allocations": [
                                          {"namespace": "example", "pod": "ex-0", "address": "10.0.2.40/24"},
                                          {"namespace": "example", "pod": "ex-0", "address": "fec0:aaaa:1::2/64"}
                                       ]
                }
        }'
```
