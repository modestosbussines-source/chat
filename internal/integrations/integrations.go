package integrations

import (
	"context"

	"github.com/zerodha/logf"
)

// Config holds all integration configurations
type Config struct {
	Evolution EvolutionConfig `json:"evolution"`
	N8n       N8nConfig       `json:"n8n"`
}

// Manager manages all integrations
type Manager struct {
	config    Config
	evolution *EvolutionClient
	n8n       *N8nClient
	log       logf.Logger
}

// NewManager creates a new integrations manager
func NewManager(config Config, log logf.Logger) *Manager {
	m := &Manager{
		config: config,
		log:    log,
	}

	// Initialize Evolution client
	if config.Evolution.Enabled {
		m.evolution = NewEvolutionClient(config.Evolution, log)
		log.Info("Evolution API integration enabled", "url", config.Evolution.APIURL)
	}

	// Initialize n8n client
	if config.N8n.Enabled {
		m.n8n = NewN8nClient(config.N8n, log)
		log.Info("n8n integration enabled", "webhook", config.N8n.WebhookURL)
	}

	return m
}

// Evolution returns the Evolution API client (nil if not enabled)
func (m *Manager) Evolution() *EvolutionClient {
	return m.evolution
}

// N8n returns the n8n client (nil if not enabled)
func (m *Manager) N8n() *N8nClient {
	return m.n8n
}

// IsEvolutionEnabled returns whether Evolution is enabled and configured
func (m *Manager) IsEvolutionEnabled() bool {
	return m.evolution != nil && m.evolution.IsEnabled()
}

// IsN8nEnabled returns whether n8n is enabled and configured
func (m *Manager) IsN8nEnabled() bool {
	return m.n8n != nil && m.n8n.IsEnabled()
}

// HealthCheck performs health checks on all enabled integrations
func (m *Manager) HealthCheck(ctx context.Context) map[string]error {
	results := make(map[string]error)

	if m.IsEvolutionEnabled() {
		results["evolution"] = m.evolution.HealthCheck(ctx)
	}

	if m.IsN8nEnabled() {
		results["n8n"] = m.n8n.HealthCheck(ctx)
	}

	return results
}

// NotifyMessageSent notifies all integrations about a sent message
func (m *Manager) NotifyMessageSent(ctx context.Context, messageID, contactID, content string, metadata map[string]interface{}) {
	if m.IsN8nEnabled() {
		if err := m.n8n.TriggerMessageSent(ctx, messageID, contactID, content, metadata); err != nil {
			m.log.Error("Failed to notify n8n about sent message", "error", err, "message_id", messageID)
		}
	}
}

// NotifyMessageReceived notifies all integrations about a received message
func (m *Manager) NotifyMessageReceived(ctx context.Context, messageID, contactID, content string, metadata map[string]interface{}) {
	if m.IsN8nEnabled() {
		if err := m.n8n.TriggerMessageReceived(ctx, messageID, contactID, content, metadata); err != nil {
			m.log.Error("Failed to notify n8n about received message", "error", err, "message_id", messageID)
		}
	}
}

// NotifyContactCreated notifies all integrations about a new contact
func (m *Manager) NotifyContactCreated(ctx context.Context, contactID, name, phone string) {
	if m.IsN8nEnabled() {
		if err := m.n8n.TriggerContactCreated(ctx, contactID, name, phone); err != nil {
			m.log.Error("Failed to notify n8n about new contact", "error", err, "contact_id", contactID)
		}
	}
}

// NotifyCampaignStarted notifies all integrations about campaign start
func (m *Manager) NotifyCampaignStarted(ctx context.Context, campaignID, campaignName string, recipientCount int) {
	if m.IsN8nEnabled() {
		if err := m.n8n.TriggerCampaignStarted(ctx, campaignID, campaignName, recipientCount); err != nil {
			m.log.Error("Failed to notify n8n about campaign start", "error", err, "campaign_id", campaignID)
		}
	}
}

// NotifyCampaignCompleted notifies all integrations about campaign completion
func (m *Manager) NotifyCampaignCompleted(ctx context.Context, campaignID string, stats map[string]interface{}) {
	if m.IsN8nEnabled() {
		if err := m.n8n.TriggerCampaignCompleted(ctx, campaignID, stats); err != nil {
			m.log.Error("Failed to notify n8n about campaign completion", "error", err, "campaign_id", campaignID)
		}
	}
}

// SendWhatsAppViaEvolution sends a WhatsApp message using Evolution API
func (m *Manager) SendWhatsAppViaEvolution(ctx context.Context, instanceName, number, text string) error {
	if !m.IsEvolutionEnabled() {
		return ErrEvolutionNotEnabled
	}

	_, err := m.evolution.SendMessage(ctx, instanceName, number, text)
	return err
}

// SendMediaViaEvolution sends media using Evolution API
func (m *Manager) SendMediaViaEvolution(ctx context.Context, instanceName string, msg *EvolutionMessage) error {
	if !m.IsEvolutionEnabled() {
		return ErrEvolutionNotEnabled
	}

	_, err := m.evolution.SendMedia(ctx, instanceName, msg)
	return err
}

// Custom errors
var (
	ErrEvolutionNotEnabled = ErrIntegration("Evolution API is not enabled")
	ErrN8nNotEnabled       = ErrIntegration("n8n is not enabled")
)

// ErrIntegration represents an integration error
type ErrIntegration string

func (e ErrIntegration) Error() string {
	return string(e)
}
