package codec

import (
	"google.golang.org/protobuf/proto"
)

type Serializer[T, V proto.Message] struct{}

func (s Serializer[T, V]) Marshal(t T) (bs []byte, err error) {
	return proto.Marshal(t)
}

func (s Serializer[T, V]) Unmarshal(bs []byte, v V) (err error) {
	return proto.Unmarshal(bs, v)
}
