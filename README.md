# ip [![Build Status](https://travis-ci.org/vinxi/ip.png)](https://travis-ci.org/vinxi/ip) [![GoDoc](https://godoc.org/github.com/vinxi/ip?status.svg)](https://godoc.org/github.com/vinxi/ip) [![Coverage Status](https://coveralls.io/repos/github/vinxi/ip/badge.svg?branch=master)](https://coveralls.io/github/vinxi/ip?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/vinxi/ip)](https://goreportcard.com/report/github.com/vinxi/ip)

IP range based filtering and multiplexer for your proxies. 
Supports IP v4/v6 expressions, CIDR ranges, subnets...

Implements a middleware layer, so you can use it as multiplexer.

## Installation

```bash
go get -u gopkg.in/vinxi/ip.v0
```

## API

See [godoc](https://godoc.org/github.com/vinxi/ip) reference.

## Example

#### Allow reserved loopback CIDR ranges

```go
package main

import (
  "fmt"

  "gopkg.in/vinxi/ip.v0"
  "gopkg.in/vinxi/vinxi.v0"
)

const port = 3100

func main() {
  // Create a new vinxi proxy
  vs := vinxi.NewServer(vinxi.ServerOptions{Port: port})

  // Attach the IP range filter
  vs.Use(ip.New("127.0.0.1/8", "::1/64"))

  // Target server to forward
  vs.Forward("http://httpbin.org")

  fmt.Printf("Server listening on port: %d\n", port)
  err := vs.Listen()
  if err != nil {
    fmt.Errorf("Error: %s\n", err)
  }
}
```

#### Use filter as multiplexer for middleware composition

```go
package main

import (
  "fmt"
  "net/http"

  "gopkg.in/vinxi/ip.v0"
  "gopkg.in/vinxi/vinxi.v0"
)

const port = 3100

func main() {
  // Create a new vinxi proxy
  vs := vinxi.NewServer(vinxi.ServerOptions{Port: port})

  // Create the IP filter
  filter := ip.New("127.0.0.1/8", "::1/64")

  // Attach a filter level middleware handler
  filter.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
    w.Header().Set("IP-Allowed", r.RemoteAddr)
    h.ServeHTTP(w, r)
  })

  // Register the filter in vinxi
  vs.Use(filter)

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
