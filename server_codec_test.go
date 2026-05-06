package codec_test

import (
	"errors"
	"reflect"
	"testing"

	cdcproto "github.com/cmd-stream/codec-protobuf-go"
	"github.com/cmd-stream/codec-protobuf-go/test"
	com "github.com/mus-format/common-go"
	"google.golang.org/protobuf/proto"

	cmock "github.com/cmd-stream/cmd-stream-go/test/mock"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestServerCodec_Encode(t *testing.T) {
	var (
		wantDTM   = 0
		result    = &test.Result1{X: 10}
		wantBs, _ = proto.Marshal(result)
		wantLen   = len(wantBs)
		wantN     = 1 + 1 + wantLen
		writer    = cmock.NewWriter()
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
	codec := cdcproto.NewServerCodec[any](
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
		writer  = cmock.NewWriter()
	)
	writer.RegisterWriteByte(func(b byte) error {
		return wantErr
	})
	codec := cdcproto.NewServerCodec[any](
		[]reflect.Type{reflect.TypeFor[*test.Cmd1]()},
		[]reflect.Type{reflect.TypeFor[*test.Result1]()},
	)
	_, err := codec.Encode(result, writer)
	assertfatal.EqualDeep(t, errors.Is(err, wantErr), true)
	assertfatal.EqualDeep(t, err.Error()[:len(cdcproto.ErrorPrefix)], cdcproto.ErrorPrefix)
}

func TestServerCodec_Decode(t *testing.T) {
	var (
		wantDTM   = 1
		wantV     = &test.Cmd2{Y: "hello"}
		wantBs, _ = proto.Marshal(wantV)
		wantLen   = len(wantBs)
		wantN     = 1 + 1 + wantLen
		reader    = cmock.NewReader()
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
	codec := cdcproto.NewServerCodec[any](
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
		reader  = cmock.NewReader()
	)
	reader.RegisterReadByte(func() (b byte, err error) {
		return 0, wantErr
	})
	codec := cdcproto.NewServerCodec[any](
		[]reflect.Type{reflect.TypeFor[*test.Cmd1]()},
		[]reflect.Type{reflect.TypeFor[*test.Result1]()},
	)
	_, _, err := codec.Decode(reader)
	assertfatal.EqualDeep(t, errors.Is(err, wantErr), true)
	assertfatal.EqualDeep(t, err.Error()[:len(cdcproto.ErrorPrefix)], cdcproto.ErrorPrefix)
}
func TestNewServerCodecWith(t *testing.T) {
	var (
		wantDTM   = 1
		wantV     = &test.Cmd2{Y: "hello"}
		wantBs, _ = proto.Marshal(wantV)
		wantLen   = len(wantBs)
		wantN     = 1 + 1 + wantLen
		reader    = cmock.NewReader()
		reg       = cdcproto.NewRegistry[any](
			cdcproto.WithCmd[any, *test.Cmd1](),
			cdcproto.WithCmd[any, *test.Cmd2](),
			cdcproto.WithResult[any, *test.Result1](),
			cdcproto.WithResult[any, *test.Result2](),
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
	codec := cdcproto.NewServerCodecWith(reg)
	v, n, err := codec.Decode(reader)
	assertfatal.EqualError(t, err, nil)
	assertfatal.Equal(t, n, wantN)
	assertfatal.EqualDeep(t, proto.Equal(v.(proto.Message), wantV), true)
}

func TestServerCodec_MaxLenOption(t *testing.T) {
	var (
		maxLen = 5
		reader = cmock.NewReader()
	)
	reader.RegisterReadByte(
		func() (b byte, err error) { return 0, nil },
	).RegisterReadByte(
		func() (b byte, err error) { return byte(10), nil },
	)
	codec := cdcproto.NewServerCodec[any](
		[]reflect.Type{reflect.TypeFor[test.Cmd2]()},
		[]reflect.Type{reflect.TypeFor[test.Result1]()},
		cdcproto.WithMaxLen(maxLen),
	)

	_, _, err := codec.Decode(reader)
	assertfatal.EqualDeep(t, errors.Is(err, com.ErrTooLargeLength), true)
}
