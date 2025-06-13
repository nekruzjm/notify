package fileman

import (
	"context"
	"io"
	"path/filepath"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

func (f *file) Upload(file io.Reader, bucket *string, dir, fileName string) error {
	_, err := f.awsS3.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: bucket,
		Key:    aws.String(dir + fileName),
		Body:   file,
	})
	if err != nil {
		f.logger.Error("err from f.awsUploader.Upload", zap.Error(err))
		return err
	}

	return nil
}

func (f *file) Remove(bucket *string, dir, fileName string) error {
	_, err := f.awsS3.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: bucket,
		Key:    aws.String(dir + fileName),
	})
	if err != nil {
		f.logger.Error("err from f.awsS3.DeleteObject", zap.Error(err))
		return err
	}

	return nil
}

func GetFileExt(filename string) string {
	return strings.ToLower(strings.Trim(filepath.Ext(filename), "."))
}

func IsImg(ext string) bool {
	return slices.Contains([]string{PNG, JPEG, JPG}, ext)
}
