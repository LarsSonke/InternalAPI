package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"InternalAPI/internal/circuitbreaker"
	"InternalAPI/internal/config"
)

// HTTPClient is the global HTTP client with timeout
var HTTPClient = &http.Client{Timeout: 30 * time.Second}

// ExternalService handles calls to external services with circuit breaker protection
type ExternalService struct {
	config *config.Config
}

// New creates a new external service client
func New(config *config.Config) *ExternalService {
	return &ExternalService{
		config: config,
	}
}

// Call makes a call to an external service with circuit breaker protection
func (es *ExternalService) Call(serviceName, method, endpoint string, data interface{}) (map[string]interface{}, error) {
	var url, authKey string

	switch serviceName {
	case "beheerder", "api-beheerder":
		url = es.config.APIBeheerderURL + endpoint
		authKey = es.config.APIBeheerderKey
	case "central", "central-mgmt":
		url = es.config.CentralMgmtURL + endpoint
		authKey = es.config.CentralMgmtKey
	default:
		return nil, fmt.Errorf("unknown service: %s", serviceName)
	}

	// Get circuit breaker for this service
	cb := circuitbreaker.Get(serviceName)
	if cb == nil {
		return nil, fmt.Errorf("circuit breaker not initialized for service: %s", serviceName)
	}

	var response map[string]interface{}
	err := cb.Call(func() error {
		return es.makeHTTPCall(method, url, authKey, data, &response)
	})

	return response, err
}

// makeHTTPCall performs the actual HTTP request
func (es *ExternalService) makeHTTPCall(method, url, authKey string, data interface{}, response *map[string]interface{}) error {
	var body []byte
	var err error

	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal request data: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Key", authKey)

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Check HTTP status
	if resp.StatusCode >= 400 {
		if errorMsg, exists := (*response)["error"]; exists {
			return fmt.Errorf("external service error: %v", errorMsg)
		}
		return fmt.Errorf("external service returned status %d", resp.StatusCode)
	}

	return nil
}