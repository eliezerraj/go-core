package aws_config

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var childLogger = log.With().Str("component","go-core").Str("package", "aws.aws_config").Logger()

type AwsConfig struct {
}

func (a *AwsConfig) NewAWSConfig(ctx context.Context, awsRegion string) (*aws.Config, error) {
	childLogger.Debug().Str("func","NewAWSConfig").Send()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}