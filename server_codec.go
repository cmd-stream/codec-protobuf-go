package codec

import (
	"reflect"
)

// NewServerCodec creates a Protobuf Codec for the server side.
//
// The cmdTypes slice lists Command types the server can handle.
// The resultTypes slice lists Result types that can be returned to the client.
//
// Note: The order of types matters — two codecs created with the same types
// in a different order are not considered equal.
func NewServerCodec[T any](cmdTypes []reflect.Type, resultTypes []reflect.Type) (
	codec Codec[Result, Cmd[T]],
) {
	return NewCodec[Result, Cmd[T]](resultTypes, cmdTypes)
}
