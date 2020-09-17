# Service Definition Language

Svcdef is a service definition language. It's similar to the protobuf language but doesn't have any code generation built-in: that's up to you.

## Language specification

### Imports

Import other `def` files using the `import` statement. Paths must be relative to the current file. The alias must be supplied and used when referring to the imported file.

```
import foo "../service.foo/foo.def"
       ⬑ an alias must be included
           ⬑ the path must be a quoted string 
```

### Service definition

Only one `service` block can exist per `def` file. If a second `service` definition is found, a parsing error will occur.

```
service Foo {}
```

#### Service options

Services can have arbitrary options.

```
foo = "bar"
⬑ must be a valid identifier
      ⬑ the value can be a quoted string, a number, or a boolean (true or false)
```

#### RPC definition

A service is made up of RPCs. Each RPC has a request (input) type and a response (output) type. The types can:
  1. refer to messages defined in the same `def` file 
     - this will be determined automatically if the name matches the name of one of the messages in the file
     - nested messages can be referenced (although not recommended)
  2. refer to messages defined in an imported file 
     - this is inferred if the format `alias.MessageName` is used
     - a parsing error will occur if no matching imported type can be found
     - nested messages can be reference (although not recommended)
  3. be a simple type (any valid identifier can be used)
     - maps, optional and repeated types are not valid as RPC types

```
rpc Foo(FooRequest) FooResponse {    
    foo = "bar"
    ⬑ An RPC can also have options
}
```

#### Service example

```
service Foo {
    foo = "bar"
    baz = 500
    bat = true

    rpc ReadUser(ReadUserRequest) ReadUserResponse {
        method = "GET"
        path = "/read"
    }
}
```

### Message definitions

Messages can be thought of as type definitions.

```
message ReadUserRequest {
        ⬑ must be a valid identifier
    foo = "bar"
        ⬑ options can be defined as in service definitions
    string name
        ⬑ type names are arbitrary unless they include a period
    user.Address address
        ⬑ type names with a period are typically
          references to a type from an imported file
    []int numbers
        ⬑ prefixing a type with [] will mark it as repeated
    *bool marketing_emails
        ⬑ prefixing a type with a * will mark it as optional
    *[]string children
        ⬑ in this case, it is the list that is optional
          []* is not valid syntax
}
```
