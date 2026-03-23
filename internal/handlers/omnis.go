package handlers

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shridarpatil/whatomate/internal/models"
	"github.com/valyala/fasthttp"
	"github.com/zerodha/fastglue"
)

// ListOmnis Returns all categories and nested scripts
func (a *App) ListOmnis(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	var categories []models.OmniCategory
	if err := a.DB.Preload("Scripts").Where("organization_id = ?", orgID).Order("created_at asc").Find(&categories).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to load omnis", nil, "")
	}

	return r.SendEnvelope(categories)
}

type SendOmniRequest struct {
	ContactID string `json:"contact_id"`
	ScriptID  string `json:"script_id"`
}

// SendOmni executes an omni script immediately
func (a *App) SendOmni(r *fastglue.Request) error {
	orgID, userID, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	var req SendOmniRequest
	if err := a.decodeRequest(r, &req); err != nil {
		return nil
	}

	scriptUUID, err := uuid.Parse(req.ScriptID)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid script ID", nil, "")
	}

	contactUUID, err := uuid.Parse(req.ContactID)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid contact ID", nil, "")
	}

	var script models.OmniScript
	if err := a.DB.Where("id = ? AND organization_id = ?", scriptUUID, orgID).First(&script).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Script not found", nil, "")
	}

	var contact models.Contact
	if err := a.DB.Where("id = ? AND organization_id = ?", contactUUID, orgID).First(&contact).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Contact not found", nil, "")
	}

	// 24h compliance check
	if contact.LastInboundAt != nil {
		hours := time.Since(*contact.LastInboundAt).Hours()
		if hours > 24 {
			return r.SendErrorEnvelope(fasthttp.StatusForbidden, "Lead is outside the 24-hour active window. Cannot send message.", nil, "")
		}
	} else {
	    // If they never sent a message, we might not be able to send them anything except templates.
		return r.SendErrorEnvelope(fasthttp.StatusForbidden, "No inbound message detected. Template required.", nil, "")
	}

	account, err := a.resolveWhatsAppAccount(orgID, contact.WhatsAppAccount)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "WhatsApp Account error", nil, "")
	}

	var caption string
	var mediaType models.MessageType

	switch script.MediaType {
	case "text":
		mediaType = models.MessageTypeText
	case "audio", "voice":
		mediaType = models.MessageTypeAudio
	case "image":
		mediaType = models.MessageTypeImage
		caption = script.Content
	case "video":
		mediaType = models.MessageTypeVideo
		caption = script.Content
	case "document":
		mediaType = models.MessageTypeDocument
		caption = script.Content
	default:
		mediaType = models.MessageTypeText
	}

	reqData := OutgoingMessageRequest{
		Account: account,
		Contact: &contact,
		Type:    mediaType,
		Caption: caption,
	}

	if mediaType == models.MessageTypeText {
		reqData.Content = script.Content
	} else {
		localPath := script.Content
		if strings.HasPrefix(localPath, "/") {
			localPath = localPath[1:]
		}
		
		mediaBytes, err := os.ReadFile(localPath)
		if err == nil {
			reqData.MediaData = mediaBytes
		}
	}

	opts := DefaultSendOptions()
	opts.SentByUserID = &userID

	if script.DelayMs > 0 {
		time.Sleep(time.Duration(script.DelayMs) * time.Millisecond)
	}

	ctx := context.Background()
	message, err := a.SendOutgoingMessage(ctx, reqData, opts)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to send omni message to WhatsApp", nil, "")
	}

	return r.SendEnvelope(message)
}
