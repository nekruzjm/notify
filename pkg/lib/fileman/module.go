package fileman

import (
	"context"
	"io"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"notifications/pkg/lib/config"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type FileManager interface {
	Upload(uploadFile io.Reader, bucket *string, dir, fileName string) error
	Remove(bucket *string, dir, fileName string) error
}

type Params struct {
	fx.In

	Config config.Config
	Logger logger.Logger
}

type file struct {
	config config.Config
	logger logger.Logger
	awsS3  *s3.Client
}

func New(p Params) FileManager {
	var f = &file{
		config: p.Config,
		logger: p.Logger,
	}

	cfg, err := awscfg.LoadDefaultConfig(
		context.Background(),
		awscfg.WithRegion("eu-central-1"),
		awscfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				p.Config.GetString("aws.access.key.id"),
				p.Config.GetString("aws.secret.access.key"), "")))
	if err != nil {
		p.Logger.Info("failed to load AWS config", zap.Error(err))
	}

	f.awsS3 = s3.NewFromConfig(cfg)

	return f
}
