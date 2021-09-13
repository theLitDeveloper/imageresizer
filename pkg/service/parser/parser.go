package parser

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gitlab.simplyadmin.de/sdCoreengines/imageresizer/pkg/service/models"
	"go.uber.org/zap"
)

// Parse extracts and validates params; also creates a fallback URI
func Parse(qryString string) (*models.ImageProperties, string, string, string) {

	//
	// Validate query string first
	//
	if valid := validate(qryString); !valid {
		return nil, "", "", ""
	}

	//
	// Parse properties
	//
	var props models.ImageProperties

	querySegments := strings.Split(qryString, "/")

	// Parse params
	var params string
	idx := -1
	for i, seg := range querySegments {
		if strings.Contains(seg, "w_") {
			idx = i
		}
	}

	if idx == -1 {
		return nil, "", "", ""
	} else {
		params = querySegments[idx]
		err := parseParams(querySegments[idx], &props)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, "", "", ""
		}
	}

	// Parse props and retrieve a resource path for creating
	// a redirect URI later on
	resourcePath := parseProps(querySegments, idx, &props)

	//
	// Create redirect URI in case of an error
	//
	redirectURI := fmt.Sprintf("https://%s/%s",
		os.Getenv("REDIRECT_HOST"),
		resourcePath)

	return &props, redirectURI, resourcePath, params
}

// validate checks for a malformed or incomplete URL
func validate(p string) bool {
	valid, _ := regexp.MatchString(`^([a-z0-9\-_+A-Z]+/)((c_fit|e_grayscale|h_\d{1,4}|w_\d{1,4}),(c_fit|e_grayscale|h_\d{1,4}|w_\d{1,4})(,(c_fit|e_grayscale|h_\d{1,4}|w_\d{1,4}))?(,(c_fit|e_grayscale|h_\d{1,4}|w_\d{1,4}))?)([a-z0-9\-_+A-Z/]+)?/([a-zA-Z0-9_\-+]+).(jpeg|JPEG|jpg|JPG|png|PNG)$`, p)
	return valid
}

// parseParams from query string
func parseParams(paramsSegment string, props *models.ImageProperties) error {

	params := strings.Split(paramsSegment, ",")

	for _, param := range params {

		// Detecting crop request
		if param == "c_fit" {
			props.Params.Crop = true
		}

		// Detecting convert to grayscale
		if param == "e_grayscale" {
			props.Params.ConvToGrayscale = true
		}

		// Retrieve width
		if strings.Contains(param, "w_") {
			props.Params.Width, _ = strconv.Atoi(strings.SplitAfter(param, "w_")[1])
		}

		// Retrieve height
		if strings.Contains(param, "h_") {
			props.Params.Height, _ = strconv.Atoi(strings.SplitAfter(param, "h_")[1])
		}
	}

	if props.Params.Width == 0 && props.Params.Height == 0 {
		return errors.New("Providing 0 for width and height isn't allowed")
	}

	return nil
}

// parseProps parses all props from query string except params and returns a resource path
func parseProps(querySegments []string, ignore int, props *models.ImageProperties) string {

	var resourcePath []string

	count := len(querySegments)
	for idx, seg := range querySegments {
		if idx != ignore {

			// Creating resource path
			resourcePath = append(resourcePath, seg)

			// Setting client UUID
			if idx == 0 {
				props.ClientUUID = seg
			}

			// Setting filename
			if idx == count-1 {
				props.ResourceName = seg
			}
		}
	}

	if count > 2 {
		var path []string
		for i, p := range querySegments {
			if i != ignore && i > 0 && i < count-1 {
				path = append(path, p)
			}
		}
		props.ResourcePath = strings.Join(path, "/")
	}

	return strings.Join(resourcePath, "/")
}
