package trace

import (
	"os"
	"testing"
	"context"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

var logger = zerolog.New(os.Stdout).
					With().
					Timestamp().
					Str("component","go-core").
					Str("package", "http").
					Logger()

// go test -v -run "TestCore_Observability"
func TestCore_Observability(t *testing.T){
	ctx, cancel := context.WithTimeout(	context.Background(), 
										time.Duration( 30 ) * time.Second)
	defer cancel()

	infoTrace := InfoTrace{}
	envTrace := EnvTrace{}

	test_tracerProvider := NewTracerProvider(	ctx, 
											envTrace, 
											infoTrace,
											&logger)
											
	test_tracerProvider.SpanCtx(ctx,"my-span-test",trace.SpanKindInternal)
}