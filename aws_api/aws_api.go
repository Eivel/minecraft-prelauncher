package aws_api

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSClient struct {
	Session       *session.Session
	DefaultBucket *string
}

func (awsClient *AWSClient) Initialize(accessKey string, secretKey string, host string) {
	awsClient.Session = session.New(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(host),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(true),
	})
	awsClient.DefaultBucket = aws.String("mods")
}

func (awsClient *AWSClient) DownloadFile(filePath string, bucketKey string) error {
	bucket := awsClient.DefaultBucket
	key := aws.String(bucketKey)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(awsClient.Session)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: bucket,
			Key:    key,
		})
	if err != nil {
		return err
	}
	fmt.Println("Downloaded file", file.Name(), numBytes, "bytes")
	return nil
}
