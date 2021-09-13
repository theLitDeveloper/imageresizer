package aws3

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"os"

	"gitlab.simplyadmin.de/sdCoreengines/imageresizer/pkg/service/models"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/disintegration/imaging"
)

type s3Client struct {
	client     *s3.Client
	uploader   *manager.Uploader
	downloader *manager.Downloader
}

// NewS3Client populates an AWS Config and creates a client,
// inits uploader and downloader
func NewS3Client(ctx context.Context) *s3Client {

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	client := s3.NewFromConfig(cfg)

	downloader := manager.NewDownloader(client)
	uploader := manager.NewUploader(client)

	return &s3Client{
		client:     client,
		uploader:   uploader,
		downloader: downloader,
	}
}

// DownloadSingleImage download original image
func (s3c *s3Client) DownloadSingleImage(ctx context.Context, resourcePath string) (image.Image, error) {
	var err error

	// 1. Pre-allocating memory by requesting download's content length
	headObject, err := s3c.getHeadObject(ctx, resourcePath)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, int(headObject.ContentLength))
	buffer := manager.NewWriteAtBuffer(buf)

	// 2. Do download
	_, err = s3c.downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(resourcePath),
	})
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	// 3. Decode image
	rawImage := bytes.NewReader(buffer.Bytes())
	img, err := imaging.Decode(rawImage)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return img, nil
}

// UploadSingleImage upload resized image
func (s3c *s3Client) UploadSingleImage(ctx context.Context, imgBytes *bytes.Reader,
	params string, props *models.ImageProperties) error {

	var err error

	// 1. Create a new resource path
	var resourcePath string
	if props.ResourcePath == "" {
		resourcePath = props.ResourceName
	} else {
		resourcePath = fmt.Sprintf("%s/%s", props.ResourcePath, props.ResourceName)
	}
	newResourcePath := fmt.Sprintf("%s/%s/%s", props.ClientUUID, params, resourcePath)

	// 2. Upload it to S3 bucket
	_, err = s3c.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("AWS_BUCKET")),
		Key:         aws.String(newResourcePath),
		Body:        imgBytes,
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	return nil
}

// getHeadObject is collecting infos before download
func (s3c *s3Client) getHeadObject(ctx context.Context, resourcePath string) (*s3.HeadObjectOutput, error) {

	input := &s3.HeadObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(resourcePath),
	}

	output, err := s3c.client.HeadObject(ctx, input)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return output, nil
}
