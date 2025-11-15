package core_trace

import(
	"fmt"
	"context"
	"time"
	"github.com/rs/zerolog"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/sdk/resource"
)

type InfoTrace struct {
	Name			string `json:"service_name,omitempty"`
	Version			string `json:"service_version,omitempty"`
	ServiceType		string `json:"service_type,omitempty"`
	Env				string `json:"enviroment,omitempty"`
	Account			string `json:"account,omitempty"`
}

type EnvTrace struct {
	OtelExportEndpoint			string
	TimeInterval            	int64    `mapstructure:"TimeInterval"`
	TimeAliveIncrementer    	int64    `mapstructure:"RandomTimeAliveIncrementer"`
	TotalHeapSizeUpperBound 	int64    `mapstructure:"RandomTotalHeapSizeUpperBound"`
	ThreadsActiveUpperBound 	int64    `mapstructure:"RandomThreadsActiveUpperBound"`
	CpuUsageUpperBound      	int64    `mapstructure:"RandomCpuUsageUpperBound"`
	SampleAppPorts          	[]string `mapstructure:"SampleAppPorts"`
	AWSCloudWatchLogGroup		[]string `mapstructure:"AWSCloudWatchLogGroup"`
	UseStdoutTracerExporter		bool	 `mapstructure:"UseStdoutTracerExporter"`
	UseOtlpCollector			bool	 `mapstructure:"UseOtlpCollector"`   
}

type TracerProvider struct {
}

// About create a http tracer provider
func (t *TracerProvider) NewTracerProvider(	ctx context.Context, 
											envTrace 	EnvTrace, 
											infoTrace 	InfoTrace,
											appLogger 	*zerolog.Logger) *sdktrace.TracerProvider {
	logger := appLogger.With().
						Str("component", "go-core.otel.trace").
						Logger()

	logger.Debug().
			Str("func","NewTracerProvider").Send()

	var stdout_export sdktrace.SpanExporter
	var err error

	if envTrace.UseStdoutTracerExporter {
		stdout_export, err = stdouttrace.New()
		if err != nil {
			logger.Warn().
					Err(nil).
					Msg("Fail create STDOUT trace exporter WARNING !!!")
		}
	}

	var auth_option otlptracegrpc.Option
	auth_option = otlptracegrpc.WithInsecure()

	exporter, err := otlptrace.New(	ctx,
									otlptracegrpc.NewClient(
										otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
											Enabled:         true,
											InitialInterval: time.Millisecond * 100,
											MaxInterval:     time.Millisecond * 500,
											MaxElapsedTime:  time.Second,
										}),
										auth_option,
										otlptracegrpc.WithEndpoint(envTrace.OtelExportEndpoint),
									),)
	if err != nil {
		logger.Error().
				Err(err).
				Msg("Erro create OTEL trace exporter")
	}

	resources, err := buildResources(ctx, infoTrace, envTrace)
	if err != nil {
		logger.Error().
				Err(err).
				Msg("Erro to build OTEL resource")
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resources),
		sdktrace.WithBatcher(stdout_export),
	)
	return tp
}

func attributes(infoTrace InfoTrace, 
				envTrace EnvTrace) []attribute.KeyValue {

	return []attribute.KeyValue{
		attribute.String("service.name", infoTrace.Name),
		attribute.String("service.version", infoTrace.Version),
		attribute.String("account", infoTrace.Account),
		attribute.String("service.type", infoTrace.ServiceType),
		attribute.String("env", infoTrace.Env),
		semconv.AWSLogGroupNamesKey.StringSlice(envTrace.AWSCloudWatchLogGroup),
		semconv.TelemetrySDKLanguageGo,
	}
}

func buildResources(ctx context.Context, 
					infoTrace InfoTrace, 
					envTrace EnvTrace) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithAttributes(attributes(infoTrace, envTrace)...),
	)
}

// About create a span and return the context
func (t *TracerProvider) SpanCtx(ctx context.Context, 
								spanName string) (context.Context, trace.Span) {
	
	// get tracer id
	trace_id := "not-informed"
	trace_id = fmt.Sprintf("%v",ctx.Value("trace-request-id"))
	tracer := otel.GetTracerProvider().Tracer("go.opentelemetry.io/otel")
	
	ctx, span := tracer.Start(	ctx,
								spanName,
								trace.WithSpanKind(trace.SpanKindConsumer),
								trace.WithAttributes(
									attribute.String("trace-request-id", trace_id),
								),
							)

	return ctx, span
}