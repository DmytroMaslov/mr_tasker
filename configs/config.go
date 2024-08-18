package configs

import (
	"os"
)

type Config struct {
	Aws AwsCredentials
}

type AwsCredentials struct {
	Key    string
	Secret string
	Region string
}

func GetConfig() (*Config, error) {
	return &Config{
		Aws: AwsCredentials{
			Key:    os.Getenv("AwsAccessKey"),
			Secret: os.Getenv("AwsAccessSecret"),
			Region: os.Getenv("AwsRegion"),
		},
	}, nil
}
