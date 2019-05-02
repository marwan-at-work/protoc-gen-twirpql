---
title: 'Multiple Services'
weight: 4
---

TwirpQL's architecture expects that for each `service` declaration in a Protobuf file, you get a full GraphQL layer under one sub-package. 

In other words, TwirpQL expects only one service for every code generation cycle. If your `.proto` file has only one service in it, then TwirpQL automatically chooses it for creating the GraphQL layer. 

However, if the `.proto` has multiple `service` declarations, then you must explicitly pick which service you want to generate a GraphQL layer for.

For example, say you have the following `service.proto` file: 

```proto
syntax = "proto3";
package hello;
option go_package = "hello";

service One {
    ...
}

service Two {
    ...
}
```

Then you must specify which service you want to create a GraphQL for in the following way:

```bash
protoc --twirpql_out=service=Two:. service.proto
```

### What happens if I want a GraphQL layer for multiple services? 

It's quite common that one server = one service. However, it's not uncommon that one *server* has multiple *services*. 

Therefore, you can generate multiple GraphQL layers in the same project. Here's how: 

Say you have 3 `.proto` files, each with 1 service: `svc1.proto, svc2.proto, svc3.proto` 

Then you can create the following command for each of them: 

```bash
protoc --twirpql_out=dest=svc1:. svc1.proto
protoc --twirpql_out=dest=svc2:. svc2.proto
protoc --twirpql_out=dest=svc3:. svc3.proto
```

This way, each generation cycle will create a *different* sub-package. By default, TwirpQL creates a sub-package called `twirpql` but in this case, we overrode this default and create different names for each sub-package. 

Once all sub-packages are created, then we can import them individually and place them under different URL paths: 


```golang
package main

import (
    "net/http"

    svc1 "./svc1"
    svc2 "./svc2"
    svc3 "./svc3"
)

func main() {
    http.Handle("/svc1", svc1.Handler(...))
    http.Handle("/svc2", svc2.Handler(...))
    http.Handle("/svc3", svc3.Handler(...))
}

```