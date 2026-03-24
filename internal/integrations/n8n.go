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

// N8nConfig holds n8n integration configuration
type N8nConfig struct {
	Enabled       bool   `json:"enabled"`
	WebhookURL    string `json:"webhook_url"`
	APIKey        string `json:"api_key"`
	WebhookSecret string `json:"webhook_secret,omitempty"`
}

// N8nClient handles communication with n8n workflows
type N8nClient struct {
	config     N8nConfig
	httpClient *http.Client
	log        logf.Logger
}

// N8nWebhookPayload represents the payload sent to n8n webhooks
type N8nWebhookPayload struct {
	Event     string                 `json:"event"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// N8nWorkflowResponse represents a workflow execution response
type N8nWorkflowResponse struct {
	ExecutionID string `json:"executionId"`
	Status      string `json:"status"`
	StartedAt   string `json:"startedAt,omitempty"`
	FinishedAt  string `json:"finishedAt,omitempty"`
}

// N8nHealthResponse represents health check response
type N8nHealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
}

// NewN8nClient creates a new n8n integration client
func NewN8nClient(config N8nConfig, log logf.Logger) *N8nClient {
	return &N8nClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		log: log,
	}
}

// IsEnabled returns whether n8n integration is enabled
func (c *N8nClient) IsEnabled() bool {
	return c.config.Enabled && c.config.WebhookURL != ""
}

// TriggerWorkflow triggers an n8n workflow via webhook
func (c *N8nClient) TriggerWorkflow(ctx context.Context, event string, data map[string]interface{}) (*N8nWorkflowResponse, error) {
	payload := N8nWebhookPayload{
		Event:     event,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
		Metadata: map[string]interface{}{
			"source":  "omni",
			"version": "1.0",
		},
	}

	return c.sendWebhook(ctx, payload)
}

// TriggerMessageSent triggers workflow when a message is sent
func (c *N8nClient) TriggerMessageSent(ctx context.Context, messageID, contactID, content string, metadata map[string]interface{}) error {
	data := map[string]interface{}{
		"message_id": messageID,
		"contact_id": contactID,
		"content":    content,
		"type":       "sent",
		"metadata":   metadata,
	}

	_, err := c.TriggerWorkflow(ctx, "message.sent", data)
	return err
}

// TriggerMessageReceived triggers workflow when a message is received
func (c *N8nClient) TriggerMessageReceived(ctx context.Context, messageID, contactID, content string, metadata map[string]interface{}) error {
	data := map[string]interface{}{
		"message_id": messageID,
		"contact_id": contactID,
		"content":    content,
		"type":       "received",
		"metadata":   metadata,
	}

	_, err := c.TriggerWorkflow(ctx, "message.received", data)
	return err
}

// TriggerContactCreated triggers workflow when a contact is created
func (c *N8nClient) TriggerContactCreated(ctx context.Context, contactID, name, phone string) error {
	data := map[string]interface{}{
		"contact_id": contactID,
		"name":       name,
		"phone":      phone,
		"type":       "created",
	}

	_, err := c.TriggerWorkflow(ctx, "contact.created", data)
	return err
}

// TriggerCampaignStarted triggers workflow when a campaign starts
func (c *N8nClient) TriggerCampaignStarted(ctx context.Context, campaignID, campaignName string, recipientCount int) error {
	data := map[string]interface{}{
		"campaign_id":     campaignID,
		"campaign_name":   campaignName,
		"recipient_count": recipientCount,
		"type":            "started",
	}

	_, err := c.TriggerWorkflow(ctx, "campaign.started", data)
	return err
}

// TriggerCampaignCompleted triggers workflow when a campaign completes
func (c *N8nClient) TriggerCampaignCompleted(ctx context.Context, campaignID string, stats map[string]interface{}) error {
	data := map[string]interface{}{
		"campaign_id": campaignID,
		"stats":       stats,
		"type":        "completed",
	}

	_, err := c.TriggerWorkflow(ctx, "campaign.completed", data)
	return err
}

// TriggerEvent is a generic event trigger
func (c *N8nClient) TriggerEvent(ctx context.Context, eventType string, data map[string]interface{}) error {
	_, err := c.TriggerWorkflow(ctx, eventType, data)
	return err
}

// sendWebhook sends payload to n8n webhook
func (c *N8nClient) sendWebhook(ctx context.Context, payload N8nWebhookPayload) (*N8nWorkflowResponse, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add authentication
	if c.config.APIKey != "" {
		req.Header.Set("X-API-Key", c.config.APIKey)
	}

	// Add webhook signature for verification
	if c.config.WebhookSecret != "" {
		req.Header.Set("X-Webhook-Secret", c.config.WebhookSecret)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("webhook request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var workflowResp N8nWorkflowResponse
	if len(respBody) > 0 {
		json.Unmarshal(respBody, &workflowResp)
	}

	if resp.StatusCode >= 400 {
		return &workflowResp, fmt.Errorf("webhook returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return &workflowResp, nil
}

// HealthCheck checks if n8n is reachable
func (c *N8nClient) HealthCheck(ctx context.Context) error {
	// n8n health endpoint
	healthURL := fmt.Sprintf("%s/healthz", getBaseURL(c.config.WebhookURL))

	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
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

// GetWorkflows retrieves available workflows (requires API key)
func (c *N8nClient) GetWorkflows(ctx context.Context) ([]map[string]interface{}, error) {
	if c.config.APIKey == "" {
		return nil, fmt.Errorf("API key required for this operation")
	}

	// Use REST API endpoint
	apiURL := fmt.Sprintf("%s/api/v1/workflows", getBaseURL(c.config.WebhookURL))

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-N8N-API-KEY", c.config.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var result struct {
		Data []map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// getBaseURL extracts base URL from webhook URL
func getBaseURL(webhookURL string) string {
	// Find protocol end
	protocolEnd := 0
	for i := 0; i < len(webhookURL)-2; i++ {
		if webhookURL[i:i+3] == "://" {
			protocolEnd = i + 3
			break
		}
	}

	// Find first slash after protocol (path start)
	for i := protocolEnd; i < len(webhookURL); i++ {
		if webhookURL[i] == '/' {
			return webhookURL[:i]
		}
	}

	return webhookURL
}
