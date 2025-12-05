package s3

import (
	"os"
	"context"
	"testing"
	"github.com/rs/zerolog"
)

var logger = zerolog.New(os.Stdout).
						With().
						Str("component", "testgocore.aws.s3").
						Logger()

func TestCore_BucketS3(t *testing.T){

	var awsRegion	= "us-east-2"

	workerBucketS3, err := NewAwsBucketS3(context.Background(),
											awsRegion,
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