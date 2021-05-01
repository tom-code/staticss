package main

import (
  "bufio"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "os"
  "strings"
)

type AllocationConfig struct {
  Namespace string  `json:"namespace"`
  Pod string        `json:"pod"`
  Address string    `json:"address"`
  Gateway string    `json:"gateway"`
}
type IpamConfig struct {
  Allocations []AllocationConfig `json:"allocations"`
}
type Config struct {
  Ipam IpamConfig `json:"ipam"`
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
  for _, c := range config.Ipam.Allocations {
    if (c.Namespace == namespace) && (c.Pod == podName) {
      addr = c.Address
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

  outf := `
  {
    "ips": [
        {
          "address": "%s"
        }
    ]
  }
  `
  fmt.Printf(outf, addr)
}
