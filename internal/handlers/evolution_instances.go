package handlers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/shridarpatil/whatomate/internal/models"
	"github.com/valyala/fasthttp"
	"github.com/zerodha/fastglue"
)

// getWebhookBaseURL returns the base URL for webhooks
func (a *App) getWebhookBaseURL() string {
	// Try environment variable first
	if url := os.Getenv("WEBHOOK_BASE_URL"); url != "" {
		return url
	}
	// Try from config
	if a.Config != nil && a.Config.Server.BasePath != "" {
		return a.Config.Server.BasePath
	}
	// Default fallback
	return "http://localhost:8080"
}

// ListEvolutionInstances returns all Evolution instances for the organization
func (a *App) ListEvolutionInstances(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	var instances []models.EvolutionInstance
	if err := a.DB.Where("organization_id = ?", orgID).Order("created_at desc").Find(&instances).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to load instances", nil, "")
	}

	return r.SendEnvelope(instances)
}

// CreateEvolutionInstanceRequest represents request body
type CreateEvolutionInstanceRequest struct {
	InstanceName string `json:"instance_name"`
	DisplayName  string `json:"display_name"`
}

// CreateEvolutionInstance creates a new Evolution WhatsApp instance
func (a *App) CreateEvolutionInstance(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	// Check if Evolution integration is enabled
	if a.IntegrationManager == nil || !a.IntegrationManager.IsEvolutionEnabled() {
		return r.SendErrorEnvelope(fasthttp.StatusServiceUnavailable, "Evolution integration not enabled", nil, "")
	}

	var req CreateEvolutionInstanceRequest
	if err := a.decodeRequest(r, &req); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid request body", nil, "")
	}

	if req.InstanceName == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Instance name is required", nil, "")
	}

	// Check if instance already exists
	var existing models.EvolutionInstance
	if err := a.DB.Where("instance_name = ?", req.InstanceName).First(&existing).Error; err == nil {
		return r.SendErrorEnvelope(fasthttp.StatusConflict, "Instance name already exists", nil, "")
	}

	// Create instance in Evolution API
	ctx := context.Background()
	evolutionClient := a.IntegrationManager.Evolution()

	_, err = evolutionClient.CreateInstance(ctx, req.InstanceName)
	if err != nil {
		a.Log.Error("Failed to create Evolution instance", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, fmt.Sprintf("Failed to create instance: %v", err), nil, "")
	}

	// Generate webhook URL (use environment variable or construct from server config)
	webhookURL := a.getWebhookBaseURL() + "/api/webhooks/evolution"

	// Set webhook for the instance
	if err := evolutionClient.SetWebhook(ctx, req.InstanceName, webhookURL); err != nil {
		a.Log.Warn("Failed to set webhook", "error", err)
		// Continue anyway - webhook can be set later
	}

	// Save instance to database
	instance := models.EvolutionInstance{
		BaseModel:      models.BaseModel{ID: uuid.New()},
		OrganizationID: orgID,
		InstanceName:   req.InstanceName,
		DisplayName:    req.DisplayName,
		Status:         models.EvolutionStatusDisconnected,
		WebhookURL:     webhookURL,
	}

	if err := a.DB.Create(&instance).Error; err != nil {
		// Try to clean up Evolution instance
		_ = evolutionClient.DeleteInstance(ctx, req.InstanceName)
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to save instance", nil, "")
	}

	return r.SendEnvelope(instance)
}

// GetEvolutionInstance returns a specific instance
func (a *App) GetEvolutionInstance(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	instanceID := r.RequestCtx.UserValue("id").(string)

	var instance models.EvolutionInstance
	if err := a.DB.Where("id = ? AND organization_id = ?", instanceID, orgID).First(&instance).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Instance not found", nil, "")
	}

	return r.SendEnvelope(instance)
}

// DeleteEvolutionInstance deletes an Evolution instance
func (a *App) DeleteEvolutionInstance(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	instanceID := r.RequestCtx.UserValue("id").(string)

	var instance models.EvolutionInstance
	if err := a.DB.Where("id = ? AND organization_id = ?", instanceID, orgID).First(&instance).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Instance not found", nil, "")
	}

	// Delete from Evolution API
	if a.IntegrationManager != nil && a.IntegrationManager.IsEvolutionEnabled() {
		ctx := context.Background()
		evolutionClient := a.IntegrationManager.Evolution()
		_ = evolutionClient.DeleteInstance(ctx, instance.InstanceName)
	}

	// Delete from database
	if err := a.DB.Delete(&instance).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to delete instance", nil, "")
	}

	return r.SendEnvelope(map[string]string{"status": "deleted"})
}

// GetEvolutionQRCode returns QR code for connecting WhatsApp
func (a *App) GetEvolutionQRCode(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	instanceID := r.RequestCtx.UserValue("id").(string)

	var instance models.EvolutionInstance
	if err := a.DB.Where("id = ? AND organization_id = ?", instanceID, orgID).First(&instance).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Instance not found", nil, "")
	}

	if a.IntegrationManager == nil || !a.IntegrationManager.IsEvolutionEnabled() {
		return r.SendErrorEnvelope(fasthttp.StatusServiceUnavailable, "Evolution integration not enabled", nil, "")
	}

	ctx := context.Background()
	evolutionClient := a.IntegrationManager.Evolution()

	// Update status to connecting
	a.DB.Model(&instance).Update("status", models.EvolutionStatusConnecting)

	// Get QR code
	qrCode, err := evolutionClient.GetQRCode(ctx, instance.InstanceName)
	if err != nil {
		a.Log.Error("Failed to get QR code", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, fmt.Sprintf("Failed to get QR code: %v", err), nil, "")
	}

	return r.SendEnvelope(map[string]string{
		"qrcode": qrCode,
		"status": instance.Status,
	})
}

