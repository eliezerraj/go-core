package core_metric

import(
	"github.com/rs/zerolog"
)

type InfoTrace struct {
	Name			string `json:"service_name,omitempty"`
	Version			string `json:"service_version,omitempty"`
	ServiceType		string `json:"service_type,omitempty"`
	Env				string `json:"enviroment,omitempty"`
	Account			string `json:"account,omitempty"`
}

// About initialize MeterProvider with Prometheus exporter
func NewMeterProvider(	ctx context.Context, 
						infoTrace InfoTrace,
						appLogger *zerolog.Logger) (*sdkmetric.MeterProvider, error) {
	logger := appLogger.With().
						Str("component", "go-core.otel.metric").
						Logger()

	logger.Info().
			Str("func","NewMeterProvider").Send()

	// 1. Configurar o Recurso OTel
	res, err := resource.New(ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(infoTrace.Name),
			attribute.String("version", infoTrace.Version),
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
		logger.Error().
				Err(err).Send()
		return nil, err
	}

	// 3. Criar o MeterProvider, usando o Prometheus Exporter como Reader.
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(exporter),
	)

	return provider, nil
}