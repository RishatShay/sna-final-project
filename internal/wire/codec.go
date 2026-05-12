package wire

import (
	"fmt"

	oldproto "github.com/golang/protobuf/proto"
	"google.golang.org/grpc/encoding"
	newproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"
)

type protobufCodec struct{}

func init() {
	encoding.RegisterCodec(protobufCodec{})
}

func (protobufCodec) Name() string {
	return "proto"
}

func (protobufCodec) Marshal(value any) ([]byte, error) {
	switch msg := value.(type) {
	case newproto.Message:
		return newproto.Marshal(msg)
	case oldproto.Message:
		return newproto.Marshal(protoadapt.MessageV2Of(msg))
	default:
		return nil, fmt.Errorf("cannot marshal %T as protobuf", value)
	}
}

func (protobufCodec) Unmarshal(data []byte, value any) error {
	switch msg := value.(type) {
	case newproto.Message:
		return newproto.Unmarshal(data, msg)
	case oldproto.Message:
		return newproto.Unmarshal(data, protoadapt.MessageV2Of(msg))
	default:
		return fmt.Errorf("cannot unmarshal protobuf into %T", value)
	}
}
