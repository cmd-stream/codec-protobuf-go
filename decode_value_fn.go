package codec

import (
	"reflect"

	gnrc "github.com/cmd-stream/codec-generic-go"
)

// decodeValue decodes a value of the given reflect.Type using the provided
// Serializer.
func decodeValue[T, V any](tp reflect.Type, ser gnrc.Serializer[T, V],
	bs []byte,
) (v V, err error) {
	tp = tp.Elem()
	ptr := reflect.New(tp)
	err = ser.Unmarshal(bs, ptr.Interface().(V))
	if err != nil {
		err = gnrc.NewFailedToUnmarshalValue(err)
		return
	}
	v = ptr.Interface().(V)
	return
}
