package dynamoDB

import(
	"context"
	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

type DatabaseDynamoDB struct {
	Client 		*dynamodb.Client
	logger 		*zerolog.Logger
}

func NewDatabaseDynamo( ctx context.Context,
						awsRegion string,
						appLogger *zerolog.Logger) (*DatabaseDynamoDB, error) {
	logger := appLogger.With().
						Str("component", "go-core.v2.aws.dynamoDB").
						Logger()
	logger.Debug().
			Str("func","NewDatabaseDynamo").Send()

	configAWS, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		logger.Error().
			   Err(err).Send()
		return nil, err
	}

	client := dynamodb.NewFromConfig(configAWS)

	return &DatabaseDynamoDB {
		Client: client,
		logger: &logger,
	}, nil
}

// About query item in Dynamo
func (d* DatabaseDynamoDB) QueryInput(ctx context.Context, 
									 tableName *string,
									 id string, 
									 sk string) ([]map[string]types.AttributeValue, error){
	d.logger.Debug().
			 Str("func","QueryInput").Send()

	var keyCond expression.KeyConditionBuilder

	keyCond = expression.KeyAnd(
		expression.Key("ID").Equal(expression.Value(id)),
		expression.Key("SK").BeginsWith(sk),
	)

	expr, err := expression.NewBuilder().
							WithKeyCondition(keyCond).
							Build()
	if err != nil {
		d.logger.Error().
				 Err(err).Send()
		return nil, err
	}

	key := &dynamodb.QueryInput{TableName: tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := d.Client.Query(ctx, key)
	if err != nil {
		d.logger.Error().
				 Err(err).Send()
		return nil, err
	}

	return result.Items, nil
}

// About put item in Dynamo
func (d* DatabaseDynamoDB) PutItem(ctx context.Context, 
									tableName *string, 
									data any) (error){
	d.logger.Debug().
			 Str("func","PutItem").Send()

	item, err := attributevalue.MarshalMap(data)
	if err != nil {
		d.logger.Error().
				 Err(err).Send()
		return err
	}

	putInput := &dynamodb.PutItemInput{
        TableName: tableName,
        Item:      item,
    }

	_, err = d.Client.PutItem(ctx, putInput)
    if err != nil {
		d.logger.Error().
				 Err(err).Send()
		return err
    }
	
	return nil
}