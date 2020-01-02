# protoc-gen-jrpc

protoc-gen-jrpc is a protobuf plugin that generates code to easily
make JSON-based RPCs between services.

A proto file for a service would look something like the following:

```
syntax = "proto3";

package fooproto;
option go_package = "github.com/jakewright/home-automation/service.foo/proto;fooproto";

// This is needed to support the custom options
import "tools/protoc-gen-jrpc/proto/jrpc.proto";

service Foo {
    // Service name; used for routing.
    option (router).name = "service.foo";

    rpc Foo (FooRequest) returns (FooResponse) {
        option (handler).method = "GET";
        option (handler).path = "/foo";
    }
}

message FooRequest {}

message FooResponse {}
``` 
