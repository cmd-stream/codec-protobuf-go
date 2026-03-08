package codec_test

import (
	"reflect"
	"testing"

	"github.com/cmd-stream/codec-protobuf-go"
	"github.com/cmd-stream/codec-protobuf-go/test/cmds"
	"github.com/cmd-stream/codec-protobuf-go/test/results"
	"google.golang.org/protobuf/proto"

	tmock "github.com/cmd-stream/transport-go/test/mock"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestClientCodec(t *testing.T) {
	t.Run("Encoding should succeed", func(t *testing.T) {
		wantDTM := 0
		cmd := &cmds.Cmd1{X: 10}
		wantBs, err := proto.Marshal(cmd)
		assertfatal.EqualError(t, err, nil)
		wantLen := len(wantBs)
		wantN := 1 + 1 + wantLen

		c := codec.NewClientCodec[any](
			[]reflect.Type{
				reflect.TypeFor[*cmds.Cmd1](),
				reflect.TypeFor[*cmds.Cmd2](),
			},
			[]reflect.Type{
				reflect.TypeFor[*results.Result1](),
				reflect.TypeFor[*results.Result2](),
			},
		)

		w := tmock.NewWriter().RegisterWriteByte(func(b byte) error {
			assertfatal.Equal(t, b, byte(wantDTM))
			return nil
		}).RegisterWriteByte(func(b byte) error {
			assertfatal.Equal(t, b, byte(wantLen))
			return nil
		}).RegisterWrite(func(p []byte) (n int, err error) {
			assertfatal.EqualDeep(t, p, wantBs)
			return len(p), nil
		})

		n, err := c.Encode(cmd, w)
		assertfatal.EqualError(t, err, nil)
		assertfatal.Equal(t, n, wantN)
	})

	t.Run("Decoding should succeed", func(t *testing.T) {
		wantDTM := 1
		wantV := &results.Result2{Y: "hello"}
		wantBs, err := proto.Marshal(wantV)
		assertfatal.EqualError(t, err, nil)
		wantLen := len(wantBs)
		wantN := 1 + 1 + wantLen

		c := codec.NewClientCodec[any](
			[]reflect.Type{
				reflect.TypeFor[*cmds.Cmd1](),
				reflect.TypeFor[*cmds.Cmd2](),
			},
			[]reflect.Type{
				reflect.TypeFor[*results.Result1](),
				reflect.TypeFor[*results.Result2](),
			},
		)

		r := tmock.NewReader().RegisterReadByte(func() (b byte, err error) {
			return byte(wantDTM), nil
		}).RegisterReadByte(func() (b byte, err error) {
			return byte(wantLen), nil
		}).RegisterRead(func(p []byte) (n int, err error) {
			copy(p, wantBs)
			return wantLen, nil
		})

		v, n, err := c.Decode(r)
		assertfatal.EqualError(t, err, nil)
		assertfatal.Equal(t, n, wantN)
		assertfatal.EqualDeep(t, proto.Equal(v.(proto.Message), wantV), true)
	})
}
