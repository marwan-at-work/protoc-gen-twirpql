---
title: 'Types'
weight: 3
---


Protocol Buffer types are not exactly mapped to GraphQL types. Therefore, TwirpQL does a number of heavy lifting when certain types are not consistent. 

### Enums

Protobuf enums are represented as `int32` types while GraphQL enums are a `String` type. Therefore, TwirpQL patches a converter between the two types so that dealing with enums feels natural. The `String` representation will be exactly how the enum was defined in the `service.proto` file. For example, if you have the following enum in `service.proto`

```proto
enum Traffic {
    RED = 0;
    YELLOW = 1;
    GREEN = 2;
}
```

Then TwirpQL will create the following schema: 

```graphql
enum Traffic {
    RED
    YELLOW
    GREEN
}
```

And therefore, the values of this GraphQL type will be one of `"RED", "YELLOW", or "GREEN"`.

### Maps

GraphQL does not yet have support for arbitrary key-value maps. See https://github.com/graphql/graphql-spec/issues/101 for more context.

However, Protocol Buffer supports arbitrary maps inside messages such as 

```proto
message MyMessage {
    map<int64, string> myMap = 1;
}
```

In the example above, TwirpQL will create a custom `Scalar` type for GraphQL to interpret as a string for input queries, and an untyped JSON Object in the response data. 

Under the hood, TwirpQL will take care of converting strings into the Go maps. 

Therefore, your query input can be: 

```json
{
    "req": {
        "myMap": "{\"33\": \"thirty three\"}"
    }
}
```

And the query response will look something like this: 

```json
{
    "data": {
        "myQuery": {
            "myMap": {
                "33": "thirty three"
            }
        }
    }
}
```


### Messages, Inputs, and Types

Defining an object in Protocol Buffers are done through the `message` type declaration. 
Therefore, both a request's input and output can refer to the same `message`. For example: 

```proto
service Service {
    rpc Hello(MyMessage) returns (MyMessage);
}

message MyMessage {
    string text = 1;
}
```

Notice, that the same "message declaration" is *both* the input and the output of the Hello RPC.

However, GraphQL makes a clear distinction between an RPC's input, and output. For example, take the following GraphQL Schema File: 

```graphql
type Query {
    Hello(req: HelloReq!): HelloResp!
}

input HelloReq {
    text: String!
}

type HelloResp {
    text: String!
}
```

Notice that the `HelloReq` object was declared with the `input` keyword while the `HelloResp` object was declared with the `type` keyword. 

Furthermore, GraphQL does not allow an `input` and a `type` to have the same name. And so converting the above Protobuf file to the following will **not** work: 

```graphql
type Query {
    Hello(req: MyMessage!): MyMessage!
}

input MyMessage {
    text: String!
}

type MyMessage {
    text: String!
}
```

Therefore, TwirpQL will need to adjust the name of either type declarations so that we can avoid name clashes. TwirpQL in this case chooses to append the word `Input` at the end of `MyMessage` as such: 

```graphql
type Query {
    Hello(req: MyMessage!): MyMessage!
}

input MyMessageInput {
    text: String!
}

type MyMessage {
    text: String!
}
```


### Empty Messages

In Protocol Buffers, a `message` declaration can be empty. However, GraphQL does not allow empty `input`/`type` declarations. 

TwirpQL does two things: 

1. If the `input` declaration is empty, then the input is removed. This is so that querying things is simpler. 
2. If the `type` declaration is empty, then TwirpQL makes up a fake field to make GraphQL happy. 

Please note that this behavior may change based on what makes the most sense in terms of maintainability. The obvious downside of this is that when you introduce a new field to the empty `message`, then the GraphQL contract has a breaking change. 