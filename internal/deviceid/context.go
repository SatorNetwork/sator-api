package deviceid

import (
	"context"
	"fmt"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

type contextKey string

var (
	deviceIDHeaderKey  string     = "Device-ID"
	deviceIDContextKey contextKey = "DeviceID"
)

// ToContext moves a Device-ID from request header to context.
func ToContext() httptransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		return context.WithValue(ctx, deviceIDContextKey, r.Header.Get(deviceIDHeaderKey))
	}
}

// FromContext gets device id from context
func FromContext(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(deviceIDContextKey))
}
