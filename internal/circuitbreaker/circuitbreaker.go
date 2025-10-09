package circuitbreaker

import (
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"
)

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker implements the circuit breaker pattern for external services
type CircuitBreaker struct {
	serviceName      string
	failureThreshold int
	timeout          time.Duration
	maxRetries       int
	retryDelay       time.Duration

	state        CircuitState
	failures     int
	lastFailTime time.Time
	mutex        sync.RWMutex
}

// ServiceMetrics tracks metrics for service calls
type ServiceMetrics struct {
	TotalCalls   int64
	SuccessCalls int64
	FailureCalls int64
	CircuitOpen  bool
	LastCallTime time.Time
	mutex        sync.RWMutex
}

// Global circuit breakers and metrics for each service
var (
	circuitBreakers map[string]*CircuitBreaker
	serviceMetrics  map[string]*ServiceMetrics
	cbMutex         sync.RWMutex
)

// Init initializes a circuit breaker for a service
func Init(serviceName string, failureThreshold int, timeout time.Duration, maxRetries int, retryDelay time.Duration) {
	cbMutex.Lock()
	defer cbMutex.Unlock()
	
	if circuitBreakers == nil {
		circuitBreakers = make(map[string]*CircuitBreaker)
		serviceMetrics = make(map[string]*ServiceMetrics)
	}
	
	circuitBreakers[serviceName] = &CircuitBreaker{
		serviceName:      serviceName,
		failureThreshold: failureThreshold,
		timeout:          timeout,
		maxRetries:       maxRetries,
		retryDelay:       retryDelay,
		state:            StateClosed,
		failures:         0,
	}
	serviceMetrics[serviceName] = &ServiceMetrics{}
}

// Get gets an existing circuit breaker for a service
func Get(serviceName string) *CircuitBreaker {
	cbMutex.RLock()
	defer cbMutex.RUnlock()
	
	cb, exists := circuitBreakers[serviceName]
	if !exists {
		return nil
	}
	return cb
}

// Call attempts to make a call through the circuit breaker
func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// Check if circuit is open
	if cb.state == StateOpen {
		if time.Since(cb.lastFailTime) < cb.timeout {
			return fmt.Errorf("circuit breaker is open for service %s", cb.serviceName)
		}
		// Transition to half-open
		cb.state = StateHalfOpen
	}

	// Attempt the call
	err := fn()

	// Update metrics
	cbMutex.RLock()
	metrics := serviceMetrics[cb.serviceName]
	cbMutex.RUnlock()

	if metrics != nil {
		metrics.mutex.Lock()
		metrics.TotalCalls++
		metrics.LastCallTime = time.Now()
		
		if err != nil {
			metrics.FailureCalls++
			cb.failures++
			cb.lastFailTime = time.Now()

			// Open circuit if failure threshold is reached
			if cb.failures >= cb.failureThreshold {
				cb.state = StateOpen
			}
			metrics.CircuitOpen = (cb.state == StateOpen)
		} else {
			metrics.SuccessCalls++
			// Reset on success
			cb.failures = 0
			if cb.state == StateHalfOpen {
				cb.state = StateClosed
			}
			metrics.CircuitOpen = false
		}
		metrics.mutex.Unlock()
	}

	return err
}

// HTTPCall makes an HTTP call through the circuit breaker
func (cb *CircuitBreaker) HTTPCall(client *http.Client, req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	callErr := cb.Call(func() error {
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
		
		// Consider 5xx responses as failures
		if resp.StatusCode >= 500 {
			return fmt.Errorf("server error: %d", resp.StatusCode)
		}
		
		return nil
	})

	if callErr != nil {
		return nil, callErr
	}

	return resp, err
}

// Reset resets the circuit breaker state
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	cb.state = StateClosed
	cb.failures = 0
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// GetAllStatus returns the status of all circuit breakers
func GetAllStatus() map[string]interface{} {
	cbMutex.RLock()
	defer cbMutex.RUnlock()

	status := make(map[string]interface{})
	
	for serviceName, cb := range circuitBreakers {
		metrics := serviceMetrics[serviceName]
		
		var successRate float64
		if metrics != nil && metrics.TotalCalls > 0 {
			successRate = math.Round((float64(metrics.SuccessCalls)/float64(metrics.TotalCalls))*10000) / 100
		}

		status[serviceName] = map[string]interface{}{
			"state":         cb.GetState(),
			"failures":      cb.failures,
			"total_calls":   metrics.TotalCalls,
			"success_calls": metrics.SuccessCalls,
			"failure_calls": metrics.FailureCalls,
			"success_rate":  successRate,
			"last_call":     metrics.LastCallTime.Unix(),
		}
	}

	return status
}

// ResetByName resets a circuit breaker by service name
func ResetByName(serviceName string) error {
	cbMutex.RLock()
	cb, exists := circuitBreakers[serviceName]
	cbMutex.RUnlock()

	if !exists {
		return fmt.Errorf("circuit breaker for service %s not found", serviceName)
	}

	cb.Reset()
	return nil
}

// String returns a string representation of the circuit state
func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}