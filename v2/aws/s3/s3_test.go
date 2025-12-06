package s3

import (
	"os"
	"context"
	"testing"
	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/config"
)

var logger = zerolog.New(os.Stdout).
						With().
						Str("component", "testgocore.aws.s3").
						Logger()

func TestCore_BucketS3(t *testing.T){

	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-2"))
	if err != nil {
		t.Errorf("failed to get aws_config : %s", err)
	}

	workerBucketS3, err := NewAwsBucketS3(&awsCfg,
											&logger)
	if err != nil {
		t.Errorf("failed QueryInput : %s", err)
	}

	bucketNameKey := "docktech-eliezer-908671954593-truststore-mtls"
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