package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_resizeHandler(t *testing.T) {
	tests := []struct {
		name   string
		method string
		url    string
		want   int
	}{
		{
			name:   "Method not allowed",
			method: "POST",
			url:    "http://localhost:8080/do?ref=simplytest/c984c70e-a9f9-4cf8-b738-a5467b4dd462/w_500,h_500/blue_marble.jpg",
			want:   405,
		},
		{
			name:   "Bad request (query string not set nor empty)",
			method: "GET",
			url:    "http://localhost:8080/do",
			want:   400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Errorf(err.Error())
			}
			rec := httptest.NewRecorder()
			ManipulateImageHandler(rec, req)
			resp := rec.Result()
			if resp.StatusCode != tt.want {
				t.Errorf("ManipulateImageHandler() want: %v, got: %v", tt.want, resp.StatusCode)
			}
		})
	}
}

func Test_healthcheckHandler(t *testing.T) {
	tests := []struct {
		name   string
		method string
		url    string
		want   int
	}{
		{
			name:   "Usual health check request",
			method: "GET",
			url:    "http://localhost:8080/health",
			want:   200,
		},
		{
			name:   "Method not allowed",
			method: "POST",
			url:    "http://localhost:8080/health",
			want:   405,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Errorf(err.Error())
			}
			rec := httptest.NewRecorder()
			HealthcheckHandler(rec, req)
			resp := rec.Result()
			if resp.StatusCode != tt.want {
				t.Errorf("healthcheckHandler() want: %v, got: %v", tt.want, resp.StatusCode)
			}
		})
	}
}
