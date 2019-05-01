---
title: 'Quick Start'
date: 2019-02-11T19:27:37+10:00
weight: 2
---


Once you installed all the required depepndencies, this page will guide you through an end-to-end example.

### Create a new project

```bash
~ export GO111MODULE=on # no need to do this if you ar eoutside of GOPATH
~ mkdir hello && cd hello && go mod init hello
```

### Define your Twirp Service

Create a file named `service.proto` and populate it with the following:

```
syntax = "proto3";
package hello;
option go_package = "hello";

service Service {
  rpc Hello(HelloReq) returns (HelloResp);
}

message HelloReq {
  string name = 1;
}

message HelloResp {
  string text = 1;
}
```

### Generate Go/Twirp Files

Before you can generate a TWIRPQL layer, you need to generate the `.pb.go` and `.twirp.go` files so that we have the contract defined: 

```bash
~ protoc --go_out=. --twirp_out=. service.proto
```

### Implement the Twirp server

Let's make a new sub-package called `server` and implement the twirp server in it: 

```bash
~ mkdir server
```

Create a main.go file inside `server` and paste the following inside it:

```golang
package main

import (
	"context"
	"hello"
	"net/http"
)

type service struct{}

func (s *service) Hello(ctx context.Context, req *hello.HelloReq) (*hello.HelloResp, error) {
	return &hello.HelloResp{Text: req.GetName()}, nil
}

func main() {
    serviceImpl := &service{}
    http.Handle("/", hello.NewServiceServer(serviceImpl, nil))

    http.ListenAndServe(":9090", nil)
}
```

### Run and test the Twirp server

1. `go run server/main.go` 
2. From another terminal run: `~ curl -X POST -d '{"name": "twirpql"}' -H 'Content-Type: application/json' localhost:9090/twirp/hello.Service/Hello`

### Create the TWIRPQL layer

Now that we verified we have a Twirp server working, let's create the TWIRPQL layer on top of it: 

```bash
# ctrl+c to stop the server, then:
~ protoc --twirpql_out=. service.proto
```

This will create a new subdirectory called `twirpql` with all the necessary files to import a GraphQL layer on top of our Hello service.

### Update server/main.go to use the new GraphQL layer

Update our main.go file to look like this:

```golang
package main

import (
	"context"
	"hello"
	"hello/twirpql"
	"net/http"
)

type service struct{}

func (s *service) Hello(ctx context.Context, req *hello.HelloReq) (*hello.HelloResp, error) {
	return &hello.HelloResp{Text: req.GetName()}, nil
}

func main() {
    serviceImpl := &service{}
	http.Handle("/", hello.NewServiceServer(serviceImpl, nil))
	http.Handle("/query", twirpql.Handler(serviceImpl, nil))
	http.Handle("/play", twirpql.Playground("twirp", "/query"))

	http.ListenAndServe(":9090", nil)
}
```

As you can see, we only had to add 3 lines of code: 

1. We imported the newly created subpackage `"hello/twirpql"`
2. We created a `/query` endpoint and called `twirpql.Handler` which knows to take your Twirp Service implementation and returns a GraphQL http.Handler. 
3. We created a `/play` endpoint which exposes a *GraphiQL* UI to discover and play with the Twirp service. 

### Rerun the server and explore the UI

```bash
~ go run server/main.go
```

And navigate to the [localhost:9090/play](http://localhost:9090/play) on the browser: 

<a href="/img/graphiql.png">
![Graphiql Image](/img/graphiql.png)
</a>


You can now use GraphQL to make queries back to the Twirp Service as you see in the screenshot above. 

### Pro Tip

To make re-generation easier, create a `gen.go` file in your project's root directory with the following content: 

```golang
package hello

//go:generate protoc --go_out=. --twirp_out=. service.proto
//go:generate protoc --twirpql_out=. service.proto
```

Then, whenever you update `service.proto`, all you have to do is run `go generate` from the root directory. 