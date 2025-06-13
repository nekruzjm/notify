package event

import (
	"context"
	"crypto/sha1"
	"encoding/binary"
	"mime/multipart"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"notifications/internal/db/tx"
	"notifications/internal/repo/event"
	"notifications/internal/repo/rom"
	"notifications/internal/repo/user"
	"notifications/internal/service/admin"
	"notifications/pkg/lib/broker/nats"
	"notifications/pkg/lib/cache"
	"notifications/pkg/lib/config"
	"notifications/pkg/lib/fileman"
	"notifications/pkg/lib/notifier/firebase"
	"notifications/pkg/lib/observer/logger"
	"notifications/pkg/lib/observer/sentry"
	"notifications/pkg/lib/tinypng"
)

var Module = fx.Provide(New)

type Service interface {
	reader
	writer
	runner
	imageManager
}

type reader interface {
	GetEvents(context.Context, Filter) ([]*Event, error)
}

type writer interface {
	Create(context.Context, admin.Admin, *Request) (*Event, error)
	Update(context.Context, admin.Admin, *Request) (*Event, error)
	Delete(ctx context.Context, adminUser admin.Admin, id int) error
}

type runner interface {
	// LoadUsers reads and processes a CSV file containing user references (userID, phone, or person external reference)
	// It divides the records into chunks and concurrently subscribes these users to the given firebase cloud messaging topic
	// The function returns a summary of how many records succeeded or failed, along with details of any failures
	// The function uses worker pool pattern to process the CSV records concurrently
	LoadUsers(context.Context, admin.Admin, int, multipart.File, *multipart.FileHeader) (*Event, error)
	LoadAllUsers(ctx context.Context, a admin.Admin, id int) (*Event, error)
	RunEvent(ctx context.Context, a admin.Admin, id int) (any, error)
	RunJob()
	SubscribeUsers(context.Context, *Event) error
	SubscribeAllUsers(context.Context, *Event) error
	UnsubscribeUsers(context.Context, *Event) error
}

type imageManager interface {
	UploadImage(context.Context, admin.Admin, string, int, multipart.File, *multipart.FileHeader) (*Event, error)
	RemoveImage(ctx context.Context, admin admin.Admin, lang string, id int) (*Event, error)
}

type Params struct {
	fx.In

	Logger      logger.Logger
	Sentry      sentry.Sentry
	Config      config.Config
	Nats        nats.Event
	Cache       cache.Cache
	FcmTopicMan firebase.TopicManager
	FcmSender   firebase.Sender
	FileManager fileman.FileManager
	TinyPng     tinypng.Resizer
	EventRepo   event.Repo
	UserRepo    user.Repo
	RomRepo     rom.Repo
	Transactor  tx.Transactor
}

type service struct {
	logger      logger.Logger
	sentry      sentry.Sentry
	config      config.Config
	nats        nats.Event
	cache       cache.Cache
	fcmTopicMan firebase.TopicManager
	fcmSender   firebase.Sender
	fileManager fileman.FileManager
	tinyPng     tinypng.Resizer
	eventRepo   event.Repo
	userRepo    user.Repo
	romRepo     rom.Repo
	transactor  tx.Transactor
	idGenerator *snowflake.Node

	storageUrl string
	bucket     string
	directory  string
}

func New(p Params) Service {
	hash := sha1.Sum([]byte(p.Config.GetString("podName")))
	podID := int64(binary.LittleEndian.Uint64(hash[:8]) % 1024)
	if podID == 0 {
		podID = 1
	}

	idGenerator, err := snowflake.NewNode(podID)
	if err != nil {
		p.Logger.Error("snowflake.NewNode", zap.Error(err))
		return nil
	}

	return &service{
		logger:      p.Logger,
		sentry:      p.Sentry,
		config:      p.Config,
		nats:        p.Nats,
		cache:       p.Cache,
		fcmTopicMan: p.FcmTopicMan,
		fcmSender:   p.FcmSender,
		fileManager: p.FileManager,
		tinyPng:     p.TinyPng,
		eventRepo:   p.EventRepo,
		userRepo:    p.UserRepo,
		romRepo:     p.RomRepo,
		transactor:  p.Transactor,
		idGenerator: idGenerator,
		storageUrl:  p.Config.GetString("fileManager.storageURL"),
		bucket:      p.Config.GetString("fileManager.bucket"),
		directory:   p.Config.GetString("fileManager.directory"),
	}
}
