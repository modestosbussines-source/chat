package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shridarpatil/whatomate/internal/models"
	"github.com/valyala/fasthttp"
	"github.com/zerodha/fastglue"
)

// EvolutionWebhook handles incoming webhooks from Evolution API
func (a *App) EvolutionWebhook(r *fastglue.Request) error {
	// Read body (PostBody returns []byte directly)
	body := r.RequestCtx.PostBody()

	// Parse webhook payload
	var payload models.EvolutionWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		a.Log.Error("Failed to parse webhook payload", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid JSON", nil, "")
	}

	a.Log.Info("Evolution webhook received", "event", payload.Event, "instance", payload.Instance)

	// Route to appropriate handler based on event type
	switch payload.Event {
	case models.EvolutionEventMessageUpsert:
		return a.handleEvolutionMessage(r, payload)
	case models.EvolutionEventConnectionUpdate:
		return a.handleEvolutionConnection(r, payload)
	default:
		a.Log.Debug("Unhandled evolution event", "event", payload.Event)
		return r.SendEnvelope(map[string]string{"status": "ignored"})
	}
}

// handleEvolutionMessage processes incoming messages from Evolution
func (a *App) handleEvolutionMessage(r *fastglue.Request, payload models.EvolutionWebhookPayload) error {
	// Find the instance in our database
	var instance models.EvolutionInstance
	if err := a.DB.Where("instance_name = ?", payload.Instance).First(&instance).Error; err != nil {
		a.Log.Error("Evolution instance not found", "instance", payload.Instance)
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Instance not found", nil, "")
	}

	// Parse message data
	messageData, ok := payload.Data["message"].(map[string]interface{})
	if !ok {
		a.Log.Error("Invalid message data in webhook")
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid message data", nil, "")
	}

	// Convert to our message payload struct
	messagePayload := parseEvolutionMessage(messageData, payload.Data)

	// Skip messages sent by us (fromMe)
	if messagePayload.IsFromMe() {
		a.Log.Debug("Skipping message from self", "instance", payload.Instance)
		return r.SendEnvelope(map[string]string{"status": "ignored"})
	}

	// Get or create contact
	contact, err := a.getOrCreateEvolutionContact(&instance, messagePayload)
	if err != nil {
		a.Log.Error("Failed to get/create contact", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to process contact", nil, "")
	}

	// Process the message
	ctx := context.Background()
	if err := a.processIncomingEvolutionMessage(ctx, &instance, contact, messagePayload); err != nil {
		a.Log.Error("Failed to process incoming message", "error", err)
		// Still return success to prevent Evolution from retrying
	}

	return r.SendEnvelope(map[string]string{"status": "ok"})
}

// handleEvolutionConnection processes connection status updates
func (a *App) handleEvolutionConnection(r *fastglue.Request, payload models.EvolutionWebhookPayload) error {
	// Find the instance
	var instance models.EvolutionInstance
	if err := a.DB.Where("instance_name = ?", payload.Instance).First(&instance).Error; err != nil {
		a.Log.Error("Evolution instance not found for connection update", "instance", payload.Instance)
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Instance not found", nil, "")
	}

	// Parse connection data
	connectionData, ok := payload.Data["connection"].(map[string]interface{})
	if !ok {
		return r.SendEnvelope(map[string]string{"status": "ignored"})
	}

	state, _ := connectionData["state"].(string)

	// Update instance status
	updates := map[string]interface{}{
		"last_activity": time.Now(),
	}

	switch state {
	case "open", "connected":
		updates["status"] = models.EvolutionStatusConnected
		now := time.Now()
		updates["connected_at"] = &now
		updates["last_error"] = nil
	case "close", "disconnected":
		updates["status"] = models.EvolutionStatusDisconnected
	case "connecting":
		updates["status"] = models.EvolutionStatusConnecting
	}

	if err := a.DB.Model(&instance).Updates(updates).Error; err != nil {
		a.Log.Error("Failed to update instance status", "error", err)
	}

	// Broadcast status update via WebSocket
	a.broadcastEvolutionStatus(&instance)

	a.Log.Info("Evolution connection updated", "instance", payload.Instance, "state", state)
	return r.SendEnvelope(map[string]string{"status": "ok"})
}

// parseEvolutionMessage converts raw webhook data to our message payload
func parseEvolutionMessage(message map[string]interface{}, data map[string]interface{}) *models.EvolutionMessagePayload {
	payload := &models.EvolutionMessagePayload{}

	// Extract key info
	if key, ok := data["key"].(map[string]interface{}); ok {
		payload.Key.RemoteJid, _ = key["remoteJid"].(string)
		payload.Key.FromMe, _ = key["fromMe"].(bool)
		payload.Key.ID, _ = key["id"].(string)
	}

	// Extract timestamp
	if ts, ok := data["messageTimestamp"].(float64); ok {
		payload.Timestamp = int64(ts)
	}

	// Extract push name
	payload.PushName, _ = data["pushName"].(string)

	// Extract status
	payload.Status, _ = data["status"].(string)

	// Store the message map for later processing
	payload.Message = message

	// Determine message type and extract content
	payload.MessageType = detectMessageType(message)

	// Extract text content
	payload.Conversation = extractConversation(message)

	// Extract media info
	payload.MediaURL = extractMediaURL(message)
	payload.MimeType = extractMimeType(message)
	payload.Caption = extractCaption(message)
	payload.FileName = extractFileName(message)

	return payload
}

// detectMessageType determines the type of message
func detectMessageType(message map[string]interface{}) string {
	// Check for direct message types
	messageTypes := []string{
		"conversation", "extendedTextMessage", "imageMessage", "videoMessage",
		"audioMessage", "documentMessage", "stickerMessage", "locationMessage",
		"contactMessage", "reactionMessage",
	}

	for _, mt := range messageTypes {
		if _, ok := message[mt]; ok {
			return mt
		}
	}

	return "unknown"
}

// extractConversation extracts text from conversation
func extractConversation(message map[string]interface{}) string {
	if conv, ok := message["conversation"].(string); ok {
		return conv
	}
	return ""
}

// extractMediaURL extracts media URL from message
func extractMediaURL(message map[string]interface{}) string {
	mediaTypes := []string{"imageMessage", "videoMessage", "audioMessage", "documentMessage", "stickerMessage"}
	for _, mt := range mediaTypes {
		if media, ok := message[mt].(map[string]interface{}); ok {
			if url, ok := media["url"].(string); ok {
				return url
			}
		}
	}
	return ""
}

// extractMimeType extracts MIME type from message
func extractMimeType(message map[string]interface{}) string {
	mediaTypes := []string{"imageMessage", "videoMessage", "audioMessage", "documentMessage", "stickerMessage"}
	for _, mt := range mediaTypes {
		if media, ok := message[mt].(map[string]interface{}); ok {
			if mime, ok := media["mimetype"].(string); ok {
				return mime
			}
		}
	}
	return ""
}

// extractCaption extracts caption from media messages
func extractCaption(message map[string]interface{}) string {
	captionTypes := []string{"imageMessage", "videoMessage", "documentMessage"}
	for _, mt := range captionTypes {
		if media, ok := message[mt].(map[string]interface{}); ok {
			if caption, ok := media["caption"].(string); ok {
				return caption
			}
		}
	}
	return ""
}

// extractFileName extracts filename from document messages
func extractFileName(message map[string]interface{}) string {
	if doc, ok := message["documentMessage"].(map[string]interface{}); ok {
		if fileName, ok := doc["fileName"].(string); ok {
			return fileName
		}
		if title, ok := doc["title"].(string); ok {
			return title
		}
	}
	return ""
}

// getOrCreateEvolutionContact finds or creates a contact for Evolution messages
func (a *App) getOrCreateEvolutionContact(instance *models.EvolutionInstance, msg *models.EvolutionMessagePayload) (*models.Contact, error) {
	phone := msg.GetSenderPhone()
	if phone == "" {
		return nil, fmt.Errorf("invalid sender phone")
	}

	// Format phone number
	formattedPhone := formatPhone(phone)

	// Try to find existing contact
	var contact models.Contact
	err := a.DB.Where("phone_number = ? AND organization_id = ?", formattedPhone, instance.OrganizationID).First(&contact).Error

	if err == nil {
		// Update last inbound time
		now := time.Now()
		a.DB.Model(&contact).Update("last_inbound_at", &now)
		return &contact, nil
	}

	// Create new contact
	contact = models.Contact{
		BaseModel:        models.BaseModel{ID: uuid.New()},
		OrganizationID:   instance.OrganizationID,
		PhoneNumber:      formattedPhone,
		Name:             msg.PushName,
		WhatsAppAccount:  instance.InstanceName,
		LastInboundAt:    func() *time.Time { t := time.Now(); return &t }(),
		AssignmentStatus: "unassigned",
	}

	if err := a.DB.Create(&contact).Error; err != nil {
		return nil, err
	}

	return &contact, nil
}

// formatPhone formats phone number to standard format
func formatPhone(phone string) string {
	// Remove any non-numeric characters
	cleaned := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	// Add country code if needed (assuming Brazil +55)
	if len(cleaned) == 11 && strings.HasPrefix(cleaned, "9") {
		cleaned = "55" + cleaned
	} else if len(cleaned) == 13 && !strings.HasPrefix(cleaned, "55") {
		// Keep as is
	}

	return cleaned
}

// processIncomingEvolutionMessage processes and saves an incoming message
func (a *App) processIncomingEvolutionMessage(ctx context.Context, instance *models.EvolutionInstance, contact *models.Contact, msg *models.EvolutionMessagePayload) error {
	// Determine message type
	msgType := msg.ToWhatsAppMessageType()
	textContent := msg.ExtractText()

	// For media messages, text might be in caption
	if textContent == "" && msgType != models.MessageTypeText {
		textContent = msg.ExtractCaption()
	}

	// Download media if present
	var localMediaPath string
	mediaURL := msg.ExtractMediaURL()
	if mediaURL != "" && msgType != models.MessageTypeText {
		// Download and save media
		mimeType := msg.MimeType
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		var err error
		localMediaPath, err = a.DownloadAndSaveMedia(ctx, mediaURL, mimeType, nil)
		if err != nil {
			a.Log.Error("Failed to download media", "error", err, "url", mediaURL)
			// Continue without media
		}
	}

	// Create message record
	message := models.Message{
		BaseModel:       models.BaseModel{ID: uuid.New()},
		OrganizationID:  contact.OrganizationID,
		ContactID:       contact.ID,
		Direction:       models.DirectionInbound,
		MessageType:     msgType,
		Content:         textContent,
		MediaURL:        localMediaPath,
		MimeType:        msg.MimeType,
		FileName:        msg.FileName,
		Status:          models.MessageStatusDelivered,
		ExternalID:      msg.Key.ID,
		WhatsAppAccount: instance.InstanceName,
		SentAt:          func() *time.Time { t := time.Unix(msg.Timestamp, 0); return &t }(),
	}

	if err := a.DB.Create(&message).Error; err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	// Update contact's last message
	a.DB.Model(contact).Update("last_message_at", message.SentAt)

	// Broadcast message via WebSocket
	a.broadcastEvolutionMessage(contact, &message)

	// Trigger chatbot processing if enabled
	go a.processChatbotForEvolution(context.Background(), contact, &message)

	a.Log.Info("Evolution message processed",
		"contact", contact.PhoneNumber,
		"type", msgType,
		"hasMedia", localMediaPath != "")

	return nil
}

// broadcastEvolutionStatus broadcasts instance status via WebSocket
func (a *App) broadcastEvolutionStatus(instance *models.EvolutionInstance) {
	// Implement WebSocket broadcast
	a.Log.Debug("Broadcasting evolution status", "instance", instance.InstanceName, "status", instance.Status)
}

// broadcastEvolutionMessage broadcasts new Evolution message via WebSocket
func (a *App) broadcastEvolutionMessage(contact *models.Contact, message *models.Message) {
	// Use the existing broadcast function
	a.broadcastNewMessage(contact.OrganizationID, message, contact)
}

// processChatbotForEvolution triggers chatbot processing for Evolution messages
func (a *App) processChatbotForEvolution(ctx context.Context, contact *models.Contact, message *models.Message) {
	// This would integrate with the existing chatbot processor
	a.Log.Debug("Chatbot processing for Evolution", "contact", contact.PhoneNumber)
}
