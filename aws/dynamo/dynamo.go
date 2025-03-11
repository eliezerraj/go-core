package dynamo

import(
	"context"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

var childLogger = log.With().Str("go-core.aws", "dynamo").Logger()

type DatabaseDynamo struct {
	Client 		*dynamodb.Client
}

func (d* DatabaseDynamo) NewDatabaseDynamo(configAWS *aws.Config) (*DatabaseDynamo){
	childLogger.Debug().Msg("NewDatabaseDynamo")

	client := dynamodb.NewFromConfig(*configAWS)

	return &DatabaseDynamo {
		Client: client,
	}
}

// About query item in Dynamo
func (d* DatabaseDynamo) QueryInput(ctx context.Context, tableName *string ,id string, sk string) ([]map[string]types.AttributeValue, error){
	childLogger.Debug().Msg("QueryInput")

	var keyCond expression.KeyConditionBuilder

	keyCond = expression.KeyAnd(
		expression.Key("ID").Equal(expression.Value(id)),
		expression.Key("SK").BeginsWith(sk),
	)

	expr, err := expression.NewBuilder().
							WithKeyCondition(keyCond).
							Build()
	if err != nil {
		return nil, err
	}

	key := &dynamodb.QueryInput{TableName: tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := d.Client.Query(ctx, key)
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

// About put item in Dynamo
func (d* DatabaseDynamo) PutItem(ctx context.Context, tableName *string, data any) (error){
	childLogger.Debug().Msg("PutItem")

	item, err := attributevalue.MarshalMap(data)
	if err != nil {
		return err
	}

	putInput := &dynamodb.PutItemInput{
        TableName: tableName,
        Item:      item,
    }

	_, err = d.Client.PutItem(ctx, putInput)
    if err != nil {
		return err
    }
	
	return nil
}