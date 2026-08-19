package tracing

import(
	"time"
	"log"
	"context"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/contrib/propagators/aws/xray"	
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Struct for tracer information
type InfoTrace struct {
	Name			string `json:"service_name,omitempty"`
	Version			string `json:"service_version,omitempty"`
	ServiceType		string `json:"service_type,omitempty"`
	Env				string `json:"environment,omitempty"`
	Account			string `json:"account,omitempty"`
}

// Struct for enviroment information
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

// Tracer provider object
type TracerProvider struct {
	TracerProvider *sdktrace.TracerProvider
	Tracer			trace.Tracer
}

var requestIdKey = "x-request-id"

// About create a http tracer provider
func NewTracerProvider(	ctx context.Context, 
						envTrace 	EnvTrace, 
						infoTrace 	InfoTrace) *TracerProvider {

	log.Println("Creating new tracer provider")

	var stdout_export sdktrace.SpanExporter
	var err error

	if envTrace.UseStdoutTracerExporter {
		stdout_export, err = stdouttrace.New()
		if err != nil {
			log.Fatalf("Fail create STDOUT trace exporter WARNING !!!: %+v", err)
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
		log.Fatalf("Error create OTEL trace exporter: %+v", err)
	}

	resources, err := buildResources(ctx, infoTrace, envTrace)
	if err != nil {
		log.Fatalf("Error create OTEL resource: %+v", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
		sdktrace.WithResource(resources),
		sdktrace.WithBatcher(stdout_export),
		sdktrace.WithIDGenerator(xray.NewIDGenerator()),
	)
	
	log.Println("OTEL Tracer Provider created successfully")
	
	return &TracerProvider{
		TracerProvider: tracerProvider,
		Tracer: tracerProvider.Tracer("github.com/eliezerraj/go-core/v3/otel/trace"),
	}
}

// About set attributes
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

// About create resourcer (attributes)
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
								 spanName string,
								 spanKind trace.SpanKind) (context.Context, trace.Span) {
	// Get request ID from context using middleware function

	id, ok := ctx.Value(requestIdKey).(string)
	if !ok {
		id = "not-informed"
	}

	ctx, span := t.Tracer.Start(ctx,
							  	spanName,
							  	trace.WithSpanKind(spanKind),
									trace.WithAttributes(
									attribute.String(string(requestIdKey), id),
								),
	)

	return ctx, span
}