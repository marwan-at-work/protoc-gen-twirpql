---
title: 'Install TwirpQL'
weight: 2
---


### In order to use TwirpQL you need to:

1. Install the [Go](https://golang.org) toolchain (1.12+). 
2. Install the latest version of the [Protobuf Compiler](https://github.com/protocolbuffers/protobuf/releases) (v3.8+)
3. Install the Twirp plugin:

    `GO111MODULE=on go install github.com/twitchtv/twirp/protoc-gen-twirp@v5.7.0`

4. Install the TwirpQL plugin:

    `GO111MODULE=on go install marwan.io/protoc-gen-twirpql`

Next: [Quick Start](/docs/quick-start)