package codec

import (
	"reflect"
)

// NewClientCodec creates a Protobuf Codec for the client side.
//
// The cmdTypes slice lists Command types the client can send.
// The resultTypes slice lists Result types the client expects to receive.
//
// Note: The order of types matters — two codecs created with the same types
// in a different order are not considered equal.
func NewClientCodec[T any](cmdTypes []reflect.Type, resultTypes []reflect.Type) (
	codec Codec[Cmd[T], Result],
) {
	return NewCodec[Cmd[T], Result](cmdTypes, resultTypes)
}
