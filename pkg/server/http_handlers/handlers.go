package http_handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/thelitdeveloper/imageresizer/pkg/server/aws3"
	"github.com/thelitdeveloper/imageresizer/pkg/server/imagemanipulation"
	"github.com/thelitdeveloper/imageresizer/pkg/server/parser"
	"go.uber.org/zap"
)

// ResizeHandler is the main entrypoint
func ResizeHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	ctx := r.Context()

	//
	// 	First things first: validate request method and
	// 	double check provided query string, cause we don't
	// 	waste time nor resources
	//
	if r.Method != "GET" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//
	// Do the query string exist and is it not empty?
	//
	if key := r.URL.Query().Get("key"); "" == key {
		badRequest(w)
		return
	}

	//
	//	Validation of the parameters in the query:
	// 	1. Path, filename, file extension
	//	2. Width and height
	//  and also create a fallback URI in case of an error
	//
	props, redirectURI := parser.Parse(r.URL.Query().Get("key"))
	if props == nil || redirectURI == "" {
		badRequest(w)
		return
	}

	//
	// 	Ok, let's get our hands dirty:
	// 	create a new AWS client
	//
	s3c := aws3.NewS3Client(ctx)

	//
	//	Now we have to strip all params to have "naked" key
	// 	for downloading the original image
	//
	nakedKey, err := parser.CreateFallbackKey(props.Prefix, props.Filename)
	if err != nil {
		zap.L().Error(err.Error())
		http.Redirect(w, r, redirectURI, http.StatusTemporaryRedirect)
		return
	}

	img, err := s3c.DownloadSingleImage(ctx, nakedKey)
	if err != nil {
		zap.L().Error(err.Error())
		http.Redirect(w, r, redirectURI, http.StatusTemporaryRedirect)
		return
	}

	//
	// 	Resize and encode the downloaded image
	//
	tmpImage, err := imagemanipulation.ResizeConvertEncode(img, props)
	if err != nil {
		http.Redirect(w, r, redirectURI, http.StatusTemporaryRedirect)
		return
	}

	//
	// 	Upload resized image and build a new key for
	// 	redirect the user-agent
	//
	newKey, err := s3c.UploadSingleImage(ctx, tmpImage, props)
	if err != nil {
		zap.L().Error(err.Error())
		http.Redirect(w, r, redirectURI, http.StatusTemporaryRedirect)
		return
	}
	redirectURI = fmt.Sprintf("%s://%s.%s.%s.amazonaws.com/%s",
		os.Getenv("AWS_ENDPOINT_SCHEME"),
		os.Getenv("AWS_BUCKET"),
		os.Getenv("AWS_S3_ENDPOINT"),
		os.Getenv("AWS_REGION"),
		newKey)

	//
	// 	If all went well, we redirect the user-agent to the newly
	// 	created resource permanently
	//
	http.Redirect(w, r, redirectURI, http.StatusMovedPermanently)
}

// HealthcheckHandler acts as endpoint for infrastructure monitoring
func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("200 OK"))
	if err != nil {
		return
	}
}

// badRequest
func badRequest(w http.ResponseWriter) {
	http.Error(w, "400 Bad request", http.StatusBadRequest)
}
