package config

import (
	"errors"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssSync struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
}

func (config OssSync) GetBucket() (*oss.Bucket, error) {
	if config.Endpoint == "" || config.AccessKeyId == "" || config.AccessKeySecret == "" || config.BucketName == "" {
		return nil, errors.New("OssSync is invalid")
	}
	client, err := oss.New(config.Endpoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}
