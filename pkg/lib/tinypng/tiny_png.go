package tinypng

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/imroc/req/v3"
	"go.uber.org/zap"
)

func (t *tinyPng) Resize(_ context.Context, file io.Reader) (string, error) {
	var (
		requestUrl = "https://api.tinify.com/shrink"
		response   Response
	)

	resp, err := req.R().
		SetHeaders(t.getHeaders()).
		SetBody(file).
		SetSuccessResult(&response).
		Post(requestUrl)
	if err != nil {
		t.logger.Error("cannot send request", zap.Error(err))
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		err = errors.New(resp.String())
		t.logger.Error("incorrect status", zap.Error(err))
		return "", err
	}

	return response.Output.Url, nil
}

func (t *tinyPng) UploadToAws(_ context.Context, url, dir, fileName string, size Size) error {
	const (
		_awsService   = "s3"
		_resizeMethod = "fit"
	)
	var request = Request{
		Store: Store{
			Service:            _awsService,
			AwsAccessKeyId:     t.awsAccessKeyID,
			AwsSecretAccessKey: t.awsSecretAccessKey,
			Region:             t.region,
			Path:               t.bucket + "/" + dir + size.Format + fileName,
		},
		Resize: Resize{
			Method: _resizeMethod,
			Width:  size.Width,
			Height: size.Height,
		},
	}

	resp, err := req.R().
		SetHeaders(t.getHeaders()).
		SetBody(request).
		Post(url)
	if err != nil {
		t.logger.Error("cannot send request", zap.Error(err))
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.String())
		t.logger.Error("incorrect status", zap.Error(err))
		return err
	}

	return nil
}
