package logging

import (
	"context"
	"path"
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/rs/xid"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ctxLogr struct{}

var ctxMarkerKey = &ctxLogr{}

// UnaryServerInterceptor returns a new unary server interceptors that adds logr.Logger to the context.
func UnaryServerInterceptor(logger logr.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		l := logger.WithValues("grpc.method", path.Base(info.FullMethod), "grpc.service", strings.TrimPrefix(path.Dir(info.FullMethod), "/"))
		ctx = context.WithValue(ctx, ctxMarkerKey, l)

		return handler(ctx, req)
	}
}

// UnaryLogRequestID returns a new unary server interceptors that adds logr.Logger with requestID to the context if a requestID doesnt exist.
func UnaryLogRequestID(requestIDKey, requestIDLogKey string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var requestID string
		data, ok := metadata.FromIncomingContext(ctx)
		if ok {
			reqID := data.Get(requestIDKey)
			if len(reqID) > 0 {
				requestID = reqID[0]
			} else {
				id := xid.New()
				requestID = id.String()
				md, _ := metadata.FromIncomingContext(ctx)
				md.Append(requestIDKey, id.String())
				ctx = metadata.NewOutgoingContext(ctx, md)
			}
		}
		logger := ExtractLogr(ctx).WithValues(requestIDLogKey, requestID)
		ctx = context.WithValue(ctx, ctxMarkerKey, logger)

		return handler(ctx, req)
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
		logger := ExtractLogr(ctx).WithValues("bmcIP", bmcIP)
		ctx = context.WithValue(ctx, ctxMarkerKey, logger)

		return handler(ctx, req)
	}
}

// ExtractLogr extracts the logr.Logger from the context and adds the grpc tags to the logger.
func ExtractLogr(ctx context.Context) logr.Logger {
	l, ok := ctx.Value(ctxMarkerKey).(logr.Logger)
	if !ok || l.GetSink() == nil {
		return logr.Discard()
	}

	return l.WithValues(grpcTags(ctx)...)
}

// grpcTags converts grpc tags to key/value pairs stored in a slice.
func grpcTags(ctx context.Context) []interface{} {
	var fields []interface{}
	tags := grpc_ctxtags.Extract(ctx)
	for k, v := range tags.Values() {
		fields = append(fields, k, v)
	}

	return fields
}
