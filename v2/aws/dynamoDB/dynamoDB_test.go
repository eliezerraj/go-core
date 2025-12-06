package dynamoDB

import (
	"os"
	"time"
	"context"
	"testing"
	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

type Credential struct {
	ID				string	`json:"ID"`
	SK				string	`json:"SK"`
	User			string	`json:"user,omitempty"`
	Password		string	`json:"password,omitempty"`
	Token			string 	`json:"token,omitempty"`
	UsagePlan		string 	`json:"usage_plan,omitempty"`
	ApiKey			string 	`json:"apikey,omitempty"`
	Updated_at  	time.Time 	`json:"updated_at,omitempty"`
}

var logger = zerolog.New(os.Stdout).
						With().
						Str("component", "testgocore.aws.dynamoDB").
						Logger()

func TestGoCore_Dynamo(t *testing.T){

	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-2"))
	if err != nil {
		t.Errorf("failed to get aws_config : %s", err)
	}

	dynamoDB, err := NewDatabaseDynamo(&awsCfg,
										&logger)
	if err != nil {
		t.Errorf("failed QueryInput : %s", err)
	}

	tableName := "user_login_2"
	id := "USER-admin"
	sk := "USER-admin"

	result, err := dynamoDB.QueryInput(context.Background(), 
										&tableName, 
										id, 
										sk )
	if err != nil {
		t.Errorf("failed QueryInput : %s", err)
	}
	if len(result) == 0 {
		t.Errorf("Data Not Found")
		return
	}

	credential := []Credential{}
	err = attributevalue.UnmarshalListOfMaps(result, &credential)
    if err != nil {
		t.Errorf("failed UnmarshalListOfMaps : %s", err)
    }
	t.Logf("=====>>>>> credential[0]: %v", credential[0])

	tableName = "user_login_2"
	credential_input := Credential{}
	credential_input.ID = "USER-teste-01"
	credential_input.SK = "USER-teste-01"
	credential_input.Updated_at = time.Now()

	err = dynamoDB.PutItem(context.Background(), &tableName, credential_input )
	if err != nil {
		t.Errorf("failed PutItem : %s", err)
	}

	tableName = "user_login_2"
	id = "USER-admin"
	sk = "USER"

	result, err = dynamoDB.QueryInput(context.Background(), &tableName, id, sk )
	if err != nil {
		t.Errorf("failed QueryInput : %s", err)
	}
	if len(result) == 0 {
		t.Errorf("Data Not Found")
		return
	}
}