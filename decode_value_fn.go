package codec

import (
	"reflect"

	codecgnrc "github.com/cmd-stream/codec-generic-go"
)

func decodeValue[T, V any](tp reflect.Type, ser codecgnrc.Serializer[T, V],
	bs []byte,
) (v V, err error) {
	tp = tp.Elem()
	ptr := reflect.New(tp)
	err = ser.Unmarshal(bs, ptr.Interface().(V))
	if err != nil {
		err = codecgnrc.NewFailedToUnmarshalValue(err)
		return
	}
	v = ptr.Interface().(V)
	return
}
