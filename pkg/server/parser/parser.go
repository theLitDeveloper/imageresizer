package parser

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/thelitdeveloper/imageresizer/pkg/server/entities"
)

// Parse extracts and validates params; also creates a fallback URI
func Parse(qryString string) (*entities.ImageProperties, string) {

	//
	// Validate URL path
	//
	if valid := validate(qryString); !valid {
		return nil, ""
	}

	//
	// Splitting URL path into segments
	//
	segments := strings.Split(qryString, "/")
	count := len(segments)

	//
	// Check, if convert is requested; default = false
	//
	doConvert := extractConvertRequest(segments[count-1])

	//
	// Extract width and height -
	// at least one of them must be greater
	// than zero
	//
	dim, err := extractDimensions(segments[count-1])
	if err != nil {
		return nil, ""
	} else if dim[0] == 0 && dim[1] == 0 {
		return nil, ""
	}

	//
	// Extract prefix
	//
	prefix := extractPrefix(segments[0 : count-1])

	//
	// Create fallback URI in case of an error
	//
	fallbackKey, err := CreateFallbackKey(prefix, segments[count-1])
	if err != nil {
		return nil, ""
	}
	redirectURI := fmt.Sprintf("%s://%s.%s.%s.amazonaws.com/%s",
		os.Getenv("AWS_ENDPOINT_SCHEME"),
		os.Getenv("AWS_BUCKET"),
		os.Getenv("AWS_S3_ENDPOINT"),
		os.Getenv("AWS_REGION"),
		fallbackKey)

	//
	// 	Return the parsed ImageProperties
	//
	return &entities.ImageProperties{
		Prefix:   prefix,
		Filename: segments[count-1],
		Width:    dim[0],
		Height:   dim[1],
		Convert:  doConvert,
	}, redirectURI
}

// validate checks for a malformed or incomplete URL path
func validate(p string) bool {
	valid, _ := regexp.MatchString(`([a-z0-9-_A-Z]*/?)*([a-z0-9A-Z-_]{3,})-\d+x\d+(-jpe?g|-JPE?G)?\.(jpe?g|JPE?G|png|PNG)$`, p)
	return valid
}

// extractSizeParams returns width and height as integer
func extractDimensions(paramSegment string) ([]int, error) {

	fnParts := strings.Split(paramSegment, ".")

	// Remove jpg, jpeg, JPG or JPEG if any
	result := removeConvertRequestParam(fnParts[0])

	// Separate width and height from the rest of the string
	tmp := strings.Split(result, "-")
	rawWidthAndHeight := tmp[len(tmp)-1]
	if !strings.Contains(rawWidthAndHeight, "x") {
		return nil, errors.New("couldn't extract width and height")
	}

	splitted := strings.Split(rawWidthAndHeight, "x")
	width, _ := strconv.Atoi(splitted[0])
	height, _ := strconv.Atoi(splitted[1])

	return []int{width, height}, nil
}

// extractPath returns a path fragment
func extractPrefix(pathSegments []string) string {

	if count := len(pathSegments); count == 1 {
		return pathSegments[0]
	} else if count > 1 {
		return strings.Join(pathSegments, "/")
	}

	return ""
}

// extractConvertRequest checks, if there is a request for converting
// into another file format (jpg -> png)
func extractConvertRequest(fname string) bool {
	f := strings.Split(fname, ".")
	if f[1] != "png" && f[1] != "PNG" {
		return false
	}

	pttrns := []string{"-jpg", "-jpeg", "-JPEG", "-JPG"}
	for _, pat := range pttrns {
		if strings.Contains(f[0], pat) {
			return true
		}
	}

	return false
}

// CreateFallbackKey strips all params from query string
func CreateFallbackKey(prefix string, fname string) (string, error) {

	// Extract file extension
	ext := strings.Split(fname, ".")[1]

	// Remove all params from filename
	strippedCnvParam := strings.Split(removeConvertRequestParam(fname), ".")[0]

	// Remove width and height params
	tmp := strings.Split(strippedCnvParam, "-")
	wh := tmp[len(tmp)-1]

	cleanedFilename := strings.TrimRight(strippedCnvParam, fmt.Sprintf("-%s", wh))

	fallbackKey := fmt.Sprintf("%s.%s", cleanedFilename, ext)
	if prefix != "" {
		fallbackKey = fmt.Sprintf("%s/%s", prefix, fallbackKey)
	}
	return fallbackKey, nil
}

//
// Helper funcs
//
func removeConvertRequestParam(fname string) string {
	fn := strings.Split(fname, ".")
	cnvReqPatterns := []string{"-jpg", "-jpeg", "-JPEG", "-JPG"}
	for _, pat := range cnvReqPatterns {
		if strings.Contains(fn[0], pat) {
			fn[0] = strings.TrimRight(fn[0], pat)
		}
	}
	return strings.Join(fn, ".")
}
