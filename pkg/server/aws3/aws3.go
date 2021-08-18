package aws3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/thelitdeveloper/imageresizer/pkg/server/entities"
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

// NewS3Client populate an AWS Config and create a client,
// uploader and downloader
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
func (s3c *s3Client) DownloadSingleImage(ctx context.Context, nkey string) (image.Image, error) {
	var err error

	// 1. Pre-allocating memory by requesting download's content length
	headObject, err := s3c.getHeadObject(ctx, nkey)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, int(headObject.ContentLength))
	buffer := manager.NewWriteAtBuffer(buf)

	// 2. Do download
	_, err = s3c.downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(nkey),
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
func (s3c *s3Client) UploadSingleImage(ctx context.Context, imgBytes *bytes.Reader, props *entities.ImageProperties) (string, error) {
	var err error

	// 1. Detect and set the right Content-Type
	contentType := detectContentType(props)

	// 2. Create a new key from properties
	newKey := createNewKeyFromProps(props)
	if newKey == "" {
		return "", errors.New("uploadSingleImage() Couldn't create a key")
	}

	// 3. Upload it
	_, err = s3c.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("AWS_BUCKET")),
		Key:         aws.String(newKey),
		Body:        imgBytes,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		zap.L().Error(err.Error())
		return "", err
	}

	return newKey, nil
}

//
// Helper functions
//

// getHeadObject is collecting infos before download
func (s3c *s3Client) getHeadObject(ctx context.Context, nkey string) (*s3.HeadObjectOutput, error) {

	input := &s3.HeadObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(nkey),
	}

	output, err := s3c.client.HeadObject(ctx, input)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return output, nil
}

// Create a new key including width and height
func createNewKeyFromProps(props *entities.ImageProperties) string {
	// 1. Build file name
	filenameParts := strings.Split(props.Filename, ".")
	if len(filenameParts) < 2 {
		return ""
	}
	if props.Convert {
		filenameParts[1] = "jpg"
		// Strip convert req from filename
		cnvReqPatterns := []string{"-jpg", "-jpeg", "-JPEG", "-JPG"}
		for _, pat := range cnvReqPatterns {
			if strings.Contains(filenameParts[0], pat) {
				filenameParts[0] = strings.TrimRight(filenameParts[0], pat)
			}
		}
	}
	newFilename := fmt.Sprintf("%s.%s", filenameParts[0],
		filenameParts[1])

	// 2. Build prefix
	var newKey string
	if props.Prefix == "" {
		newKey = newFilename
	} else {
		newKey = fmt.Sprintf("%s/%s", props.Prefix, newFilename)
	}

	return newKey
}

// detectContentType
func detectContentType(props *entities.ImageProperties) string {
	cntType := "image/jpeg"

	if !props.Convert && strings.ToLower(strings.Split(props.Filename, ".")[1]) == "png" {
		return "image/png"
	}

	return cntType
}
