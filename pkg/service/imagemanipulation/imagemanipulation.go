package imagemanipulation

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"strings"

	"gitlab.simplyadmin.de/sdCoreengines/imageresizer/pkg/service/models"

	"github.com/disintegration/imaging"
	"go.uber.org/zap"
)

// ManipulateImage batches all requested steps and returns a changed image
func ManipulateImage(img image.Image, settings *models.ImageProperties) (*bytes.Reader, error) {
	var err error
	var manipulatedImage *image.NRGBA

	// detecting image format
	format, err := imaging.FormatFromFilename(settings.ResourceName)
	if err != nil {
		zap.L().Error(err.Error())
	}

	// resizing image (always)
	manipulatedImage = resize(img, settings.Params.Width, settings.Params.Height)

	// cropping (if requested)
	if settings.Params.Crop {
		manipulatedImage = crop(manipulatedImage, settings.Params.Width, settings.Params.Height)
	}

	// convert to grayscale (if requested)
	if settings.Params.ConvToGrayscale {
		manipulatedImage = convertToGrayscale(manipulatedImage)
	}

	// convert to jpg when png is detected
	if format == imaging.PNG {

		manipulatedImage = convertToJpg(manipulatedImage)
		if err != nil {
			zap.L().Error(err.Error())
		}
		format = imaging.JPEG

		// change resource name (extension)
		rn := strings.Split(settings.ResourceName, ".")
		rn[1] = "jpg"
		settings.ResourceName = strings.Join(rn, ".")
	}

	// Encode and return the changed image
	buffer := new(bytes.Buffer)
	if err = imaging.Encode(buffer, manipulatedImage, format); err != nil {
		zap.L().Fatal(err.Error())
	}

	return bytes.NewReader(buffer.Bytes()), nil
}

// resize the image to width and height
func resize(img image.Image, width int, height int) *image.NRGBA {
	return imaging.Resize(img, width, height, imaging.Lanczos)
}

// convertToGrayscale
func convertToGrayscale(img *image.NRGBA) *image.NRGBA {
	return imaging.Grayscale(img)
}

// crop is centering the image and crops to width and height then
func crop(img *image.NRGBA, width int, height int) *image.NRGBA {
	return imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
}

// convertToJpg converts png images to jpg format
func convertToJpg(img *image.NRGBA) *image.NRGBA {
	destRect := image.NewNRGBA(img.Bounds())
	draw.Draw(destRect, destRect.Bounds(), &image.Uniform{C: color.White}, image.Point{
		X: 0,
		Y: 0,
	}, draw.Src)
	draw.Draw(destRect, destRect.Bounds(), img, img.Bounds().Min, draw.Over)
	return destRect
}
