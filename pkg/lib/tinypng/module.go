package tinypng

import (
	"context"
	"io"

	"go.uber.org/fx"

	"notifications/pkg/lib/config"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Resizer interface {
	Resize(context.Context, io.Reader) (string, error)
	UploadToAws(_ context.Context, url, dir, fileName string, size Size) error
}

type Params struct {
	fx.In

	Logger logger.Logger
	Config config.Config
}

type tinyPng struct {
	logger             logger.Logger
	apiKey             string
	awsAccessKeyID     string
	awsSecretAccessKey string
	region             string
	bucket             string
}

func New(p Params) Resizer {
	return &tinyPng{
		logger:             p.Logger,
		apiKey:             p.Config.GetString("tinyPng.apiKey"),
		awsAccessKeyID:     p.Config.GetString("aws.access.key.id"),
		awsSecretAccessKey: p.Config.GetString("aws.secret.access.key"),
		region:             p.Config.GetString("aws.region"),
		bucket:             p.Config.GetString("fileManager.bucket"),
	}
}
