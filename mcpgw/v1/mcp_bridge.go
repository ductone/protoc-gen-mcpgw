package v1

import (
	"context"
	"crypto/tls"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type contextKey string

func (c contextKey) String() string {
	return "mcpgw_v1 context key " + string(c)
}

const (
	contextKeyMethodDesc = contextKey("methodDesc")
)

func MethodDescContext(ctx context.Context) *MethodDesc {
	return ctx.Value(contextKeyMethodDesc).(*MethodDesc)
}

func NewMethodDescContext(ctx context.Context, methodDesc *MethodDesc) context.Context {
	return context.WithValue(ctx, contextKeyMethodDesc, methodDesc)
}

type RequestMetadata interface {
	Host() string
	RemoteAddr() string
}

func MetadataForRequest(req RequestMetadata, methodFulLName string) metadata.MD {
	// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
	rv := metadata.MD{}

	// emulate grpc-go's behavior of setting these headers.
	// TODO(pquerna): should we do this or preserve the original headers?
	rv.Set("content-type ", "application/grpc+proto")
	rv.Set(":method", "POST")
	rv.Set(":path", methodFulLName)
	rv.Set(":authority", req.Host())
	rv.Set(":scheme", "https")

	return rv
}

func PeerForRequest(req RequestMetadata) *peer.Peer {
	// TODO(pquerna): grpc-server uses a raw conn address here.
	pr := &peer.Peer{
		Addr: strAddr(req.RemoteAddr()),
	}
	pr.AuthInfo = credentials.TLSInfo{
		State:          tls.ConnectionState{},
		CommonAuthInfo: credentials.CommonAuthInfo{SecurityLevel: credentials.PrivacyAndIntegrity},
	}
	return pr
}
