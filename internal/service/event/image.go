package event

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"slices"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"notifications/internal/api/resp"
	"notifications/internal/api/transport/broker/stream"
	"notifications/internal/api/transport/broker/subject"
	"notifications/internal/lib/language"
	"notifications/internal/repo/repomodel"
	"notifications/internal/service/admin"
	"notifications/pkg/lib/fileman"
	"notifications/pkg/lib/tinypng"
	"notifications/pkg/util/strset"
	"notifications/pkg/util/uid"
)

func (s *service) UploadImage(ctx context.Context, a admin.Admin, lang string, id int, file multipart.File, fileHeader *multipart.FileHeader) (*Event, error) {
	if !slices.Contains(language.GetAll(), lang) {
		s.logger.Error("invalid language", zap.String("lang", lang))
		return nil, resp.Wrap(resp.ErrBadRequest, "invalid language")
	}

	tx := s.transactor.New()
	err := tx.Begin(ctx)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("failed to begin transaction", zap.Error(err))
		return nil, err
	}

	defer func() {
		if err != nil {
			if errX := tx.Rollback(ctx); errX != nil {
				err = errors.Join(err, fmt.Errorf("errX: %w", errX))
			}
			return
		}
		err = tx.Commit(ctx)
		if err != nil {
			err = errors.Join(err, errors.New("tx commit failed"))
		}
	}()

	selectedEvent, err := tx.EventRepo().GetByID(ctx, id)
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("failed to get events by id", zap.Error(err), zap.Int("id", id))
		}
		return nil, err
	}

	var (
		oldEvent = *selectedEvent
		fileExt  = fileman.GetFileExt(fileHeader.Filename)
	)

	if !fileman.IsImg(fileExt) {
		s.logger.Warning("image extension is not valid", zap.String("file ext", fileExt))
		return nil, resp.Wrap(resp.ErrBadRequest, "image extension is not valid")
	}

	url, err := s.tinyPng.Resize(ctx, file)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("failed to resize image", zap.Error(err))
		return nil, err
	}

	var (
		sizes   = tinypng.Sizes()
		eg      errgroup.Group
		builder strings.Builder
	)

	builder.WriteString(uid.ShortUID())
	builder.WriteString(_underscoreDelim)
	builder.WriteString(lang)
	builder.WriteString(_dotDelim)
	builder.WriteString(fileExt)

	fileName := builder.String()

	for _, size := range sizes {
		eg.Go(func() error {
			if err = s.tinyPng.UploadToAws(ctx, url, s.directory, fileName, size); err != nil {
				s.sentry.CaptureException(err)
				s.logger.Error("failed to upload image", zap.Error(err))
				return err
			}
			return nil
		})
	}
	if err = eg.Wait(); err != nil {
		return nil, resp.Wrap(resp.ErrInternalErr, err.Error())
	}

	var img = selectedEvent.Image.Get(lang)
	if !strset.IsEmpty(img) {
		for _, size := range sizes {
			eg.Go(func() error {
				if err = s.fileManager.Remove(&s.bucket, s.directory+size.Format, img); err != nil {
					s.sentry.CaptureException(err)
					s.logger.Error("failed to remove image", zap.Error(err))
					return err
				}
				return nil
			})
		}
		if err = eg.Wait(); err != nil {
			return nil, resp.Wrap(resp.ErrInternalErr, err.Error())
		}
	}

	selectedEvent.Image.Set(lang, fileName)
	updatedEvent, err := tx.EventRepo().UpdateImage(ctx, id, selectedEvent.Image)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("failed to update image", zap.Error(err))
		return nil, err
	}

	err = s.nats.Publish(stream.Audit, subject.AuditAdd, admin.Audit{
		AdminId:   a.ID,
		IpAddress: a.IP,
		EventName: admin.UploadImageEvent,
		OldData:   oldEvent,
		NewData:   *selectedEvent,
		CreatedAt: time.Now(),
	})
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("failed to publish audit event", zap.Error(err))
	}

	var event = new(Event)
	event.toService(updatedEvent)
	s.setImgURL(event)

	return event, nil
}

func (s *service) RemoveImage(ctx context.Context, a admin.Admin, lang string, id int) (*Event, error) {
	if !slices.Contains(language.GetAll(), lang) {
		s.logger.Error("invalid language", zap.String("lang", lang))
		return nil, resp.Wrap(resp.ErrBadRequest, "invalid language")
	}

	tx := s.transactor.New()
	err := tx.Begin(ctx)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("failed to begin transaction", zap.Error(err))
		return nil, err
	}

	defer func() {
		if err != nil {
			if errX := tx.Rollback(ctx); errX != nil {
				err = errors.Join(err, fmt.Errorf("errX: %w", errX))
			}
			return
		}
		err = tx.Commit(ctx)
		if err != nil {
			err = errors.Join(err, errors.New("tx commit failed"))
		}
	}()

	selectedEvent, err := tx.EventRepo().GetByID(ctx, id)
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("failed to get events by id", zap.Error(err), zap.Int("id", id))
		}
		return nil, err
	}

	var (
		oldEvent = *selectedEvent
		fileName = selectedEvent.Image.Get(lang)
		sizes    = tinypng.Sizes()
		eg       errgroup.Group
	)

	for _, size := range sizes {
		eg.Go(func() error {
			if err = s.fileManager.Remove(&s.bucket, s.directory+size.Format, fileName); err != nil {
				s.sentry.CaptureException(err)
				s.logger.Error("failed to remove image", zap.Error(err))
				return err
			}
			return nil
		})
	}
	if err = eg.Wait(); err != nil {
		return nil, resp.Wrap(resp.ErrInternalErr, err.Error())
	}

	selectedEvent.Image.Set(lang, _empty)
	updatedEvent, err := tx.EventRepo().UpdateImage(ctx, id, selectedEvent.Image)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("failed to update image", zap.Error(err))
		return nil, err
	}

	var event = new(Event)
	event.toService(updatedEvent)
	s.setImgURL(event)

	err = s.nats.Publish(stream.Audit, subject.AuditAdd, admin.Audit{
		AdminId:   a.ID,
		IpAddress: a.IP,
		EventName: admin.RemoveImageEvent,
		OldData:   oldEvent,
		NewData:   *selectedEvent,
		CreatedAt: time.Now(),
	})
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("failed to publish audit event", zap.Error(err))
	}

	return event, nil
}

func (s *service) setImgURL(e *Event) {
	images := e.Image.GetAllWithLang()
	for i := 0; i < len(images); i++ {
		if !strset.IsEmpty(images[i].Val) {
			e.Image.Set(images[i].Key, s.storageUrl+s.directory+"1x/"+images[i].Val)
		}
	}
}
