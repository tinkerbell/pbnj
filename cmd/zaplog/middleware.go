package zaplog

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryLogRequestID returns a new unary server interceptors that adds zap.Logger with requestID to the context.
func UnaryLogRequestID(logger *zap.Logger, requestIDKey, requestIDLogKey string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		t := grpc_ctxtags.Extract(ctx)

		data, ok := metadata.FromIncomingContext(ctx)
		if ok {
			reqID := data.Get(requestIDKey)
			if len(reqID) > 0 {
				t.Set(requestIDLogKey, reqID[0])
			} else {
				id := xid.New()
				t.Set(requestIDLogKey, id.String())
				md, _ := metadata.FromIncomingContext(ctx)
				md.Append(requestIDKey, id.String())
				ctx = metadata.NewOutgoingContext(ctx, md)
			}
		}

		newCtx := ctxzap.ToContext(ctx, ctxzap.Extract(ctx))

		resp, err := handler(newCtx, req)
		return resp, err
	}
}
