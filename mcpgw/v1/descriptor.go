package v1

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type ServiceDesc struct {
	Name        string
	HandlerType interface{}
	Methods     []MethodDescInterface
}

type methodHandler[TReq, TResp proto.Message] func(srv interface{}, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error)
type decoderHandler func(ctx context.Context, input DecoderInput, out proto.Message) error
type inputSchemaHandler func() map[string]any

// MethodDesc is a generic type that captures the specific request and response types
type MethodDesc[TReq, TResp proto.Message] struct {
	Method        string
	Handler       methodHandler[TReq, TResp]
	Decoder       decoderHandler
	InputSchema   inputSchemaHandler
	Title         string
	Description   string
	ReadOnlyHint  bool
	Destructive   bool
	Idempotent    bool
	OpenWorldHint bool
}

// MethodDescInterface provides a common interface for all MethodDesc types
type MethodDescInterface interface {
	GetMethod() string
	GetHandler() func(srv interface{}, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error)
	GetDecoder() decoderHandler
	GetInputSchema() inputSchemaHandler
	GetTitle() string
	GetDescription() string
	GetReadOnlyHint() bool
	GetDestructive() bool
	GetIdempotent() bool
	GetOpenWorldHint() bool
	AssertRequestType(ctx context.Context, in proto.Message) (proto.Message, error)
	AssertResponseType(ctx context.Context, in proto.Message) (proto.Message, error)

	// Get concrete types at runtime
	GetRequestType() reflect.Type
	GetResponseType() reflect.Type
	CreateRequest() proto.Message
	CreateResponse() proto.Message
}

// Implement MethodDescInterface for MethodDesc[TReq, TResp]
func (md *MethodDesc[TReq, TResp]) GetMethod() string {
	return md.Method
}

func (md *MethodDesc[TReq, TResp]) GetHandler() func(srv interface{}, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error) {
	return md.Handler
}

func (md *MethodDesc[TReq, TResp]) GetDecoder() decoderHandler {
	return md.Decoder
}

func (md *MethodDesc[TReq, TResp]) GetInputSchema() inputSchemaHandler {
	return md.InputSchema
}

func (md *MethodDesc[TReq, TResp]) GetTitle() string {
	return md.Title
}

func (md *MethodDesc[TReq, TResp]) GetDescription() string {
	return md.Description
}

func (md *MethodDesc[TReq, TResp]) GetReadOnlyHint() bool {
	return md.ReadOnlyHint
}

func (md *MethodDesc[TReq, TResp]) GetDestructive() bool {
	return md.Destructive
}

func (md *MethodDesc[TReq, TResp]) GetIdempotent() bool {
	return md.Idempotent
}

func (md *MethodDesc[TReq, TResp]) GetOpenWorldHint() bool {
	return md.OpenWorldHint
}

// AssertRequestType converts a proto.Message to the specific request type for this method
func (md *MethodDesc[TReq, TResp]) AssertRequestType(ctx context.Context, in proto.Message) (proto.Message, error) {
	if in == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Try to convert to the specific type
	if req, ok := in.(TReq); ok {
		return req, nil
	}

	return nil, fmt.Errorf("expected request type %T, got %T", (*TReq)(nil), in)
}

// AssertResponseType converts a proto.Message to the specific response type for this method
func (md *MethodDesc[TReq, TResp]) AssertResponseType(ctx context.Context, in proto.Message) (proto.Message, error) {
	if in == nil {
		return nil, errors.New("response cannot be nil")
	}

	// Try to convert to the specific type
	if resp, ok := in.(TResp); ok {
		return resp, nil
	}

	return nil, fmt.Errorf("expected response type %T, got %T", (*TResp)(nil), in)
}

// GetRequestType returns the reflect.Type of the request type
func (md *MethodDesc[TReq, TResp]) GetRequestType() reflect.Type {
	// Get the type of TReq (which is a pointer type)
	return reflect.TypeOf((*TReq)(nil)).Elem()
}

// GetResponseType returns the reflect.Type of the response type
func (md *MethodDesc[TReq, TResp]) GetResponseType() reflect.Type {
	// Get the type of TResp (which is a pointer type)
	return reflect.TypeOf((*TResp)(nil)).Elem()
}

// CreateRequest creates a new instance of the request type
func (md *MethodDesc[TReq, TResp]) CreateRequest() proto.Message {
	// TReq is a pointer type, so get the element type and create a new instance
	ptrType := reflect.TypeOf((*TReq)(nil)).Elem() // e.g., *v1.CreateGenreRequest
	val := reflect.New(ptrType.Elem()).Interface() // new(v1.CreateGenreRequest)
	return val.(proto.Message)
}

// CreateResponse creates a new instance of the response type
func (md *MethodDesc[TReq, TResp]) CreateResponse() proto.Message {
	ptrType := reflect.TypeOf((*TResp)(nil)).Elem()
	val := reflect.New(ptrType.Elem()).Interface()
	return val.(proto.Message)
}

// Helper function to convert proto.Message to concrete type using reflection
func ConvertToConcreteType(method MethodDescInterface, in proto.Message, isRequest bool) (proto.Message, error) {
	if in == nil {
		return nil, errors.New("input cannot be nil")
	}

	var expectedType reflect.Type
	if isRequest {
		expectedType = method.GetRequestType()
	} else {
		expectedType = method.GetResponseType()
	}

	// Check if the input is already the correct type
	if reflect.TypeOf(in) == expectedType {
		return in, nil
	}

	// Try to convert using reflection
	val := reflect.ValueOf(in)
	if val.Type() == expectedType {
		return in, nil
	}

	// If it's a pointer to the correct type, return it
	if val.Type().Elem() == expectedType.Elem() {
		return in, nil
	}

	typeName := "response"
	if isRequest {
		typeName = "request"
	}

	return nil, fmt.Errorf("expected %s type %v, got %v", typeName, expectedType, reflect.TypeOf(in))
}

type ServiceRegistrar interface {
	RegisterService(sd *ServiceDesc, ss interface{})
}
