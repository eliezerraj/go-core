package feature_flag

import (
	"context"
	"testing"
	"encoding/json"

	go_core_aws_config "github.com/eliezerraj/go-core/aws/aws_config"
)

type FeatureFlag struct {
	Feature01 struct {
		Enabled bool   	`json:"enabled"`
		Value string 	`json:"value"`
	} `json:"feature_01"`
	Feature02 struct {
		Value string `json:"value"`
	} `json:"feature_02"`
}

func TestCore_FeatureFlag(t *testing.T){

	ctx := context.Background()

	var awsConfig				go_core_aws_config.AwsConfig
	var awsConfigFeatureFlag	AwsConfigFeatureFlag

	aws_config, err := awsConfig.NewAWSConfig(context.Background(), "us-east-2")
	if err != nil {
		t.Errorf("failed to get aws_config : %s", err)
	}

	client_feature_flag := awsConfigFeatureFlag.NewAwsConfigFeatureFlag(aws_config)

	applicationID := "i8udatk" 			//go-account-fg
	configurationProfileID := "4ows5ui" // go-account-fg-profile
	environmentID := "dev"

	feature_flag_output, err := client_feature_flag.StartConfigurationSession(ctx,
																			applicationID, 
																			configurationProfileID, 
																			environmentID)
	if err != nil {
		t.Errorf("failed to start configuration session, %v", err)
	}

	var flag FeatureFlag
	err = json.Unmarshal(feature_flag_output.Configuration, &flag)
	if err != nil {
		t.Errorf("failed to unmarshal feature flags, %v", err)
	}

	t.Logf("feature_flag_output.Configuration: %v", string(feature_flag_output.Configuration))
}