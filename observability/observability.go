package observability

import(
	"context"
	"time"
	"github.com/rs/zerolog/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/sdk/resource"
)

var childLogger = log.With().Str("go-core", "observability").Logger()

type InfoTrace struct {
	PodName				string `json:"pod_name"`
	PodVersion			string `json:"pod_version"`
	ServiceType			string `json:"service_type"`
	Env					string `json:"enviroment,omitempty"`
	AccountID			string `json:"account_id,omitempty"`
}

type ConfigOTEL struct {
	OtelExportEndpoint		string
	TimeInterval            int64    `mapstructure:"TimeInterval"`
	TimeAliveIncrementer    int64    `mapstructure:"RandomTimeAliveIncrementer"`
	TotalHeapSizeUpperBound int64    `mapstructure:"RandomTotalHeapSizeUpperBound"`
	ThreadsActiveUpperBound int64    `mapstructure:"RandomThreadsActiveUpperBound"`
	CpuUsageUpperBound      int64    `mapstructure:"RandomCpuUsageUpperBound"`
	SampleAppPorts          []string `mapstructure:"SampleAppPorts"`
}

type TracerProvider struct {
}

func attributes(ctx context.Context, infoTrace *InfoTrace) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("service.name", infoTrace.PodName),
		attribute.String("service.version", infoTrace.PodVersion),
		attribute.String("account", infoTrace.AccountID),
		attribute.String("service.type", infoTrace.ServiceType),
		attribute.String("env", infoTrace.Env),
		semconv.TelemetrySDKLanguageGo,
	}
}

func buildResources(ctx context.Context, infoTrace *InfoTrace) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithAttributes(attributes(ctx, infoTrace)...),
	)
}

func (t *TracerProvider) NewTracerProvider(	ctx context.Context, 
											configOTEL *ConfigOTEL, 
											infoTrace 	*InfoTrace) *sdktrace.TracerProvider {
	log.Debug().Msg("NewTracerProvider")

	var authOption otlptracegrpc.Option
	authOption = otlptracegrpc.WithInsecure()

	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
				Enabled:         true,
				InitialInterval: time.Millisecond * 100,
				MaxInterval:     time.Millisecond * 500,
				MaxElapsedTime:  time.Second,
			}),
			authOption,
			otlptracegrpc.WithEndpoint(configOTEL.OtelExportEndpoint),
		),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create OTEL trace exporter")
	}

	resources, err := buildResources(ctx, infoTrace)
	if err != nil {
		log.Error().Err(err).Msg("failed to load OTEL resource")
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resources),
	)
	return tp
}

func (t *TracerProvider) Event(span trace.Span, attributeSpan string) {
	span.AddEvent("Executing SQL query", trace.WithAttributes(attribute.String("db.statement", attributeSpan)))
}

func (t *TracerProvider) Span(ctx context.Context, spanName string) trace.Span {
	cID, rID := "unknown", "unknown"
	/*if id, ok := logger.ClientUUID(ctx); ok {
		cID = id
	}
	if id, ok := logger.RequestUUID(ctx); ok {
		rID = id
	}*/

	tracer := otel.GetTracerProvider().Tracer("go.opentelemetry.io/otel")
	
	_, span := tracer.Start(
							ctx,
							spanName,
							trace.WithSpanKind(trace.SpanKindConsumer),
							trace.WithAttributes(
								attribute.String("id", cID),
								attribute.String("request_id", rID)),
	)

	return span
}
