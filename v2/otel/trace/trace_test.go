package trace

import (
	"os"
	"testing"
	"context"
	"time"

	"github.com/rs/zerolog"
)

var logger = zerolog.New(os.Stdout).
					With().
					Timestamp().
					Str("component","go-core").
					Str("package", "http").
					Logger()

func TestCore_Observability(t *testing.T){
	var tracerProvider TracerProvider

	ctx, cancel := context.WithTimeout(	context.Background(), 
										time.Duration( 30 ) * time.Second)
	defer cancel()

	infoTrace := InfoTrace{}
	envTrace := EnvTrace{}

	tracerProvider.NewTracerProvider(ctx, 
									envTrace, 
									infoTrace,
									&logger)
									
	tracerProvider.SpanCtx(ctx,"my-span-test")
}