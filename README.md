# ip [![Build Status](https://travis-ci.org/vinxi/ip.png)](https://travis-ci.org/vinxi/ip) [![GoDoc](https://godoc.org/github.com/vinxi/ip?status.svg)](https://godoc.org/github.com/vinxi/ip) [![Coverage Status](https://coveralls.io/repos/github/vinxi/ip/badge.svg?branch=master)](https://coveralls.io/github/vinxi/ip?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/vinxi/ip)](https://goreportcard.com/report/github.com/vinxi/ip)

IP range based filtering and multiplexer for your proxies. 

Supports IP v4/v6 expressions, CIDR ranges, subnets...

## Installation

```bash
go get -u gopkg.in/vinxi/ip.v0
```

## API

See [godoc](https://godoc.org/github.com/vinxi/ip) reference.

## Example

#### CIDR range filtering

```go
package main

import (
  "fmt"
  "time"

  "gopkg.in/vinxi/ip.v0"
  "gopkg.in/vinxi/vinxi.v0"
)

const port = 3100

func main() {
  // Create a new vinxi proxy
  vs := vinxi.NewServer(vinxi.ServerOptions{Port: port})

  // Attach the rate limit middleware for 10 req/min
  vs.Use(ip.New("127.0.0.1/8"))

  // Target server to forward
  vs.Forward("http://httpbin.org")

  fmt.Printf("Server listening on port: %d\n", port)
  err := vs.Listen()
  if err != nil {
    fmt.Errorf("Error: %s\n", err)
  }
}
```

## License

MIT
