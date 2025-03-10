package secret_manager

import (
	"context"
	
	"github.com/rs/zerolog/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var childLogger = log.With().Str("go-core.aws", "aws_secret_manager").Logger()

type AwsSecretManager struct {
	Client *secretsmanager.Client
}

// About create a new client
func (p *AwsSecretManager) NewAwsSecretManager(configAWS *aws.Config) (*AwsSecretManager, error) {
	childLogger.Debug().Msg("NewAwsSecretManager")
		
	client := secretsmanager.NewFromConfig(*configAWS)

	return &AwsSecretManager{
		Client: client,
	}, nil
}

// About get a secret
func (p *AwsSecretManager) GetSecret(ctx context.Context, secretName string) (*string, error) {
	childLogger.Debug().Msg("GetSecret")

	result, err := p.Client.GetSecretValue(ctx, 
										&secretsmanager.GetSecretValueInput{
											SecretId:		aws.String(secretName),
											VersionStage:	aws.String("AWSCURRENT"),
										})
	if err != nil {
		return nil, err
	}
	return result.SecretString, nil
}