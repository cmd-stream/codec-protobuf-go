package codec_test

import (
	"testing"

	cdctest "github.com/cmd-stream/codec-go/test"
	cdcproto "github.com/cmd-stream/codec-protobuf-go"
	"github.com/cmd-stream/codec-protobuf-go/test"
	"google.golang.org/protobuf/proto"
)

func FuzzClientCodec_Decode(f *testing.F) {
	reg := cdcproto.NewRegistry(
		cdcproto.WithCmd[any, *test.Cmd1](),
		cdcproto.WithCmd[any, *test.Cmd2](),
		cdcproto.WithResult[any, *test.Result1](),
		cdcproto.WithResult[any, *test.Result2](),
	)

	// Seed with valid Protobuf data.
	f.Add([]byte{0, 2, 8, 10}, 100)                          // Cmd1{X: 10}
	f.Add([]byte{1, 7, 10, 5, 104, 101, 108, 108, 111}, 100) // Cmd2{Y: "hello"}
	f.Add([]byte{255, 0}, 10)                                // Invalid DTM

	f.Fuzz(func(t *testing.T, data []byte, maxLen int) {
		if maxLen <= 0 || maxLen > 10*1024*1024 {
			return
		}
		c := cdcproto.NewClientCodecWith(reg, cdcproto.WithMaxLen(maxLen))
		cdctest.FuzzDecode(c, data)
	})
}

func FuzzServerCodec_Decode(f *testing.F) {
	reg := cdcproto.NewRegistry(
		cdcproto.WithCmd[any, *test.Cmd1](),
		cdcproto.WithCmd[any, *test.Cmd2](),
		cdcproto.WithResult[any, *test.Result1](),
		cdcproto.WithResult[any, *test.Result2](),
	)

	f.Add([]byte{0, 2, 8, 100}, 100)                         // Result1{X: 100}
	f.Add([]byte{1, 7, 10, 5, 119, 111, 114, 108, 100}, 100) // Result2{Y: "world"}

	f.Fuzz(func(t *testing.T, data []byte, maxLen int) {
		if maxLen <= 0 || maxLen > 10*1024*1024 {
			return
		}
		s := cdcproto.NewServerCodecWith(reg, cdcproto.WithMaxLen(maxLen))
		cdctest.FuzzDecode(s, data)
	})
}

func FuzzRoundTrip_Cmd(f *testing.F) {
	var (
		reg = cdcproto.NewRegistry(
			cdcproto.WithCmd[any, *test.Cmd1](),
			cdcproto.WithResult[any, *test.Result1](),
		)
		client = cdcproto.NewClientCodecWith(reg)
		server = cdcproto.NewServerCodecWith(reg)
	)

	f.Add(int64(10))
	f.Fuzz(func(t *testing.T, x int64) {
		cmd := &test.Cmd1{X: x}
		cdctest.VerifyRoundTripCmdWith(t, client, server, cmd, func(expected, actual any) bool {
			return proto.Equal(expected.(proto.Message), actual.(proto.Message))
		})
	})
}

func FuzzRoundTrip_Result(f *testing.F) {
	var (
		reg = cdcproto.NewRegistry(
			cdcproto.WithCmd[any, *test.Cmd1](),
			cdcproto.WithResult[any, *test.Result1](),
		)
		client = cdcproto.NewClientCodecWith(reg)
		server = cdcproto.NewServerCodecWith(reg)
	)

	f.Add(int64(10))
	f.Fuzz(func(t *testing.T, x int64) {
		res := &test.Result1{X: x}
		cdctest.VerifyRoundTripResultWith(t, client, server, res, func(expected, actual any) bool {
			return proto.Equal(expected.(proto.Message), actual.(proto.Message))
		})
	})
}
