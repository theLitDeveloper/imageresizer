package handlers

import (
	"fmt"
	"net/http"
	"os"

	"gitlab.simplyadmin.de/sdCoreengines/imageresizer/pkg/service/aws3"
	"gitlab.simplyadmin.de/sdCoreengines/imageresizer/pkg/service/imagemanipulation"
	"gitlab.simplyadmin.de/sdCoreengines/imageresizer/pkg/service/parser"
	"go.uber.org/zap"
)

// ManipulateImageHandler is the main entrypoint
func ManipulateImageHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	//
	// 	Validate request method (GET)
	//
	if r.Method != "GET" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//
	// Do the query string exist and isn't empty?
	//
	var ref string
	if ref = r.URL.Query().Get("ref"); "" == ref {
		badRequest(w)
		return
	}

	//
	//	Parse query string into properties
	//  and also create a redirect URI in case of an error
	//
	props, redirectURI, resourcePath, params := parser.Parse(ref)
	if props == nil || redirectURI == "" {
		badRequest(w)
		return
	}

	//
	// 	Download the image from S3 bucket
	//
	s3c := aws3.NewS3Client(ctx)
	img, err := s3c.DownloadSingleImage(ctx, resourcePath)
	if err != nil {
		zap.L().Error(err.Error())
		http.Redirect(w, r, redirectURI, http.StatusTemporaryRedirect)
		return
	}

	//
	// 	Change the image with given parameter
	//
	changedImage, err := imagemanipulation.ManipulateImage(img, props)
	if err != nil {
		http.Redirect(w, r, redirectURI, http.StatusTemporaryRedirect)
		return
	}

	//
	// 	Upload changed image to S3 bucket
	// 	build a new redirect URL for the user-agent
	//
	err = s3c.UploadSingleImage(ctx, changedImage, params, props)
	if err != nil {
		zap.L().Error(err.Error())
		http.Redirect(w, r, redirectURI, http.StatusTemporaryRedirect)
		return
	}

	//
	// 	Redirect the user-agent to the newly
	// 	created resource permanently
	//
	if props.ResourcePath == "" {
		resourcePath = props.ResourceName
	} else {
		resourcePath = fmt.Sprintf("%s/%s", props.ResourcePath, props.ResourceName)
	}

	redirectURI = fmt.Sprintf("https://%s/%s/%s/%s",
		os.Getenv("REDIRECT_HOST"),
		props.ClientUUID,
		params,
		resourcePath)

	http.Redirect(w, r, redirectURI, http.StatusMovedPermanently)
}

// HealthcheckHandler endpoint for infrastructure monitoring
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

// badRequest returns status code 400
func badRequest(w http.ResponseWriter) {
	http.Error(w, "400 Bad request", http.StatusBadRequest)
}
