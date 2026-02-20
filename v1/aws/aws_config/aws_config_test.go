package aws_config

import (
	"context"
	"testing"
)

func TestCore_AWSConfig(t *testing.T){
	var awsConfig	AwsConfig
	
	_, err := awsConfig.NewAWSConfig(context.Background(), "us-east-2")
	if err != nil {
		t.Errorf("failed to get aws_config : %s", err)
	}
}