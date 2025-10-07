package feature_flag

import (
	"context"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appconfigdata"
)

var childLogger = 	log.With().
					Str("component","go-core").
					Str("package", "aws.config.feature-flag").
					Logger()

type AwsConfigFeatureFlag struct {
	Client *appconfigdata.Client
}

// About create a new client
func (p *AwsConfigFeatureFlag) NewAwsConfigFeatureFlag(configAWS *aws.Config) (*AwsConfigFeatureFlag) {
	childLogger.Debug().Str("func","NewAwsConfigFeatureFlag").Send()

	client := appconfigdata.NewFromConfig(*configAWS)

	return &AwsConfigFeatureFlag{
		Client: client,
	}
}

// About create a session input
func (p *AwsConfigFeatureFlag) StartConfigurationSession(	ctx context.Context,
															applicationID string,
															configurationProfileID string,
															environmentID string) (*appconfigdata.GetLatestConfigurationOutput, error) {
	childLogger.Debug().Str("func","StartConfigurationSessionInput").Send()

	featureFlagConfig := &appconfigdata.StartConfigurationSessionInput{
							ApplicationIdentifier: &applicationID,
							ConfigurationProfileIdentifier: &configurationProfileID,
							EnvironmentIdentifier: &environmentID,
							RequiredMinimumPollIntervalInSeconds: func(i int32) *int32 { return &i }(15), // Poll every 15 seconds
	}

	startSessionOutput, err := p.Client.StartConfigurationSession(ctx, featureFlagConfig)
	if err != nil {
		return nil, err
	}

	token := startSessionOutput.InitialConfigurationToken

	getConfOutput, err := p.Client.GetLatestConfiguration(ctx, &appconfigdata.GetLatestConfigurationInput{
		ConfigurationToken: token,
	})
	if err != nil {
		return nil, err
	}

	return getConfOutput, nil
}