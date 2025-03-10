package parameter_store

import (
	"context"
	"testing"

	go_core_aws_config "github.com/eliezerraj/go-core/aws/aws_config"
)

func TestCore_ParameterStore(t *testing.T){
	var awsConfig	go_core_aws_config.AwsConfig
	var awsParameterStore	AwsParameterStore

	aws_config, err := awsConfig.NewAWSConfig(context.Background(), "us-east-2")
	if err != nil {
		t.Errorf("failed to get aws_config : %s", err)
	}

	workerParameterStore, err := awsParameterStore.NewAwsParameterStore(aws_config)
	if err != nil {
		t.Errorf("failed create a NewParameterStore : %s", err)
	}

	parameterName := "key-secret"
	parameter, err := workerParameterStore.GetParameter(context.Background(), parameterName)
	if err != nil {
		t.Errorf("failed GetParameter : %s", err)
	}

	t.Logf("=====>>>>> parameter: %v", parameter)
}