package codec_test

import (
	"errors"
	"reflect"
	"testing"

	cdc "github.com/cmd-stream/codec-protobuf-go"
	"github.com/cmd-stream/codec-protobuf-go/test"
	"google.golang.org/protobuf/proto"

	tmock "github.com/cmd-stream/cmd-stream-go/test/mock/transport"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestServerCodec_Encode(t *testing.T) {
	var (
		wantDTM   = 0
		result    = &test.Result1{X: 10}
		wantBs, _ = proto.Marshal(result)
		wantLen   = len(wantBs)
		wantN     = 1 + 1 + wantLen
		writer    = tmock.NewWriter()
	)
	writer.RegisterWriteByte(
		func(b byte) error {
			assertfatal.Equal(t, b, byte(wantDTM))
			return nil
		},
	).RegisterWriteByte(
		func(b byte) error {
			assertfatal.Equal(t, b, byte(wantLen))
			return nil
		},
	).RegisterWrite(
		func(p []byte) (n int, err error) {
			assertfatal.EqualDeep(t, p, wantBs)
			return len(p), nil
		},
	)
	codec := cdc.NewServerCodec[any](
		[]reflect.Type{
			reflect.TypeFor[*test.Cmd1](),
			reflect.TypeFor[*test.Cmd2](),
		},
		[]reflect.Type{
			reflect.TypeFor[*test.Result1](),
			reflect.TypeFor[*test.Result2](),
		},
	)
	n, err := codec.Encode(result, writer)
	assertfatal.EqualError(t, err, nil)
	assertfatal.Equal(t, n, wantN)
}

func TestServerCodec_EncodeError(t *testing.T) {
	var (
		result  = &test.Result1{X: 10}
		wantErr = errors.New("write error")
		writer  = tmock.NewWriter()
	)
	writer.RegisterWriteByte(func(b byte) error {
		return wantErr
	})
	codec := cdc.NewServerCodec[any](
		[]reflect.Type{reflect.TypeFor[*test.Cmd1]()},
		[]reflect.Type{reflect.TypeFor[*test.Result1]()},
	)
	_, err := codec.Encode(result, writer)
	assertfatal.EqualDeep(t, errors.Is(err, wantErr), true)
	assertfatal.EqualDeep(t, err.Error()[:len(cdc.ErrorPrefix)], cdc.ErrorPrefix)
}

func TestServerCodec_Decode(t *testing.T) {
	var (
		wantDTM   = 1
		wantV     = &test.Cmd2{Y: "hello"}
		wantBs, _ = proto.Marshal(wantV)
		wantLen   = len(wantBs)
		wantN     = 1 + 1 + wantLen
		reader    = tmock.NewReader()
	)
	reader.RegisterReadByte(
		func() (b byte, err error) { return byte(wantDTM), nil },
	).RegisterReadByte(
		func() (b byte, err error) { return byte(wantLen), nil },
	).RegisterRead(
		func(p []byte) (n int, err error) {
			copy(p, wantBs)
			return wantLen, nil
		},
	)
	codec := cdc.NewServerCodec[any](
		[]reflect.Type{
			reflect.TypeFor[*test.Cmd1](),
			reflect.TypeFor[*test.Cmd2](),
		},
		[]reflect.Type{
			reflect.TypeFor[*test.Result1](),
			reflect.TypeFor[*test.Result2](),
		},
	)
	v, n, err := codec.Decode(reader)
	assertfatal.EqualError(t, err, nil)
	assertfatal.Equal(t, n, wantN)
	assertfatal.EqualDeep(t, proto.Equal(v.(proto.Message), wantV), true)
}

func TestServerCodec_DecodeError(t *testing.T) {
	var (
		wantErr = errors.New("read error")
		reader  = tmock.NewReader()
	)
	reader.RegisterReadByte(func() (b byte, err error) {
		return 0, wantErr
	})
	codec := cdc.NewServerCodec[any](
		[]reflect.Type{reflect.TypeFor[*test.Cmd1]()},
		[]reflect.Type{reflect.TypeFor[*test.Result1]()},
	)
	_, _, err := codec.Decode(reader)
	assertfatal.EqualDeep(t, errors.Is(err, wantErr), true)
	assertfatal.EqualDeep(t, err.Error()[:len(cdc.ErrorPrefix)], cdc.ErrorPrefix)
}
func TestNewServerCodecWith(t *testing.T) {
	var (
		wantDTM   = 1
		wantV     = &test.Cmd2{Y: "hello"}
		wantBs, _ = proto.Marshal(wantV)
		wantLen   = len(wantBs)
		wantN     = 1 + 1 + wantLen
		reader    = tmock.NewReader()
		reg       = cdc.NewRegistry[any](
			cdc.WithCmd[any, *test.Cmd1](),
			cdc.WithCmd[any, *test.Cmd2](),
			cdc.WithResult[any, *test.Result1](),
			cdc.WithResult[any, *test.Result2](),
		)
	)
	reader.RegisterReadByte(
		func() (b byte, err error) { return byte(wantDTM), nil },
	).RegisterReadByte(
		func() (b byte, err error) { return byte(wantLen), nil },
	).RegisterRead(
		func(p []byte) (n int, err error) {
			copy(p, wantBs)
			return wantLen, nil
		},
	)
	codec := cdc.NewServerCodecWith(reg)
	v, n, err := codec.Decode(reader)
	assertfatal.EqualError(t, err, nil)
	assertfatal.Equal(t, n, wantN)
	assertfatal.EqualDeep(t, proto.Equal(v.(proto.Message), wantV), true)
}
