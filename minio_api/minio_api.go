package minio_api

import (
	"fmt"

	"github.com/minio/minio-go"
)

func InitializeClient(accessKey string, secretKey string, host string) (*minio.Client, error) {
	client, err := minio.New(host, accessKey, secretKey, true)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func DownloadFile(client *minio.Client, bucket string, filePath string, key string) error {
	err := client.FGetObject(bucket, key, filePath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	fmt.Println("Downloaded file", key)
	return nil
}
