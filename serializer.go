package codec

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
)

// Serializer implements the gnrc.Serializer interface using Protobuf.
type Serializer[T, V any] struct{}

// Marshal encodes the given value of type T into a Protobuf byte slice.
// If T does not implement proto.Message, it panics.
func (s Serializer[T, V]) Marshal(t T) (bs []byte, err error) {
	pt, ok := any(t).(proto.Message)
	if !ok {
		panic(fmt.Sprintf(ErrorPrefix+"%v does not implement proto.Message",
			reflect.TypeOf(t)))
	}
	return proto.Marshal(pt)
}

// Unmarshal decodes the Protobuf byte slice into the given value of type V.
// If V does not implement proto.Message, it panics.
func (s Serializer[T, V]) Unmarshal(bs []byte, v V) (err error) {
	pv, ok := any(v).(proto.Message)
	if !ok {
		panic(fmt.Sprintf(ErrorPrefix+"%v does not implement proto.Message",
			reflect.TypeOf(v)))
	}
	return proto.Unmarshal(bs, pv)
}
