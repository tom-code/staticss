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
  f, err := os.OpenFile("/tmp/d1.txt", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0666)
  if err != nil {
    os.Exit(1)
  }
  defer f.Close()
  f.WriteString("log1\n")
  f.WriteString(os.Getenv("CNI_COMMAND"))
  f.WriteString("\n")
  f.WriteString(os.Getenv("CNI_CONTAINERID"))
  f.WriteString("\n")
  f.WriteString(os.Getenv("CNI_ARGS"))
  f.WriteString("\n")

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
  f.Write(data)
  f.WriteString("\n")

  var config Config
  err = json.Unmarshal(data, &config)
  if err != nil {
    f.WriteString(err.Error())
  }

  addr := "2.2.2.3/16"
  for _, c := range config.Ipam.Allocations {
    if (c.Namespace == namespace) && (c.Pod == podName) {
      addr = c.Address
    }
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
