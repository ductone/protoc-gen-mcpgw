package v1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type ServiceDesc struct {
	Name        string
	HandlerType interface{}
	Methods     []*MethodDesc
}

type methodHandler func(srv interface{}, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error)
type decoderHandler func(ctx context.Context, input DecoderInput, out proto.Message) error
type inputSchemaHandler func() (map[string]any)

type MethodDesc struct {
	Method        string
	Handler       methodHandler
	Decoder       decoderHandler
	InputSchema   inputSchemaHandler
	Title         string
	Description   string
	ReadOnlyHint  bool
	Destructive   bool
	Idempotent    bool
	OpenWorldHint bool
}

type ServiceRegistrar interface {
	RegisterService(sd *ServiceDesc, ss interface{})
}
