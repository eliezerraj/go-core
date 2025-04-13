package observability

import(
	"fmt"
	"context"
	"time"
	"github.com/rs/zerolog/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/metadata"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/sdk/resource"
)

var childLogger = log.With().Str("component","go-core").Str("package", "observability").Logger()

type InfoTrace struct {
	PodName				string `json:"pod_name"`
	PodVersion			string `json:"pod_version"`
	ServiceType			string `json:"service_type"`
	Env					string `json:"enviroment,omitempty"`
	AccountID			string `json:"account_id,omitempty"`
}

type ConfigOTEL struct {
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

func attributes(ctx context.Context, infoTrace *InfoTrace, configOTEL *ConfigOTEL) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("service.name", infoTrace.PodName),
		attribute.String("service.version", infoTrace.PodVersion),
		attribute.String("account", infoTrace.AccountID),
		attribute.String("service.type", infoTrace.ServiceType),
		attribute.String("env", infoTrace.Env),
		semconv.AWSLogGroupNamesKey.StringSlice(configOTEL.AWSCloudWatchLogGroup),
		semconv.TelemetrySDKLanguageGo,
	}
}

func buildResources(ctx context.Context, infoTrace *InfoTrace, configOTEL *ConfigOTEL) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithAttributes(attributes(ctx, infoTrace, configOTEL)...),
	)
}

// About create a http tracer provider
func (t *TracerProvider) NewTracerProvider(	ctx context.Context, 
											configOTEL *ConfigOTEL, 
											infoTrace 	*InfoTrace) *sdktrace.TracerProvider {
	childLogger.Debug().Str("func","NewTracerProvider").Send()
	
	// turn on/off the tracer
	if !configOTEL.UseStdoutTracerExporter && !configOTEL.UseOtlpCollector {
		return nil
	}

	var stdout_export sdktrace.SpanExporter
	var err error
	if configOTEL.UseStdoutTracerExporter {
		stdout_export, err = stdouttrace.New()
		if err != nil {
			childLogger.Error().Err(err).Msg("error creating stdout logging exporter")
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
										otlptracegrpc.WithEndpoint(configOTEL.OtelExportEndpoint),
									),)
	if err != nil {
		childLogger.Error().Err(err).Msg("failed to create OTEL trace exporter")
	}

	resources, err := buildResources(ctx, infoTrace, configOTEL)
	if err != nil {
		childLogger.Error().Err(err).Msg("failed to build OTEL resource")
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

// About create a event
func (t *TracerProvider) Event(span trace.Span, attributeSpan string) {
	span.AddEvent("Executing SQL query", trace.WithAttributes(attribute.String("db.statement", attributeSpan)))
}

// About create a span
func (t *TracerProvider) Span(ctx context.Context, spanName string) trace.Span {
	
	// get tracer id
	trace_id := "not-informed"
	trace_id = fmt.Sprintf("%v",ctx.Value("trace-request-id"))

	tracer := otel.GetTracerProvider().Tracer("go.opentelemetry.io/otel")
	
	_, span := tracer.Start(
							ctx,
							spanName,
							trace.WithSpanKind(trace.SpanKindConsumer),
							trace.WithAttributes(
								attribute.String("trace-request-id", trace_id)),
	)

	return span
}

// About create a span and return the context
func (t *TracerProvider) SpanCtx(ctx context.Context, spanName string) (context.Context, trace.Span) {
	
	// get tracer id
	trace_id := "not-informed"
	trace_id = fmt.Sprintf("%v",ctx.Value("trace-request-id"))

	tracer := otel.GetTracerProvider().Tracer("go.opentelemetry.io/otel")
	
	ctx, span := tracer.Start(
							ctx,
							spanName,
							trace.WithSpanKind(trace.SpanKindConsumer),
							trace.WithAttributes(
								attribute.String("trace-request-id", trace_id)),
	)

	return ctx, span
}

// For Grpc 
type MetadataCarrier struct {
	metadata.MD
}

func (mc MetadataCarrier) Get(key string) string {
	values := mc.MD.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (mc MetadataCarrier) Set(key, value string) {
	mc.MD.Set(key, value)
}

func (mc MetadataCarrier) Keys() []string {
	keys := make([]string, 0, len(mc.MD))
	for k := range mc.MD {
		keys = append(keys, k)
	}
	return keys
}