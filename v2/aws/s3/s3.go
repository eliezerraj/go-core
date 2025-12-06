package s3

import (
	"context"
	"io"
	"bytes"

	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsBucketS3 struct {
	Client *s3.Client
	logger 	*zerolog.Logger
}

// About create a new client
func NewAwsBucketS3(awsConfig *aws.Config,
					appLogger *zerolog.Logger) (*AwsBucketS3, error) {

	logger := appLogger.With().
						Str("component", "go-core.v2.aws.s3").
						Logger()
	logger.Debug().
			Str("func","NewAwsBucketS3").Send()

	client := s3.NewFromConfig(*awsConfig)

	return &AwsBucketS3{
		Client: client,
		logger: &logger,
	}, nil
}

// About get a object
func (b *AwsBucketS3) GetObject(ctx context.Context, 	
								bucketName 	string,
								filePath 	string,
								fileKey 	string) (*string, error) {
	b.logger.Debug().
			Str("func","GetObject").Send()

	getObjectInput := &s3.GetObjectInput{
						Bucket: aws.String(bucketName+filePath),
						Key:    aws.String(fileKey),
	}

	getObjectOutput, err := b.Client.GetObject(ctx, getObjectInput)
	if err != nil {
		b.logger.Error().
				 Err(err).Send()
		return nil, err
	}
	defer getObjectOutput.Body.Close()

	bodyBytes, err := io.ReadAll(getObjectOutput.Body)
	if err != nil {
		b.logger.Error().
				 Err(err).Send()
		return nil, err
	}

	res := string(bodyBytes)
	return &res, nil
}

// About put file
func (b *AwsBucketS3) PutObject(ctx context.Context,
								bucketName 	string,
								filePath 	string,	 	
								fileKey string,
								file 	[]byte) error {
	b.logger.Debug().
			Str("func","PutObject").Send()

	putObjectInput := &s3.PutObjectInput{
						Bucket: aws.String(bucketName+filePath),
						Key:    aws.String(fileKey),
						Body: 	bytes.NewReader(file),
	}

	_, err := b.Client.PutObject(ctx, putObjectInput)
	if err != nil {
		b.logger.Error().
				 Err(err).Send()
		return err
	}

	return nil
}