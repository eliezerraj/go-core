package parameter_store

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var childLogger = log.With().Str("go-core.aws", "aws_parameter_store").Logger()

type AwsParameterStore struct {
	Client *ssm.Client
}

// About create a new client
func (p *AwsParameterStore) NewAwsParameterStore(configAWS *aws.Config) (*AwsParameterStore, error) {
	childLogger.Debug().Msg("NewAwsParameterStore")

	client := ssm.NewFromConfig(*configAWS)
	return &AwsParameterStore{
		Client: client,
	}, nil
}

// About get a parameter
func (p *AwsParameterStore) GetParameter(ctx context.Context, parameterName string) (*string, error) {
	childLogger.Debug().Msg("GetParameter")

	result, err := p.Client.GetParameter(ctx, 
										&ssm.GetParameterInput{
											Name:	aws.String(parameterName),
											WithDecryption:	aws.Bool(false), // Set to true for SecureString parameters
										})
	if err != nil {
		return nil, err
	}
	return result.Parameter.Value, nil
}