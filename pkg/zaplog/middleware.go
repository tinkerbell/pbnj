package zaplog

import (
	"context"
	"reflect"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/rs/xid"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryLogRequestID returns a new unary server interceptors that adds zap.Logger with requestID to the context.
func UnaryLogRequestID(_ *zap.Logger, requestIDKey, requestIDLogKey string) grpc.UnaryServerInterceptor {
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

// UnaryLogBMCIP returns a new unary server interceptors that adds the BMC IP to the logger.
func UnaryLogBMCIP() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (h interface{}, err error) {
		// handle panics from reflect by just moving on. Calls to the task endpoint will always panic,
		// as they don't contain Authn.
		defer func() {
			if r := recover(); r != nil {
				h, err = handler(ctx, req)
			}
		}()
		// get BMC IP
		var bmcIP string
		if a, ok := reflect.ValueOf(req).Elem().FieldByName("Authn").Interface().(*v1.Authn); ok {
			bmcIP = a.GetDirectAuthn().GetHost().GetHost()
		}

		// Get logger from context and add ip to logger
		logger := ctxzap.Extract(ctx).With(zap.String("bmcIP", bmcIP))

		// Add the logger back to the context
		newCtx := ctxzap.ToContext(ctx, logger)

		return handler(newCtx, req)
	}
}
