# consistenthash

A simple, extensible **consistent hashing** library for Go.

* 64-bit hash space (Murmur3 by default)
* Virtual nodes (replicas) for even distribution
* Pluggable hash functions via `Hasher` interface

## Install

```bash
go get github.com/prfc0/consistenthash
```

## Usage

```go
package main

import (
    "fmt"

    "github.com/prfc0/consistenthash"
    "github.com/prfc0/consistenthash/hasher"
)

func main() {
    ring := consistenthash.New(50, hasher.Murmur64{})

    for i := 1; i <= 20; i++ {
  	  node := fmt.Sprintf("node%02d", i)
  	  ring.AddNode(node)
    }

	  for i := 1; i <= 100; i++ {
		  key := fmt.Sprintf("key%03d", i)
		  owner, _ := ring.GetOwner(key)
		  fmt.Println("Primary owner:", owner)
	  }

    key := "key042"
    owner, _ := ring.GetOwner(key)
    fmt.Println("Primary owner:", owner)
    // Primary owner: node17

    replicas, _ := ring.GetReplicas(key, 3)
    fmt.Println("Replicas:", replicas)
    // Replicas: [node17 node15 node01]

    ring.RemoveNode("node15")
    fmt.Println("Nodes:", ring.Nodes())
    // Nodes: [node09 node11 node04 node13 node16 node17 node18 node10 node12 node14 node02 node03 node05 node06 node08 node19 node20 node01 node07]
}
```

## API

* `New(replicas int, h hasher.Hasher) *Ring`
* `AddNode(node string)`
* `RemoveNode(node string)`
* `Nodes() []string`
* `Slots() []Slot`
* `GetOwner(key string) (string, error)`
* `GetReplicas(key string, n int) ([]string, error)`

## Acknowledgment

This package was developed with help of the OpenAI ChatGPT-5 model.
