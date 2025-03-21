package bucket_s3

import (
	"context"
	"io"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var childLogger = log.With().Str("component","go-core").Str("package", "aws.bucket_s3").Logger()

type AwsBucketS3 struct {
	Client *s3.Client
}

// About create a new client
func (b *AwsBucketS3) NewAwsS3Bucket(configAWS *aws.Config) *AwsBucketS3 {
	childLogger.Debug().Str("func","NewAwsS3Bucket").Send()

	client := s3.NewFromConfig(*configAWS)
	
	return &AwsBucketS3{
		Client: client,
	}
}

// About get a object
func (b *AwsBucketS3) GetObject(ctx context.Context, 	
								bucketNameKey 	string,
								filePath 		string,
								fileKey 		string) (*string, error) {
	childLogger.Debug().Str("func","GetObject").Send()

	getObjectInput := &s3.GetObjectInput{
						Bucket: aws.String(bucketNameKey+filePath),
						Key:    aws.String(fileKey),
	}

	getObjectOutput, err := b.Client.GetObject(ctx, getObjectInput)
	if err != nil {
		return nil, err
	}
	defer getObjectOutput.Body.Close()

	bodyBytes, err := io.ReadAll(getObjectOutput.Body)
	if err != nil {
		return nil, err
	}

	res := string(bodyBytes)
	return &res, nil
}