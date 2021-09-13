package service

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.simplyadmin.de/sdCoreengines/imageresizer/pkg/service/handlers"

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
	// to work with AWS services first:
	//
	requiredEnvVars := []string{"REDIRECT_HOST", "AWS_BUCKET", "AWS_REGION"}
	for _, envar := range requiredEnvVars {
		if val, ok := os.LookupEnv(envar); !ok || val == "" {
			zap.L().Fatal(fmt.Sprintf("%s is missing", envar))
		}
	}
}

// registerHandlers registering all routes
func registerHandlers() {

	http.HandleFunc("/health", handlers.HealthcheckHandler)
	http.HandleFunc("/do", handlers.ManipulateImageHandler)
	http.Handle("/metrics", promhttp.Handler())
}

// Run registering handlers and spin up HTTP server
func Run() {

	registerHandlers()

	port := os.Getenv("PORT")
	if port == "" {
		port = "4321"
	}

	http.ListenAndServe(":"+port, nil)
}
