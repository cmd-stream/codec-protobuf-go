package codec

import (
	"reflect"

	"github.com/cmd-stream/cmd-stream-go/core"
)

// ServerCodec defines a Protobuf codec for the server side.
type ServerCodec[T any] = codec[core.Result, core.Cmd[T]]

// NewServerCodec creates a Protobuf codec for the server side.
//
// The cmdTypes slice lists Command types the server can handle.
// The resultTypes slice lists Result types that can be returned to the client.
//
// Note: All Command and Result types must implement proto.Message.
// The codec will panic if a type does not satisfy this requirement.
//
// Note: The order of types matters — two codecs created with the same types
// in a different order are not considered equal.
func NewServerCodec[T any](cmdTypes []reflect.Type, resultTypes []reflect.Type) (
	codec ServerCodec[T],
) {
	return newCodec[core.Result, core.Cmd[T]](resultTypes, cmdTypes)
}
