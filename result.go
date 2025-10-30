package codec

import (
	"github.com/cmd-stream/core-go"

	"google.golang.org/protobuf/proto"
)

// Result is a type that can be encoded and decoded by the codec.
type Result interface {
	core.Result
	proto.Message
}
