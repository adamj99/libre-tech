package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

var (
	loginRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "login_requests_total",
		Help: "The total number of /login requests",
	})
	secretRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "secret_requests_total",
		Help: "The total number of /secret requests",
	})
)

type response struct {
	Secret string `json:"secret"`
}

func main() {
	// Tell Viper to read environment variables
	viper.AutomaticEnv()
	// Get the value of the APP_SECRET environment variable
	app_secret := viper.GetString("APP_SECRET")
	// If APP_SECRET is not set, log an error and stop the application
	if app_secret == "" {
		log.Fatal("APP_SECRET environment variable is not set")
	}
	// Create a new router
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/secret", func(w http.ResponseWriter, r *http.Request) {
		secretHandler(w, r, app_secret)
	})
	// Start the HTTP server on port 8080, and log an error if it fails to start
	err := http.ListenAndServe((":8080"), r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Handler for /login route
func loginHandler(writer http.ResponseWriter, request *http.Request) {
	// Increment the loginRequests counter
	loginRequests.Inc()
	writer.WriteHeader(http.StatusOK)
}

// Handler for /secret route
func secretHandler(w http.ResponseWriter, r *http.Request, app_secret string) {
	// Increment the secretRequests counter
	secretRequests.Inc()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response{Secret: app_secret}); err != nil {
		log.Printf("Failed to marshal secret for endpoint /secret: %v", err)
		http.Error(w, "Failed to create JSON response for /secret", http.StatusInternalServerError)
		return
	}
}
