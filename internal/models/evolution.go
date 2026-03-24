package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// EvolutionInstance represents a WhatsApp instance managed by Evolution API
type EvolutionInstance struct {
	BaseModel
	OrganizationID uuid.UUID      `gorm:"type:uuid;not null;index" json:"organization_id"`
	InstanceName   string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"instance_name"`
	DisplayName    string         `gorm:"type:varchar(255)" json:"display_name"`
	Phone          string         `gorm:"type:varchar(20)" json:"phone"`
	Status         string         `gorm:"type:varchar(50);default:'disconnected'" json:"status"` // connected, disconnected, connecting
	ProfileName    sql.NullString `gorm:"type:varchar(255)" json:"profile_name"`
	ProfilePic     sql.NullString `gorm:"type:text" json:"profile_pic"`
	WebhookURL     string         `gorm:"type:text" json:"webhook_url"`
	LastError      sql.NullString `gorm:"type:text" json:"last_error"`
	ConnectedAt    *time.Time     `json:"connected_at"`
	LastActivity   *time.Time     `json:"last_activity"`
}

// EvolutionWebhookPayload represents incoming webhook data from Evolution
type EvolutionWebhookPayload struct {
	Event    string                 `json:"event"`
	Instance string                 `json:"instance"`
	Data     map[string]interface{} `json:"data"`
}

// EvolutionMessagePayload represents a message received via webhook
type EvolutionMessagePayload struct {
	Key struct {
		RemoteJid string `json:"remoteJid"`
		FromMe    bool   `json:"fromMe"`
		ID        string `json:"id"`
	} `json:"key"`
	Message      map[string]interface{} `json:"message"`
	MessageType  string                 `json:"messageType"`
	Timestamp    int64                  `json:"timestamp"`
	PushName     string                 `json:"pushName"`
	Status       string                 `json:"status"`
	Conversation string                 `json:"conversation,omitempty"`
	Caption      string                 `json:"caption,omitempty"`
	MediaURL     string                 `json:"mediaUrl,omitempty"`
	MimeType     string                 `json:"mimeType,omitempty"`
	FileName     string                 `json:"fileName,omitempty"`
}

// EvolutionConnectionPayload represents connection status update
type EvolutionConnectionPayload struct {
	State string `json:"state"` // open, close, connecting
}

// EvolutionStatus constants
const (
	EvolutionStatusConnected    = "connected"
	EvolutionStatusDisconnected = "disconnected"
	EvolutionStatusConnecting   = "connecting"
)

// EvolutionEvent constants
const (
	EvolutionEventMessageUpsert    = "messages.upsert"
	EvolutionEventConnectionUpdate = "connection.update"
	EvolutionEventQrCode           = "qrcode"
	EvolutionEventInstanceCreated  = "instance.created"
)

// ToWhatsAppMessageType converts Evolution message type to internal MessageType
func (p *EvolutionMessagePayload) ToWhatsAppMessageType() MessageType {
	switch p.MessageType {
	case "conversation", "extendedTextMessage":
		return MessageTypeText
	case "imageMessage":
		return MessageTypeImage
	case "videoMessage":
		return MessageTypeVideo
	case "audioMessage":
		return MessageTypeAudio
	case "documentMessage":
		return MessageTypeDocument
	case "stickerMessage":
		return MessageTypeImage
	default:
		return MessageTypeText
	}
}

// ExtractText extracts text content from message
func (p *EvolutionMessagePayload) ExtractText() string {
	if p.Conversation != "" {
		return p.Conversation
	}

	// Try to extract from extendedTextMessage
	if extText, ok := p.Message["extendedTextMessage"].(map[string]interface{}); ok {
		if text, ok := extText["text"].(string); ok {
			return text
		}
	}

	return ""
}

// ExtractMediaURL extracts media URL from message
func (p *EvolutionMessagePayload) ExtractMediaURL() string {
	if p.MediaURL != "" {
		return p.MediaURL
	}

	// Try to extract from different message types
	mediaTypes := []string{"imageMessage", "videoMessage", "audioMessage", "documentMessage"}
	for _, mt := range mediaTypes {
		if media, ok := p.Message[mt].(map[string]interface{}); ok {
			if url, ok := media["url"].(string); ok {
				return url
			}
		}
	}

	return ""
}

// ExtractCaption extracts caption from media messages
func (p *EvolutionMessagePayload) ExtractCaption() string {
	if p.Caption != "" {
		return p.Caption
	}

	mediaTypes := []string{"imageMessage", "videoMessage", "documentMessage"}
	for _, mt := range mediaTypes {
		if media, ok := p.Message[mt].(map[string]interface{}); ok {
			if caption, ok := media["caption"].(string); ok {
				return caption
			}
		}
	}

	return ""
}

// ExtractFileName extracts filename from document messages
func (p *EvolutionMessagePayload) ExtractFileName() string {
	if p.FileName != "" {
		return p.FileName
	}

	if doc, ok := p.Message["documentMessage"].(map[string]interface{}); ok {
		if title, ok := doc["title"].(string); ok {
			return title
		}
		if fileName, ok := doc["fileName"].(string); ok {
			return fileName
		}
	}

	return ""
}

// IsFromMe checks if message is sent by the instance owner
func (p *EvolutionMessagePayload) IsFromMe() bool {
	return p.Key.FromMe
}

// GetSenderPhone extracts phone number from remoteJid
func (p *EvolutionMessagePayload) GetSenderPhone() string {
	// remoteJid format: "5511999999999@s.whatsapp.net"
	jid := p.Key.RemoteJid
	if len(jid) > 0 {
		// Remove @s.whatsapp.net suffix
		for i, c := range jid {
			if c == '@' {
				return jid[:i]
			}
		}
		return jid
	}
	return ""
}
