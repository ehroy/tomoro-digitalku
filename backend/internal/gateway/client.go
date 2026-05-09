package gateway

import (
	"bytes"
	"coffee-order/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TomoroClient is the HTTP client for Tomoro Coffee API
type TomoroClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// RequestContext contains authentication and location data
type RequestContext struct {
	Token      string
	DeviceCode string
	WToken     string
	UCDE       string
	Latitude   float64
	Longitude  float64
}

// NewTomoroClient creates a new Tomoro API client
func NewTomoroClient() *TomoroClient {
	return &TomoroClient{
		BaseURL: config.TomoroAPIBaseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Request makes an HTTP request to Tomoro API with proper headers
func (c *TomoroClient) Request(ctx RequestContext, method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set standard headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "id-ID;q=1.0")
	req.Header.Set("User-Agent", config.APIVersion)
	req.Header.Set("Connection", "keep-alive")

	// Set app-specific headers
	req.Header.Set("appLanguage", config.AppLanguage)
	req.Header.Set("countryCode", config.CountryCode)
	req.Header.Set("timeZone", config.TimeZone)
	req.Header.Set("appChannel", config.AppChannel)
	req.Header.Set("revision", config.APIVersion)

	// Set authentication headers if available
	if ctx.Token != "" {
		req.Header.Set("token", ctx.Token)
	}
	if ctx.DeviceCode != "" {
		req.Header.Set("deviceCode", ctx.DeviceCode)
	}
	if ctx.WToken != "" {
		req.Header.Set("wToken", ctx.WToken)
	}
	if ctx.UCDE != "" {
		req.Header.Set("ucde", ctx.UCDE)
	}

	// Set location headers if available
	if ctx.Latitude != 0 {
		req.Header.Set("latitude", fmt.Sprintf("%f", ctx.Latitude))
	}
	if ctx.Longitude != 0 {
		req.Header.Set("longitude", fmt.Sprintf("%f", ctx.Longitude))
	}

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
