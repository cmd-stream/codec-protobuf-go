# codec-protobuf-go

[![Go Reference](https://pkg.go.dev/badge/github.com/cmd-stream/codec-protobuf-go.svg)](https://pkg.go.dev/github.com/cmd-stream/codec-protobuf-go)
[![GoReportCard](https://goreportcard.com/badge/cmd-stream/codec-protobuf-go)](https://goreportcard.com/report/github.com/cmd-stream/codec-protobuf-go)
[![codecov](https://codecov.io/gh/cmd-stream/codec-protobuf-go/graph/badge.svg?token=3ACEL5m6LC)](https://codecov.io/gh/cmd-stream/codec-protobuf-go)

**codec-protobuf-go** provides a Protobuf-based codec for [cmd-stream-go](https://github.com/cmd-stream/cmd-stream-go).

It maps concrete Command and Result types to internal identifiers,
allowing type-safe serialization across network boundaries.

## How To

```go
import (
  "reflect"
  codec "github.com/cmd-stream/codec-protobuf-go"
)

var (
  // Note: The order of types matters — two codecs created with the same types
  // in a different order are not considered equal.
  cmdTypes = []reflect.Type{
    reflect.TypeFor[*YourCmd](),
    // ...
  }
  resultTypes = []reflect.Type{
    reflect.TypeFor[*YourResult](),
    // ...
  }
  serverCodec = codec.NewServerCodec(cmdTypes, resultTypes)
  clientCodec = codec.NewClientCodec(cmdTypes, resultTypes)
)
```
