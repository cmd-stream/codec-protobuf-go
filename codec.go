// Package codec provides a Protobuf-based codec implementation for cmd-stream-go.
package codec

import (
	"fmt"
	"reflect"

	codecgnrc "github.com/cmd-stream/codec-generic-go"
	"github.com/cmd-stream/transport-go"
)

const errorPrefix = "codecproto: "

// NewCodec constructs a Protobuf Codec.
//
// Parameters:
//   - types1 lists the Go types that can be encoded.
//   - types2 lists the Go types that can be decoded.
func NewCodec[T, V any](types1 []reflect.Type,
	types2 []reflect.Type,
) Codec[T, V] {
	return Codec[T, V]{
		codecgnrc.NewCodecWithDecoder(types1, types2, Serializer[T, V]{},
			decodeValue),
	}
}

// Codec represents a generic type-safe codec for encoding and decoding values.
// T is the type used for encoding, V is the type used for decoding.
type Codec[T, V any] struct {
	codecgnrc.Codec[T, V]
}

// Encode writes a value of type T to the given transport.Writer.
// Returns the total number of bytes written and any error.
func (c Codec[T, V]) Encode(t T, w transport.Writer) (n int, err error) {
	n, err = c.Codec.Encode(t, w)
	if err != nil {
		err = fmt.Errorf(errorPrefix+"%w", err)
	}
	return
}

// Decode reads a value of type V from the given transport.Reader.
// Returns the decoded value, total bytes read, and any error.
func (c Codec[T, V]) Decode(r transport.Reader) (v V, n int, err error) {
	v, n, err = c.Codec.Decode(r)
	if err != nil {
		err = fmt.Errorf(errorPrefix+"%w", err)
	}
	return
}
