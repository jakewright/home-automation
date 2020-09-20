# JRPC

JRPC generates code from `def` files to assist with JSON-based RPCs.

### Types

A message field can have one of the following types. The table shows the corresponding generated types. Fields can be marked as repeated by prepending the type with `[]`.

| JRPC type   | Golang type   |
| ----------- | ------------- |
| `bool`      | `bool`        |
| `string`    | `string`      |
| `int32`     | `int32`       |
| `int64`     | `int64`       |
| `uint32`    | `uint32`      |
| `uint64`    | `uint64`      |
| `float32`   | `float32`     |
| `float64`   | `float64`     |
| `bytes`     | `[]byte`      |
| `time`      | `time.Time`   |
| `any`       | `interface{}` |
| `map[x]y`   | `map[x]y`     |
| `rgb`       | `util.RGB`    |

### Field options

Message fields can take various options which are used to generate validation functions. The router code (`template_router.go`) automatically calls the validation functions in the generated handlers.

**`required`** If set, the _zero value_ is not allowed. For pointer (`*`) fields, they must be set. For non-pointer fields, they must not be zero.

**`min`** Can be used on numeric fields to enforce a minimum allowed value.

**`max`** Can be used on numeric fields to enforce a maximum allowed value.
