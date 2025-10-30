package codec

import (
	"github.com/cmd-stream/core-go"

	"google.golang.org/protobuf/proto"
)

// Cmd is a type that can be encoded and decoded by the codec.
type Cmd[T any] interface {
	core.Cmd[T]
	proto.Message
}
