# JRPC

JRPC generates code from .def files to assist with JSON-based RPCs.

### Types

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
