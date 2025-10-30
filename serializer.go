package codec

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
)

type Serializer[T, V any] struct{}

func (s Serializer[T, V]) Marshal(t T) (bs []byte, err error) {
	pt, ok := any(t).(proto.Message)
	if !ok {
		panic(fmt.Sprintf(errorPrefix+"%v does not implement proto.Message",
			reflect.TypeOf(t)))
	}
	return proto.Marshal(pt)
}

func (s Serializer[T, V]) Unmarshal(bs []byte, v V) (err error) {
	pv, ok := any(v).(proto.Message)
	if !ok {
		panic(fmt.Sprintf(errorPrefix+"%v does not implement proto.Message",
			reflect.TypeOf(v)))
	}
	return proto.Unmarshal(bs, pv)
}
