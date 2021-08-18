package main

import (
	"fmt"
	"os"

	"github.com/thelitdeveloper/imageresizer/pkg/server/http_handlers"
	"go.uber.org/zap"
)

func init() {
	//
	// Init logging
	//
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	//
	// Let's check the presence of required env vars
	// to work with AWS first:
	//
	requiredEnvVars := []string{"AWS_REGION", "AWS_BUCKET",
		"AWS_S3_ENDPOINT", "AWS_ENDPOINT_SCHEME"}
	for _, envar := range requiredEnvVars {
		if val, ok := os.LookupEnv(envar); !ok || val == "" {
			zap.L().Fatal(fmt.Sprintf("%s is missing", envar))
		}
	}
}

func main() {

	//
	// Run the service
	//
	if err := http_handlers.LaunchHttpServer(); err != nil {
		zap.L().Fatal(err.Error())
	}

}
