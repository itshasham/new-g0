package logger

import (
	"context"
	"testing"
)

func TestLogger(t *testing.T) {
	ctx := context.Background()
	ctx = WithCorrelationID(ctx, "test-correlation-id")

	fields := Fields{
		"test_field": "test_value",
		FieldError:   "test error message",
	}

	Info(ctx, "Test info message", fields)
	Debug(ctx, "Test debug message", fields)
	Warn(ctx, "Test warning message", fields)
	Error(ctx, "Test error message", fields)
}

func TestCorrelationID(t *testing.T) {
	ctx := context.Background()

	// Test empty correlation ID
	correlationID := extractCorrelationID(ctx)
	if correlationID != "" {
		t.Errorf("Expected empty correlation ID, got %s", correlationID)
	}

	// Test with correlation ID
	ctx = WithCorrelationID(ctx, "test-id")
	correlationID = extractCorrelationID(ctx)
	if correlationID != "test-id" {
		t.Errorf("Expected 'test-id', got %s", correlationID)
	}
}
