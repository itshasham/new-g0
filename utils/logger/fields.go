package logger

import "context"

type Fields map[string]interface{}

const (
	FieldError         = "error"
	FieldCorrelationID = "correlation_id"
	FieldMethod        = "method"
	FieldRequest       = "request"
	FieldResponse      = "response"
)

type correlationIDKey struct{}

func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDKey{}, correlationID)
}

func extractCorrelationID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if correlationID, ok := ctx.Value(correlationIDKey{}).(string); ok {
		return correlationID
	}
	return ""
}
