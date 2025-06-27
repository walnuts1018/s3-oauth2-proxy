package config

type S3Config struct {
	Bucket       string `env:"S3_BUCKET,required"`
	UsePathStyle bool   `env:"S3_USE_PATH_STYLE" envDefault:"false"`
}
