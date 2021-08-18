package imagemanipulation

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"

	"github.com/disintegration/imaging"
	"github.com/thelitdeveloper/imageresizer/pkg/server/entities"
	"go.uber.org/zap"
)

// ResizeConvertEncode resizes, converts (if requested), and encode given image
func ResizeConvertEncode(img image.Image, props *entities.ImageProperties) (*bytes.Reader, error) {

	// Do resizing
	resizedImg := imaging.Resize(img, props.Width, props.Height, imaging.Lanczos)

	// Evaluate format
	format, err := imaging.FormatFromFilename(props.Filename)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	//	If conversion is requested, convert now to jpg
	if props.Convert {
		format = imaging.JPEG
		destRect := image.NewNRGBA(resizedImg.Bounds())
		draw.Draw(destRect, destRect.Bounds(), &image.Uniform{C: color.White}, image.Point{
			X: 0,
			Y: 0,
		}, draw.Src)
		draw.Draw(destRect, destRect.Bounds(), resizedImg, resizedImg.Bounds().Min, draw.Over)
		resizedImg = destRect
	}

	// Last but not least: encode
	buffer := new(bytes.Buffer)
	if err = imaging.Encode(buffer, resizedImg, format); err != nil {
		zap.L().Fatal(err.Error())
	}

	return bytes.NewReader(buffer.Bytes()), nil
}