// GetEvolutionConnectionStatus returns connection status
func (a *App) GetEvolutionConnectionStatus(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	instanceID := r.RequestCtx.UserValue("id").(string)

	var instance models.EvolutionInstance
	if err := a.DB.Where("id = ? AND organization_id = ?", instanceID, orgID).First(&instance).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Instance not found", nil, "")
	}

	// Get status from Evolution API
	status := map[string]interface{}{
		"status":       instance.Status,
		"phone":        instance.Phone,
		"profile_name": instance.ProfileName.String,
		"connected_at": instance.ConnectedAt,
	}

	if a.IntegrationManager != nil && a.IntegrationManager.IsEvolutionEnabled() {
		ctx := context.Background()
		evolutionClient := a.IntegrationManager.Evolution()

		connData, err := evolutionClient.GetConnections(ctx, instance.InstanceName)
		if err == nil {
			status["evolution_data"] = connData
		}
	}

	return r.SendEnvelope(status)
}

// DisconnectEvolutionInstance disconnects WhatsApp from instance
func (a *App) DisconnectEvolutionInstance(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	instanceID := r.RequestCtx.UserValue("id").(string)

	var instance models.EvolutionInstance
	if err := a.DB.Where("id = ? AND organization_id = ?", instanceID, orgID).First(&instance).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Instance not found", nil, "")
	}

	// Delete instance from Evolution (disconnects it)
	if a.IntegrationManager != nil && a.IntegrationManager.IsEvolutionEnabled() {
		ctx := context.Background()
		evolutionClient := a.IntegrationManager.Evolution()

		// Delete and recreate instance
		_ = evolutionClient.DeleteInstance(ctx, instance.InstanceName)

		// Recreate instance
		_, err := evolutionClient.CreateInstance(ctx, instance.InstanceName)
		if err != nil {
			a.Log.Error("Failed to recreate instance", "error", err)
		}

		// Reset webhook
		webhookURL := a.getWebhookBaseURL() + "/api/webhooks/evolution"
		_ = evolutionClient.SetWebhook(ctx, instance.InstanceName, webhookURL)
	}

	// Update status in database
	now := time.Now()
	a.DB.Model(&instance).Updates(map[string]interface{}{
		"status":        models.EvolutionStatusDisconnected,
		"phone":         "",
		"profile_name":  nil,
		"connected_at":  nil,
		"last_activity": &now,
	})

	return r.SendEnvelope(map[string]string{"status": "disconnected"})
}

// SendEvolutionMessage sends a message via Evolution
func (a *App) SendEvolutionMessage(r *fastglue.Request) error {
	orgID, userID, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	var req struct {
		InstanceID string `json:"instance_id"`
		ContactID  string `json:"contact_id"`
		Content    string `json:"content"`
	}
	if err := a.decodeRequest(r, &req); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid request body", nil, "")
	}

	// Get instance
	var instance models.EvolutionInstance
	if err := a.DB.Where("id = ? AND organization_id = ?", req.InstanceID, orgID).First(&instance).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Instance not found", nil, "")
	}

	if instance.Status != models.EvolutionStatusConnected {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Instance not connected", nil, "")
	}

	// Get contact
	var contact models.Contact
	if err := a.DB.Where("id = ? AND organization_id = ?", req.ContactID, orgID).First(&contact).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Contact not found", nil, "")
	}

	// 24h compliance check
	if contact.LastInboundAt != nil {
		hours := time.Since(*contact.LastInboundAt).Hours()
		if hours > 24 {
			return r.SendErrorEnvelope(fasthttp.StatusForbidden, "Lead is outside the 24-hour active window", nil, "")
		}
	}

	// Send via Evolution
	ctx := context.Background()
	evolutionClient := a.IntegrationManager.Evolution()

	resp, err := evolutionClient.SendMessage(ctx, instance.InstanceName, contact.PhoneNumber, req.Content)
	if err != nil {
		a.Log.Error("Failed to send Evolution message", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, fmt.Sprintf("Failed to send: %v", err), nil, "")
	}

	// Save message to database
	message := models.Message{
		BaseModel:       models.BaseModel{ID: uuid.New()},
		OrganizationID:  orgID,
		ContactID:       contact.ID,
		Direction:       models.DirectionOutbound,
		MessageType:     models.MessageTypeText,
		Content:         req.Content,
		Status:          models.MessageStatusSent,
		WhatsAppAccount: instance.InstanceName,
		SentByUserID:    &userID,
		SentAt:          func() *time.Time { t := time.Now(); return &t }(),
	}

	if err := a.DB.Create(&message).Error; err != nil {
		a.Log.Error("Failed to save message", "error", err)
	}

	// Update contact last message
	a.DB.Model(&contact).Update("last_message_at", message.SentAt)

	// Broadcast via WebSocket
	a.broadcastNewMessage(orgID, &message, &contact)

	return r.SendEnvelope(map[string]interface{}{
		"status":             "sent",
		"message_id":         message.ID,
		"evolution_response": resp,
	})
}
