package metric

import(
	"context"

	"github.com/rs/zerolog"

	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type InfoMetric struct {
	Name			string `json:"service_name,omitempty"`
	Version			string `json:"service_version,omitempty"`
	ServiceType		string `json:"service_type,omitempty"`
	Env				string `json:"enviroment,omitempty"`
	Account			string `json:"account,omitempty"`
}

// About initialize MeterProvider with Prometheus exporter
func NewMeterProvider(	ctx context.Context, 
						infoMetric InfoMetric,
						appLogger *zerolog.Logger) (*sdkmetric.MeterProvider, error) {

	logger := appLogger.With().
						Str("component", "go-core.v2.otel.metric").
						Logger()

	logger.Info().
			Str("func","NewMeterProvider").Send()

	// 1. Configurar o Recurso OTel
	res, err := resource.New(ctx,
							resource.WithSchemaURL(semconv.SchemaURL),
							resource.WithAttributes(
								semconv.ServiceNameKey.String(infoMetric.Name),
								attribute.String("version", infoMetric.Version),
							),
	)
	if err != nil {
		logger.Error().
				Err(err).Send()
		return nil, err
	}

	// 2. Criar o Prometheus Exporter
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	// 3. Criar o MeterProvider, usando o Prometheus Exporter como Reader.
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(exporter),
	)

	return meterProvider, nil
}