package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	clientS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/walnuts1018/s3-oauth2-proxy/config"
	"github.com/walnuts1018/s3-oauth2-proxy/domain/model"
	"github.com/walnuts1018/s3-oauth2-proxy/domain/repository"
)

type s3Repository struct {
	client *clientS3.Client
	bucket string
}

func NewS3Repository(cfg aws.Config, s3Config config.S3Config) repository.S3Repository {
	return &s3Repository{
		client: clientS3.NewFromConfig(cfg, func(o *clientS3.Options) {
			o.UsePathStyle = s3Config.UsePathStyle
		}),
		bucket: s3Config.Bucket,
	}
}

func (r *s3Repository) GetObject(ctx context.Context, key string) (*model.S3Object, error) {
	input := &clientS3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	}

	result, err := r.client.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}

	return &model.S3Object{
		Body:          result.Body,
		ContentLength: *result.ContentLength,
		ContentType:   *result.ContentType,
	}, nil
}
