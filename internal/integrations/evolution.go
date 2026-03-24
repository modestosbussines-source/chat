package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zerodha/logf"
)

// EvolutionConfig holds Evolution API configuration
type EvolutionConfig struct {
	Enabled  bool   `json:"enabled"`
	APIURL   string `json:"api_url"`
	APIKey   string `json:"api_key"`
	Instance string `json:"instance,omitempty"`
}

// EvolutionClient handles communication with Evolution API
type EvolutionClient struct {
	config     EvolutionConfig
	httpClient *http.Client
	log        logf.Logger
}

// EvolutionInstance represents a WhatsApp instance in Evolution
type EvolutionInstance struct {
	Instance struct {
		InstanceName string `json:"instanceName"`
		State        string `json:"state"`
	} `json:"instance"`
}

// EvolutionMessage represents a message to send via Evolution
type EvolutionMessage struct {
	Number  string                 `json:"number"`
	Text    string                 `json:"text,omitempty"`
	Media   *EvolutionMedia        `json:"media,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// EvolutionMedia represents media attachment
type EvolutionMedia struct {
	Media    string `json:"media"` // URL or base64
	Caption  string `json:"caption,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

// EvolutionResponse represents API response
type EvolutionResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewEvolutionClient creates a new Evolution API client
func NewEvolutionClient(config EvolutionConfig, log logf.Logger) *EvolutionClient {
	return &EvolutionClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		log: log,
	}
}

// IsEnabled returns whether Evolution integration is enabled
func (c *EvolutionClient) IsEnabled() bool {
	return c.config.Enabled && c.config.APIURL != "" && c.config.APIKey != ""
}

// CreateInstance creates a new WhatsApp instance
func (c *EvolutionClient) CreateInstance(ctx context.Context, instanceName string) (*EvolutionInstance, error) {
	url := fmt.Sprintf("%s/instance/create", c.config.APIURL)

	body := map[string]string{
		"instanceName": instanceName,
	}

	resp, err := c.doRequest(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}

	var instance EvolutionInstance
	if data, ok := resp.Data.(map[string]interface{}); ok {
		jsonData, _ := json.Marshal(data)
		json.Unmarshal(jsonData, &instance)
	}

	return &instance, nil
}

// GetInstance gets instance status
func (c *EvolutionClient) GetInstance(ctx context.Context, instanceName string) (*EvolutionInstance, error) {
	url := fmt.Sprintf("%s/instance/%s", c.config.APIURL, instanceName)

	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var instance EvolutionInstance
	if data, ok := resp.Data.(map[string]interface{}); ok {
		jsonData, _ := json.Marshal(data)
		json.Unmarshal(jsonData, &instance)
	}

	return &instance, nil
}

// DeleteInstance deletes an instance
func (c *EvolutionClient) DeleteInstance(ctx context.Context, instanceName string) error {
	url := fmt.Sprintf("%s/instance/%s", c.config.APIURL, instanceName)
	_, err := c.doRequest(ctx, "DELETE", url, nil)
	return err
}

// GetQRCode gets QR code for instance connection
func (c *EvolutionClient) GetQRCode(ctx context.Context, instanceName string) (string, error) {
	url := fmt.Sprintf("%s/instance/%s/qrcode", c.config.APIURL, instanceName)

	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	if qrData, ok := resp.Data.(map[string]interface{}); ok {
		if code, ok := qrData["qrcode"].(string); ok {
			return code, nil
		}
	}

	return "", fmt.Errorf("QR code not available")
}

// SendMessage sends a text message
func (c *EvolutionClient) SendMessage(ctx context.Context, instanceName, number, text string) (*EvolutionResponse, error) {
	url := fmt.Sprintf("%s/instance/%s/message/text", c.config.APIURL, instanceName)

	msg := EvolutionMessage{
		Number: number,
		Text:   text,
	}

	return c.doRequest(ctx, "POST", url, msg)
}

// SendMedia sends a media message (image, video, document)
func (c *EvolutionClient) SendMedia(ctx context.Context, instanceName string, msg *EvolutionMessage) (*EvolutionResponse, error) {
	url := fmt.Sprintf("%s/instance/%s/message/media", c.config.APIURL, instanceName)
	return c.doRequest(ctx, "POST", url, msg)
}

// SendAudio sends an audio message
func (c *EvolutionClient) SendAudio(ctx context.Context, instanceName, number, audioURL string) (*EvolutionResponse, error) {
	url := fmt.Sprintf("%s/instance/%s/message/audio", c.config.APIURL, instanceName)

	msg := EvolutionMessage{
		Number: number,
		Media: &EvolutionMedia{
			Media: audioURL,
		},
	}

	return c.doRequest(ctx, "POST", url, msg)
}

// SetWebhook sets webhook URL for instance events
func (c *EvolutionClient) SetWebhook(ctx context.Context, instanceName, webhookURL string) error {
	url := fmt.Sprintf("%s/instance/%s/webhook", c.config.APIURL, instanceName)

	body := map[string]interface{}{
		"url":    webhookURL,
		"events": []string{"messages.upsert", "connection.update"},
	}

	_, err := c.doRequest(ctx, "POST", url, body)
	return err
}

// GetConnections gets connection status
func (c *EvolutionClient) GetConnections(ctx context.Context, instanceName string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/instance/%s/connection", c.config.APIURL, instanceName)

	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	if data, ok := resp.Data.(map[string]interface{}); ok {
		return data, nil
	}

	return nil, fmt.Errorf("invalid connection data")
}

// doRequest performs an HTTP request to Evolution API
func (c *EvolutionClient) doRequest(ctx context.Context, method, url string, body interface{}) (*EvolutionResponse, error) {
	var reqBody io.Reader

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", c.config.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var evolutionResp EvolutionResponse
	if err := json.Unmarshal(respBody, &evolutionResp); err != nil {
		// Try to parse as direct data
		evolutionResp.Success = resp.StatusCode >= 200 && resp.StatusCode < 300
		json.Unmarshal(respBody, &evolutionResp.Data)
	}

	if resp.StatusCode >= 400 {
		return &evolutionResp, fmt.Errorf("API error: %s", evolutionResp.Error)
	}

	return &evolutionResp, nil
}

// HealthCheck checks if Evolution API is reachable
func (c *EvolutionClient) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/health", c.config.APIURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}
