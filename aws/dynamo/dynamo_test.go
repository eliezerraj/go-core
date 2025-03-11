package dynamo

import (
	"time"
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	go_core_aws_config "github.com/eliezerraj/go-core/aws/aws_config"
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

func TestCore_Dynamo(t *testing.T){
	var awsConfig	go_core_aws_config.AwsConfig
	var databaseDynamo	DatabaseDynamo

	aws_config, err := awsConfig.NewAWSConfig(context.Background(), "us-east-2")
	if err != nil {
		t.Errorf("failed to get aws_config : %s", err)
	}

	dynamoDB, err := databaseDynamo.NewDatabaseDynamo(aws_config)
	if err != nil {
		t.Errorf("failed create a NewDatabaseDynamo : %s", err)
	}

	tableName := "user_login_2"
	id := "USER-admin1"
	sk := "USER-admin"

	result, err := dynamoDB.QueryInput(context.Background(), &tableName, id, sk )
	if err != nil {
		t.Errorf("failed QueryInput : %s", err)
	}
	if len(result) == 0 {
		t.Errorf("Not Found")
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
}