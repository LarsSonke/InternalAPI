package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// PluginRegistration represents the registration payload sent to the broker
type PluginRegistration struct {
	Description   string   `json:"description"`
	Version       string   `json:"version"`
	Slug          string   `json:"slug"`
	Name          string   `json:"name"`
	Category      string   `json:"category,omitempty"`
	Host          string   `json:"host"`
	BaseAPIRoute  string   `json:"base-api-route"`
	SettingsRoute string   `json:"settings-route,omitempty"`
	APIRoutes     []string `json:"api-routes,omitempty"`
	Enabled       bool     `json:"enabled"`
}

// RegisterWithBroker registers InternalAPI with the broker on startup
// This is non-blocking and won't fail the application if broker is unavailable
func RegisterWithBroker(host, port string) {
	brokerURL := os.Getenv("BROKER_URL")
	if brokerURL == "" {
		brokerURL = "http://localhost:8081" // Default broker URL
		log.Info("BROKER_URL not set, using default: http://localhost:8081")
	}

	brokerAuthToken := os.Getenv("BROKER_AUTH_TOKEN")
	if brokerAuthToken == "" {
		log.Warn("⚠️  BROKER_AUTH_TOKEN not set - broker registration may fail if authentication is required")
		// Don't return - attempt registration anyway in case broker allows unauthenticated registration
	}

	// Construct the full host URL
	serviceHost := fmt.Sprintf("http://%s:%s", host, port)

	registration := PluginRegistration{
		Description:   "Hotel Internal API - Gateway for user portal and admin services",
		Version:       "2.0.0",
		Slug:          "internal-api",
		Name:          "Hotel Internal API",
		Category:      "gateway",
		Host:          serviceHost,
		BaseAPIRoute:  "/api/v1",
		SettingsRoute: "/admin/system/stats",
		APIRoutes: []string{
			"/api/v1/albums",
			"/api/v1/guests",
			"/api/v1/reservations",
			"/api/auth/login",
			"/api/auth/logout",
			"/api/auth/refresh",
			"/admin/users",
			"/admin/roles",
			"/health",
		},
		Enabled: true,
	}

	// Run registration in background to not block startup
	go func() {
		// Wait a moment for InternalAPI to be fully ready
		time.Sleep(2 * time.Second)

		if err := attemptRegistration(brokerURL, brokerAuthToken, registration); err != nil {
			log.WithError(err).Error("Failed to register with broker - service will continue running but won't receive proxied traffic")
		} else {
			log.WithFields(logrus.Fields{
				"broker_url":  brokerURL,
				"plugin_slug": registration.Slug,
				"host":        registration.Host,
			}).Info("✓ Successfully registered with broker")
		}
	}()
}

// attemptRegistration performs the actual HTTP request to register with the broker
func attemptRegistration(brokerURL, authToken string, registration PluginRegistration) error {
	payload, err := json.Marshal(registration)
	if err != nil {
		return fmt.Errorf("failed to marshal registration payload: %w", err)
	}

	req, err := http.NewRequest("POST", brokerURL+"/api/v1/route", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create registration request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send registration request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("registration failed with status %d: %v", resp.StatusCode, errResp)
	}

	return nil
}
