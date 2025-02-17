package observability

import (
	"testing"
	"context"
	"time"
)

func TestCore_Observability(t *testing.T){
	var tracerProvider TracerProvider

	ctx, cancel := context.WithTimeout(	context.Background(), 
										time.Duration( 30 ) * time.Second)
	defer cancel()

	configOTEL := ConfigOTEL{}
	infoTrace := InfoTrace{}

	tracerProvider.NewTracerProvider(ctx, &configOTEL, &infoTrace)
	tracerProvider.Span(ctx,"my-span-test")
}