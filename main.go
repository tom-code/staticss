package main

import (
  "bufio"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "strings"
)

type AllocationConfig struct {
  Namespace string  `json:"namespace"`
  Pod string        `json:"pod"`
  Address string    `json:"address"`
  Gateway string    `json:"gateway"`
}

type RouteConfig struct {
  Dst string `json:"dst,omitempty"`
  Gateway string `json:"gw,omitempty"`
}

type IpamConfig struct {
  Allocations []AllocationConfig `json:"allocations"`
  Gateway string `json:"gateway,omitempty"`
  Routes []RouteConfig `json:"routes,omitempty"`
}
type Config struct {
  Ipam IpamConfig `json:"ipam"`
}
type CNIIpConfig struct {
  Address string `json:"address"`
  Gateway string ` json:"gateway,omitempty"`
}
type CNIIpam struct {
  CniVersion string `json:"cniVersion"`
  Ips []CNIIpConfig `json:"ips"`
  Routes []RouteConfig `json:"routes,omitempty"`
}

func main() {
  command := os.Getenv("CNI_COMMAND")
  if command != "ADD" {
    fmt.Println(`{"cniVersion": "1.0.0"}`)
    return
  }

  args := os.Getenv("CNI_ARGS")
  argss := strings.Split(args, ";")
  var podName string
  var namespace string
  for _, a := range argss {
    as := strings.Split(a, "=")
    if len(as) != 2 {
      continue
    }
    if as[0] == "K8S_POD_NAMESPACE" {
      namespace = as[1]
    }
    if as[0] == "K8S_POD_NAME" {
      podName = as[1]
    }
  }

  si := bufio.NewReader(os.Stdin)
  data, err := ioutil.ReadAll(si)

  var config Config
  err = json.Unmarshal(data, &config)
  if err != nil {
    out := `
    {
      "cniVersion": "1.0.0",
      "code": 7,
      "msg": "Invalid Configuration",
      "details": "%s"
    }
    `
    fmt.Printf(out, err.Error())
    os.Exit(1)
  }

  addr := ""
  gateway := config.Ipam.Gateway
  routes := config.Ipam.Routes
  for _, c := range config.Ipam.Allocations {
    if (c.Namespace == namespace) && (c.Pod == podName) {
      addr = c.Address
      if len(c.Gateway) > 0 {
        gateway = c.Gateway
      }
    }
  }

  if len(addr) == 0 {
    out := `
    {
      "cniVersion": "1.0.0",
      "code": 7,
      "msg": "Invalid Configuration",
      "details": "no ip configured for %s/%s"
    }
    `
    fmt.Printf(out, namespace, podName)
    os.Exit(1)
  }

  outStruct := CNIIpam {
    CniVersion: "1.0.0",
    Ips: []CNIIpConfig{
      {
        Address: addr,
        Gateway: gateway,
      },
    },
    Routes: routes,
  }

  out, err := json.Marshal(&outStruct)
  if err != nil {
    log.Println(err)
  }
  fmt.Println(string(out))
  log.Printf(string(out))
}
