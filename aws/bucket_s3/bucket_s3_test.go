package bucket_s3

import (
	"context"
	"testing"

	go_core_aws_config "github.com/eliezerraj/go-core/aws/aws_config"
)

func TestCore_BucketS3(t *testing.T){

	var awsConfig	go_core_aws_config.AwsConfig
	var awsBucketS3	AwsBucketS3

	aws_config, err := awsConfig.NewAWSConfig(context.Background(), "us-east-2")
	if err != nil {
		t.Errorf("failed to get aws_config : %s", err)
	}

	workerBucketS3 := awsBucketS3.NewAwsS3Bucket(aws_config)

	bucketNameKey := "eliezerraj-908671954593-mtls-truststore"
	filePath := "/"
	fileKey := "server-private.key"

	obj, err := workerBucketS3.GetObject(context.Background(), 
										bucketNameKey,
										filePath,
										fileKey )
	if err != nil {
		t.Errorf("failed GetObject : %s", err)
	}

	t.Logf("=====>>>>> obj: %v", obj)

	bucketNameKey = "eliezerraj-908671954593-mtls-truststore"
	filePath = "/"
	fileKey = "text.txt"
	file := []byte("my test")

	err = workerBucketS3.PutObject(context.Background(), 
										bucketNameKey,
										filePath,
										fileKey,
										file )
	if err != nil {
		t.Errorf("failed PutObject : %s", err)
	}
}