package codec

import (
	"reflect"

	"github.com/cmd-stream/core-go"
)

// NewServerCodec creates a Protobuf Codec for the server side.
//
// The cmdTypes slice lists Command types the server can handle.
// The resultTypes slice lists Result types that can be returned to the client.
//
// Note: The order of types matters — two codecs created with the same types
// in a different order are not considered equal.
//
// Note: All Command and Result types must implement proto.Message.
// The codec will panic if a type does not satisfy this requirement.
func NewServerCodec[T any](cmdTypes []reflect.Type, resultTypes []reflect.Type) (
	codec Codec[core.Result, core.Cmd[T]],
) {
	return NewCodec[core.Result, core.Cmd[T]](resultTypes, cmdTypes)
}
