package aws_secret_manager

import (
	"context"
	"testing"
	"encoding/json"

	go_core_aws_config "github.com/eliezerraj/go-core/aws/aws_config"
)

func TestCore_SecretManager(t *testing.T){
	var awsConfig	go_core_aws_config.AwsConfig
	var awsClientSecretManager	AwsClientSecretManager

	aws_config, err := awsConfig.NewAWSConfig(context.Background(), "us-east-2")
	if err != nil {
		t.Errorf("failed to get aws_config : %s", err)
	}

	clientSecretManager, err := awsClientSecretManager.NewAwsClientSecretManager(aws_config)
	if err != nil {
		t.Errorf("failed create a NewAwsClientSecretManager : %s", err)
	}

	secretName := "key-jwt-auth"
	secret, err := clientSecretManager.GetSecret(context.Background(), secretName)
	if err != nil {
		t.Errorf("failed GetSecret : %s", err)
	}

	var secretData map[string]string
	if err := json.Unmarshal([]byte(*secret), &secretData); err != nil {
		t.Errorf("failed unmarshal secret : %s", err)
	}
	
	t.Logf("=====>>>>> secretData: %v", secretData)
}