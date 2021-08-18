package http_handlers

import (
	"net/http"
)

// LaunchHttpServer wraps handler registration and the launch of the http server
func LaunchHttpServer() error {

	//
	// 	Just two endpoints
	//
	http.HandleFunc("/health", HealthcheckHandler)
	http.HandleFunc("/resize", ResizeHandler)

	//
	// Spin up the http server
	//
	return http.ListenAndServe(":4321", nil)
}
