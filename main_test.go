package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestLoginHandler(t *testing.T) {
	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	recorder := httptest.NewRecorder()

	// Create a mock HTTP request.
	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() failed: %v", err)
	}

	// Call the loginHandler function, passing in our mock HTTP request.
	loginHandler(recorder, req)

	// Check the HTTP status code of the response.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("loginHandler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}
}

func TestSecretHandler(t *testing.T) {
	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	recorder := httptest.NewRecorder()

	// Create a mock HTTP request.
	req, err := http.NewRequest("GET", "/secret", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() failed: %v", err)
	}

	// Define a mock secret.
	mockSecret := "test"

	// Call the secretHandler function, passing in a mock HTTP request.
	secretHandler(recorder, req, mockSecret)

	// Check the HTTP status code of the response.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("secretHandler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check the response body.
	expected := `{"secret":"test"}`
	if strings.TrimSuffix(recorder.Body.String(), "\n") != expected {
		t.Errorf("secretHandler returned unexpected body: got %v, want %v", recorder.Body.String(), expected)
	}
}

func TestMetricsHandler(t *testing.T) {
	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	recorder := httptest.NewRecorder()

	// Create a mock HTTP request.
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() failed: %v", err)
	}

	// Call the promhttp.Handler() function, passing in our mock HTTP request.
	promhttp.Handler().ServeHTTP(recorder, req)

	// Check the HTTP status code of the response.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("metricsHandler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check the response body.
	// We're just checking that it includes the metrics we defined.
	body := recorder.Body.String()
	if !strings.Contains(body, "login_requests_total") {
		t.Errorf("metricsHandler did not include login_requests_total metric")
	}
	if !strings.Contains(body, "secret_requests_total") {
		t.Errorf("metricsHandler did not include secret_requests_total metric")
	}
}
