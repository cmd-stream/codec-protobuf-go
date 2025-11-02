package codec_test

import (
	"reflect"
	"testing"

	"github.com/cmd-stream/codec-protobuf-go"
	"github.com/cmd-stream/codec-protobuf-go/internal"
	"google.golang.org/protobuf/proto"

	tmocks "github.com/cmd-stream/testkit-go/mocks/transport"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestClientCodec(t *testing.T) {
	t.Run("Encoding should work", func(t *testing.T) {
		wantDTM := 0
		cmd := &internal.Cmd1{X: 10}
		wantBs, err := proto.Marshal(cmd)
		assertfatal.EqualError(err, nil, t)
		wantLen := len(wantBs)
		wantN := 1 + 1 + wantLen

		c := codec.NewClientCodec[any](
			[]reflect.Type{
				reflect.TypeFor[*internal.Cmd1](),
				reflect.TypeFor[*internal.Cmd2](),
			},
			[]reflect.Type{
				reflect.TypeFor[*internal.Result1](),
				reflect.TypeFor[*internal.Result2](),
			},
		)

		w := tmocks.NewWriter().RegisterWriteByte(func(b byte) error {
			assertfatal.Equal(b, byte(wantDTM), t)
			return nil
		}).RegisterWriteByte(func(b byte) error {
			assertfatal.Equal(b, byte(wantLen), t)
			return nil
		}).RegisterWrite(func(p []byte) (n int, err error) {
			assertfatal.EqualDeep(p, wantBs, t)
			return len(p), nil
		})

		n, err := c.Encode(cmd, w)
		assertfatal.EqualError(err, nil, t)
		assertfatal.Equal(n, wantN, t)
	})

	t.Run("Decoding should work", func(t *testing.T) {
		wantDTM := 1
		wantV := &internal.Result2{Y: "hello"}
		wantBs, err := proto.Marshal(wantV)
		assertfatal.EqualError(err, nil, t)
		wantLen := len(wantBs)
		wantN := 1 + 1 + wantLen

		c := codec.NewClientCodec[any](
			[]reflect.Type{
				reflect.TypeFor[*internal.Cmd1](),
				reflect.TypeFor[*internal.Cmd2](),
			},
			[]reflect.Type{
				reflect.TypeFor[*internal.Result1](),
				reflect.TypeFor[*internal.Result2](),
			},
		)

		r := tmocks.NewReader().RegisterReadByte(func() (b byte, err error) {
			return byte(wantDTM), nil
		}).RegisterReadByte(func() (b byte, err error) {
			return byte(wantLen), nil
		}).RegisterRead(func(p []byte) (n int, err error) {
			copy(p, wantBs)
			return wantLen, nil
		})

		v, n, err := c.Decode(r)
		assertfatal.EqualError(err, nil, t)
		assertfatal.Equal(n, wantN, t)
		assertfatal.EqualDeep(proto.Equal(v.(proto.Message), wantV), true, t)
	})
}
