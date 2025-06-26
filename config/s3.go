package config

type S3Config struct {
	Bucket string `env:"S3_BUCKET,required"`
}
