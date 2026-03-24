package handlers

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/omni-platform/omni/internal/models"
	"github.com/valyala/fasthttp"
	"github.com/zerodha/fastglue"
	"gorm.io/gorm"
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

// CreateCategoryRequest represents request body for creating a category
type CreateCategoryRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// CreateCategory creates a new omni category
func (a *App) CreateCategory(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	var req CreateCategoryRequest
	if err := a.decodeRequest(r, &req); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid request body", nil, "")
	}

	if req.Name == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Category name is required", nil, "")
	}

	if req.Color == "" {
		req.Color = "#3b82f6" // Default blue color
	}

	category := models.OmniCategory{
		BaseModel:      models.BaseModel{ID: uuid.New()},
		OrganizationID: orgID,
		Name:           req.Name,
		Color:          req.Color,
	}

	if err := a.DB.Create(&category).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create category", nil, "")
	}

	return r.SendEnvelope(category)
}

// UpdateCategoryRequest represents request body for updating a category
type UpdateCategoryRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// UpdateCategory updates an existing omni category
func (a *App) UpdateCategory(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	categoryID, ok := r.RequestCtx.UserValue("id").(string)
	if !ok || categoryID == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Category ID is required", nil, "")
	}

	var req UpdateCategoryRequest
	if err := a.decodeRequest(r, &req); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid request body", nil, "")
	}

	var category models.OmniCategory
	if err := a.DB.Where("id = ? AND organization_id = ?", categoryID, orgID).First(&category).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Category not found", nil, "")
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Color != "" {
		updates["color"] = req.Color
	}

	if len(updates) > 0 {
		if err := a.DB.Model(&category).Updates(updates).Error; err != nil {
			return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to update category", nil, "")
		}
	}

	return r.SendEnvelope(category)
}

// DeleteCategory deletes an omni category and all its scripts
func (a *App) DeleteCategory(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	categoryID, ok := r.RequestCtx.UserValue("id").(string)
	if !ok || categoryID == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Category ID is required", nil, "")
	}

	// First delete all scripts in the category
	if err := a.DB.Where("category_id = ? AND organization_id = ?", categoryID, orgID).Delete(&models.OmniScript{}).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to delete scripts", nil, "")
	}

	// Then delete the category
	if err := a.DB.Where("id = ? AND organization_id = ?", categoryID, orgID).Delete(&models.OmniCategory{}).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to delete category", nil, "")
	}

	return r.SendEnvelope(map[string]string{"status": "deleted"})
}

// CreateScriptRequest represents request body for creating a script
type CreateScriptRequest struct {
	CategoryID string `json:"category_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	MediaType  string `json:"media_type"`
	DelayMs    int    `json:"delay_ms"`
}

// CreateScript creates a new omni script
func (a *App) CreateScript(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	var req CreateScriptRequest
	if err := a.decodeRequest(r, &req); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid request body", nil, "")
	}

	if req.Title == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Script title is required", nil, "")
	}

	if req.CategoryID == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Category ID is required", nil, "")
	}

	// Verify category exists
	var category models.OmniCategory
	if err := a.DB.Where("id = ? AND organization_id = ?", req.CategoryID, orgID).First(&category).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Category not found", nil, "")
	}

	// Default media type to text
	if req.MediaType == "" {
		req.MediaType = "text"
	}

	script := models.OmniScript{
		BaseModel:      models.BaseModel{ID: uuid.New()},
		OrganizationID: orgID,
		CategoryID:     uuid.MustParse(req.CategoryID),
		Title:          req.Title,
		Content:        req.Content,
		MediaType:      req.MediaType,
		DelayMs:        req.DelayMs,
	}

	if err := a.DB.Create(&script).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create script", nil, "")
	}

	return r.SendEnvelope(script)
}

// UpdateScriptRequest represents request body for updating a script
type UpdateScriptRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	MediaType string `json:"media_type"`
	DelayMs   int    `json:"delay_ms"`
}

// UpdateScript updates an existing omni script
func (a *App) UpdateScript(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	scriptID, ok := r.RequestCtx.UserValue("id").(string)
	if !ok || scriptID == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Script ID is required", nil, "")
	}

	var req UpdateScriptRequest
	if err := a.decodeRequest(r, &req); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid request body", nil, "")
	}

	var script models.OmniScript
	if err := a.DB.Where("id = ? AND organization_id = ?", scriptID, orgID).First(&script).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Script not found", nil, "")
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.MediaType != "" {
		updates["media_type"] = req.MediaType
	}
	if req.DelayMs >= 0 {
		updates["delay_ms"] = req.DelayMs
	}

	if len(updates) > 0 {
		if err := a.DB.Model(&script).Updates(updates).Error; err != nil {
			return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to update script", nil, "")
		}
	}

	return r.SendEnvelope(script)
}

// DeleteScript deletes an omni script
func (a *App) DeleteScript(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	scriptID, ok := r.RequestCtx.UserValue("id").(string)
	if !ok || scriptID == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Script ID is required", nil, "")
	}

	// If script has a media file, try to delete it
	var script models.OmniScript
	if err := a.DB.Where("id = ? AND organization_id = ?", scriptID, orgID).First(&script).Error; err == nil {
		if script.FileName != nil && script.Content != "" {
			// Try to delete the media file (ignore errors)
			os.Remove(script.Content)
		}
	}

	if err := a.DB.Where("id = ? AND organization_id = ?", scriptID, orgID).Delete(&models.OmniScript{}).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to delete script", nil, "")
	}

	return r.SendEnvelope(map[string]string{"status": "deleted"})
}

// UploadMediaRequest handles media upload for omni scripts
func (a *App) UploadOmniMedia(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	scriptID, ok := r.RequestCtx.UserValue("id").(string)
	if !ok || scriptID == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Script ID is required", nil, "")
	}

	// Verify script exists
	var script models.OmniScript
	if err := a.DB.Where("id = ? AND organization_id = ?", scriptID, orgID).First(&script).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Script not found", nil, "")
	}

	fileHeader, err := r.RequestCtx.FormFile("file")
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "No file provided", nil, "")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to open uploaded file", nil, "")
	}
	defer file.Close()

	// Read file content
	mediaBytes := make([]byte, fileHeader.Size)
	_, err = file.Read(mediaBytes)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to read file", nil, "")
	}

	// Determine MIME type from filename
	mimeType := "application/octet-stream"
	fileName := fileHeader.Filename
	if strings.HasSuffix(fileName, ".mp4") {
		mimeType = "video/mp4"
	} else if strings.HasSuffix(fileName, ".ogg") || strings.HasSuffix(fileName, ".mp3") {
		mimeType = "audio/mpeg"
	} else if strings.HasSuffix(fileName, ".png") {
		mimeType = "image/png"
	} else if strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".jpeg") {
		mimeType = "image/jpeg"
	} else if strings.HasSuffix(fileName, ".pdf") {
		mimeType = "application/pdf"
	}

	// Save media locally
	localPath, err := a.saveMediaLocally(mediaBytes, mimeType, fileName)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to save media", nil, "")
	}

	// Update script with media info
	updates := map[string]interface{}{
		"content":   localPath,
		"file_name": fileName,
	}
	if err := a.DB.Model(&script).Updates(updates).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to update script", nil, "")
	}

	return r.SendEnvelope(map[string]interface{}{
		"file_name":  fileName,
		"local_path": localPath,
	})
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

// ListSequences returns all sequences with their steps
func (a *App) ListSequences(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	var sequences []models.OmniSequence
	if err := a.DB.Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_order ASC")
	}).Where("organization_id = ?", orgID).Order("created_at asc").Find(&sequences).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to load sequences", nil, "")
	}

	return r.SendEnvelope(sequences)
}

// CreateSequenceRequest represents request body for creating a sequence
type CreateSequenceRequest struct {
	CategoryID  string                      `json:"category_id"`
	Title       string                      `json:"title"`
	Description string                      `json:"description"`
	Steps       []CreateSequenceStepRequest `json:"steps"`
}

// CreateSequenceStepRequest represents a step in sequence creation
type CreateSequenceStepRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	MediaType string `json:"media_type"`
	DelayMs   int    `json:"delay_ms"`
}

// CreateSequence creates a new omni sequence with steps
func (a *App) CreateSequence(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	var req CreateSequenceRequest
	if err := a.decodeRequest(r, &req); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid request body", nil, "")
	}

	if req.Title == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Sequence title is required", nil, "")
	}

	if req.CategoryID == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Category ID is required", nil, "")
	}

	// Verify category exists
	var category models.OmniCategory
	if err := a.DB.Where("id = ? AND organization_id = ?", req.CategoryID, orgID).First(&category).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Category not found", nil, "")
	}

	tx := a.DB.Begin()

	sequence := models.OmniSequence{
		BaseModel:      models.BaseModel{ID: uuid.New()},
		OrganizationID: orgID,
		CategoryID:     uuid.MustParse(req.CategoryID),
		Title:          req.Title,
		Description:    req.Description,
	}

	if err := tx.Create(&sequence).Error; err != nil {
		tx.Rollback()
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create sequence", nil, "")
	}

	// Create steps
	for i, stepReq := range req.Steps {
		step := models.OmniSequenceStep{
			BaseModel:      models.BaseModel{ID: uuid.New()},
			OrganizationID: orgID,
			SequenceID:     sequence.ID,
			StepOrder:      i + 1,
			Title:          stepReq.Title,
			Content:        stepReq.Content,
			MediaType:      stepReq.MediaType,
			DelayMs:        stepReq.DelayMs,
		}

		if step.MediaType == "" {
			step.MediaType = "text"
		}

		if err := tx.Create(&step).Error; err != nil {
			tx.Rollback()
			return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create sequence step", nil, "")
		}
	}

	tx.Commit()

	// Reload with steps
	a.DB.Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_order ASC")
	}).First(&sequence, sequence.ID)

	return r.SendEnvelope(sequence)
}

// UpdateSequenceRequest represents request body for updating a sequence
type UpdateSequenceRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateSequence updates an existing omni sequence
func (a *App) UpdateSequence(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	sequenceID, ok := r.RequestCtx.UserValue("id").(string)
	if !ok || sequenceID == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Sequence ID is required", nil, "")
	}

	var req UpdateSequenceRequest
	if err := a.decodeRequest(r, &req); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid request body", nil, "")
	}

	var sequence models.OmniSequence
	if err := a.DB.Where("id = ? AND organization_id = ?", sequenceID, orgID).First(&sequence).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Sequence not found", nil, "")
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	updates["description"] = req.Description

	if len(updates) > 0 {
		if err := a.DB.Model(&sequence).Updates(updates).Error; err != nil {
			return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to update sequence", nil, "")
		}
	}

	return r.SendEnvelope(sequence)
}

// DeleteSequence deletes an omni sequence and all its steps
func (a *App) DeleteSequence(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	sequenceID, ok := r.RequestCtx.UserValue("id").(string)
	if !ok || sequenceID == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Sequence ID is required", nil, "")
	}

	// First delete all steps
	if err := a.DB.Where("sequence_id = ? AND organization_id = ?", sequenceID, orgID).Delete(&models.OmniSequenceStep{}).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to delete sequence steps", nil, "")
	}

	// Then delete the sequence
	if err := a.DB.Where("id = ? AND organization_id = ?", sequenceID, orgID).Delete(&models.OmniSequence{}).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to delete sequence", nil, "")
	}

	return r.SendEnvelope(map[string]string{"status": "deleted"})
}

// AddSequenceStepRequest represents request body for adding a step
type AddSequenceStepRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	MediaType string `json:"media_type"`
	DelayMs   int    `json:"delay_ms"`
}

// AddSequenceStep adds a new step to an existing sequence
func (a *App) AddSequenceStep(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	sequenceID, ok := r.RequestCtx.UserValue("id").(string)
	if !ok || sequenceID == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Sequence ID is required", nil, "")
	}

	var req AddSequenceStepRequest
	if err := a.decodeRequest(r, &req); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid request body", nil, "")
	}

	// Verify sequence exists
	var sequence models.OmniSequence
	if err := a.DB.Where("id = ? AND organization_id = ?", sequenceID, orgID).First(&sequence).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Sequence not found", nil, "")
	}

	// Get max step order
	var maxOrder int
	a.DB.Model(&models.OmniSequenceStep{}).
		Where("sequence_id = ?", sequenceID).
		Select("COALESCE(MAX(step_order), 0)").
		Scan(&maxOrder)

	step := models.OmniSequenceStep{
		BaseModel:      models.BaseModel{ID: uuid.New()},
		OrganizationID: orgID,
		SequenceID:     uuid.MustParse(sequenceID),
		StepOrder:      maxOrder + 1,
		Title:          req.Title,
		Content:        req.Content,
		MediaType:      req.MediaType,
		DelayMs:        req.DelayMs,
	}

	if step.MediaType == "" {
		step.MediaType = "text"
	}

	if err := a.DB.Create(&step).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create step", nil, "")
	}

	return r.SendEnvelope(step)
}

// DeleteSequenceStep deletes a step from a sequence
func (a *App) DeleteSequenceStep(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	stepID, ok := r.RequestCtx.UserValue("step_id").(string)
	if !ok || stepID == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Step ID is required", nil, "")
	}

	// Verify step exists and belongs to organization
	var step models.OmniSequenceStep
	if err := a.DB.Where("id = ? AND organization_id = ?", stepID, orgID).First(&step).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Step not found", nil, "")
	}

	if err := a.DB.Delete(&step).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to delete step", nil, "")
	}

	return r.SendEnvelope(map[string]string{"status": "deleted"})
}

// SendSequenceRequest represents request body for sending a sequence
type SendSequenceRequest struct {
	ContactID  string `json:"contact_id"`
	SequenceID string `json:"sequence_id"`
}

// SendSequence executes an omni sequence (sends all steps with delays)
func (a *App) SendSequence(r *fastglue.Request) error {
	orgID, userID, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	var req SendSequenceRequest
	if err := a.decodeRequest(r, &req); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid request body", nil, "")
	}

	sequenceUUID, err := uuid.Parse(req.SequenceID)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid sequence ID", nil, "")
	}

	contactUUID, err := uuid.Parse(req.ContactID)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid contact ID", nil, "")
	}

	// Load sequence with steps
	var sequence models.OmniSequence
	if err := a.DB.Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_order ASC")
	}).Where("id = ? AND organization_id = ?", sequenceUUID, orgID).First(&sequence).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "Sequence not found", nil, "")
	}

	if len(sequence.Steps) == 0 {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Sequence has no steps", nil, "")
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
		return r.SendErrorEnvelope(fasthttp.StatusForbidden, "No inbound message detected. Template required.", nil, "")
	}

	account, err := a.resolveWhatsAppAccount(orgID, contact.WhatsAppAccount)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "WhatsApp Account error", nil, "")
	}

	ctx := context.Background()
	opts := DefaultSendOptions()
	opts.SentByUserID = &userID

	// Send each step with delays
	sentMessages := []interface{}{}
	for _, step := range sequence.Steps {
		// Apply step delay
		if step.DelayMs > 0 {
			time.Sleep(time.Duration(step.DelayMs) * time.Millisecond)
		}

		var caption string
		var mediaType models.MessageType

		switch step.MediaType {
		case "text":
			mediaType = models.MessageTypeText
		case "audio", "voice":
			mediaType = models.MessageTypeAudio
		case "image":
			mediaType = models.MessageTypeImage
			caption = step.Content
		case "video":
			mediaType = models.MessageTypeVideo
			caption = step.Content
		case "document":
			mediaType = models.MessageTypeDocument
			caption = step.Content
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
			reqData.Content = step.Content
		} else {
			localPath := step.Content
			if strings.HasPrefix(localPath, "/") {
				localPath = localPath[1:]
			}

			mediaBytes, err := os.ReadFile(localPath)
			if err == nil {
				reqData.MediaData = mediaBytes
			}
		}

		message, err := a.SendOutgoingMessage(ctx, reqData, opts)
		if err != nil {
			// Log error but continue with next steps
			a.Log.Error("Failed to send sequence step", "error", err, "step", step.StepOrder)
			continue
		}

		sentMessages = append(sentMessages, message)
	}

	return r.SendEnvelope(map[string]interface{}{
		"status":         "sent",
		"messages_count": len(sentMessages),
		"messages":       sentMessages,
	})
}

// ExportOmnisRequest represents request body for exporting omnis
type ExportOmnisRequest struct {
	IncludeMedia bool `json:"include_media"`
}

// OmniExportData represents the complete export structure
type OmniExportData struct {
	Version    string               `json:"version"`
	ExportedAt string               `json:"exported_at"`
	Categories []OmniCategoryExport `json:"categories"`
	Sequences  []OmniSequenceExport `json:"sequences"`
	MediaFiles map[string][]byte    `json:"media_files,omitempty"`
}

// OmniCategoryExport represents a category with its scripts for export
type OmniCategoryExport struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Color   string             `json:"color"`
	Scripts []OmniScriptExport `json:"scripts"`
}

// OmniScriptExport represents a script for export
type OmniScriptExport struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MediaType   string `json:"media_type"`
	FileName    string `json:"file_name,omitempty"`
	DelayMs     int    `json:"delay_ms"`
	MediaBase64 string `json:"media_base64,omitempty"`
}

// OmniSequenceExport represents a sequence with its steps for export
type OmniSequenceExport struct {
	ID          string                   `json:"id"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	Steps       []OmniSequenceStepExport `json:"steps"`
}

// OmniSequenceStepExport represents a sequence step for export
type OmniSequenceStepExport struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	MediaType   string `json:"media_type"`
	FileName    string `json:"file_name,omitempty"`
	DelayMs     int    `json:"delay_ms"`
	MediaBase64 string `json:"media_base64,omitempty"`
}

// ExportOmnis exports all omnis (categories, scripts, sequences) as JSON
func (a *App) ExportOmnis(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	var req ExportOmnisRequest
	// Try to decode request body, but if it fails, use defaults
	_ = a.decodeRequest(r, &req)

	// Load categories with scripts
	var categories []models.OmniCategory
	if err := a.DB.Preload("Scripts").Where("organization_id = ?", orgID).Order("created_at asc").Find(&categories).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to load categories", nil, "")
	}

	// Load sequences with steps
	var sequences []models.OmniSequence
	if err := a.DB.Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_order ASC")
	}).Where("organization_id = ?", orgID).Order("created_at asc").Find(&sequences).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to load sequences", nil, "")
	}

	exportData := OmniExportData{
		Version:    "1.0",
		ExportedAt: time.Now().UTC().Format(time.RFC3339),
		Categories: make([]OmniCategoryExport, 0, len(categories)),
		Sequences:  make([]OmniSequenceExport, 0, len(sequences)),
		MediaFiles: make(map[string][]byte),
	}

	// Export categories and scripts
	for _, cat := range categories {
		catExport := OmniCategoryExport{
			ID:      cat.ID.String(),
			Name:    cat.Name,
			Color:   cat.Color,
			Scripts: make([]OmniScriptExport, 0, len(cat.Scripts)),
		}

		for _, script := range cat.Scripts {
			scriptExport := OmniScriptExport{
				ID:        script.ID.String(),
				Title:     script.Title,
				Content:   script.Content,
				MediaType: script.MediaType,
				DelayMs:   script.DelayMs,
			}

			if script.FileName != nil {
				scriptExport.FileName = *script.FileName
			}

			// Include media as base64 if requested
			if req.IncludeMedia && script.MediaType != "text" && script.Content != "" {
				localPath := script.Content
				if strings.HasPrefix(localPath, "/") {
					localPath = localPath[1:]
				}
				if mediaBytes, err := os.ReadFile(localPath); err == nil {
					scriptExport.MediaBase64 = base64.StdEncoding.EncodeToString(mediaBytes)
					// Store in media files map for ZIP export
					if script.FileName != nil {
						exportData.MediaFiles[*script.FileName] = mediaBytes
					}
				}
			}

			catExport.Scripts = append(catExport.Scripts, scriptExport)
		}

		exportData.Categories = append(exportData.Categories, catExport)
	}

	// Export sequences
	for _, seq := range sequences {
		seqExport := OmniSequenceExport{
			ID:          seq.ID.String(),
			Title:       seq.Title,
			Description: seq.Description,
			Steps:       make([]OmniSequenceStepExport, 0, len(seq.Steps)),
		}

		for _, step := range seq.Steps {
			stepExport := OmniSequenceStepExport{
				Title:     step.Title,
				Content:   step.Content,
				MediaType: step.MediaType,
				DelayMs:   step.DelayMs,
			}

			if step.FileName != nil {
				stepExport.FileName = *step.FileName
			}

			// Include media as base64 if requested
			if req.IncludeMedia && step.MediaType != "text" && step.Content != "" {
				localPath := step.Content
				if strings.HasPrefix(localPath, "/") {
					localPath = localPath[1:]
				}
				if mediaBytes, err := os.ReadFile(localPath); err == nil {
					stepExport.MediaBase64 = base64.StdEncoding.EncodeToString(mediaBytes)
					// Store in media files map for ZIP export
					if step.FileName != nil {
						exportData.MediaFiles[*step.FileName] = mediaBytes
					}
				}
			}

			seqExport.Steps = append(seqExport.Steps, stepExport)
		}

		exportData.Sequences = append(exportData.Sequences, seqExport)
	}

	// Return JSON export
	return r.SendEnvelope(exportData)
}

// ExportOmnisZIP exports all omnis as a ZIP file containing JSON and media files
func (a *App) ExportOmnisZIP(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	// Load categories with scripts
	var categories []models.OmniCategory
	if err := a.DB.Preload("Scripts").Where("organization_id = ?", orgID).Order("created_at asc").Find(&categories).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to load categories", nil, "")
	}

	// Load sequences with steps
	var sequences []models.OmniSequence
	if err := a.DB.Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_order ASC")
	}).Where("organization_id = ?", orgID).Order("created_at asc").Find(&sequences).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to load sequences", nil, "")
	}

	exportData := OmniExportData{
		Version:    "1.0",
		ExportedAt: time.Now().UTC().Format(time.RFC3339),
		Categories: make([]OmniCategoryExport, 0, len(categories)),
		Sequences:  make([]OmniSequenceExport, 0, len(sequences)),
		MediaFiles: make(map[string][]byte),
	}

	// Export categories and scripts with media
	for _, cat := range categories {
		catExport := OmniCategoryExport{
			ID:      cat.ID.String(),
			Name:    cat.Name,
			Color:   cat.Color,
			Scripts: make([]OmniScriptExport, 0, len(cat.Scripts)),
		}

		for _, script := range cat.Scripts {
			scriptExport := OmniScriptExport{
				ID:        script.ID.String(),
				Title:     script.Title,
				Content:   script.Content,
				MediaType: script.MediaType,
				DelayMs:   script.DelayMs,
			}

			if script.FileName != nil {
				scriptExport.FileName = *script.FileName
			}

			// Read media file for ZIP
			if script.MediaType != "text" && script.Content != "" {
				localPath := script.Content
				if strings.HasPrefix(localPath, "/") {
					localPath = localPath[1:]
				}
				if mediaBytes, err := os.ReadFile(localPath); err == nil {
					if script.FileName != nil {
						exportData.MediaFiles[*script.FileName] = mediaBytes
					}
				}
			}

			catExport.Scripts = append(catExport.Scripts, scriptExport)
		}

		exportData.Categories = append(exportData.Categories, catExport)
	}

	// Export sequences with media
	for _, seq := range sequences {
		seqExport := OmniSequenceExport{
			ID:          seq.ID.String(),
			Title:       seq.Title,
			Description: seq.Description,
			Steps:       make([]OmniSequenceStepExport, 0, len(seq.Steps)),
		}

		for _, step := range seq.Steps {
			stepExport := OmniSequenceStepExport{
				Title:     step.Title,
				Content:   step.Content,
				MediaType: step.MediaType,
				DelayMs:   step.DelayMs,
			}

			if step.FileName != nil {
				stepExport.FileName = *step.FileName
			}

			// Read media file for ZIP
			if step.MediaType != "text" && step.Content != "" {
				localPath := step.Content
				if strings.HasPrefix(localPath, "/") {
					localPath = localPath[1:]
				}
				if mediaBytes, err := os.ReadFile(localPath); err == nil {
					if step.FileName != nil {
						exportData.MediaFiles[*step.FileName] = mediaBytes
					}
				}
			}

			seqExport.Steps = append(seqExport.Steps, stepExport)
		}

		exportData.Sequences = append(exportData.Sequences, seqExport)
	}

	// Remove media files from JSON (they'll be in the ZIP)
	exportData.MediaFiles = nil

	// Create ZIP in memory
	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	// Add config.json to ZIP
	configJSON, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to marshal config", nil, "")
	}

	configFile, err := zipWriter.Create("config.json")
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create config file in ZIP", nil, "")
	}
	configFile.Write(configJSON)

	// Add media files to ZIP
	for filename, mediaData := range exportData.MediaFiles {
		mediaFile, err := zipWriter.Create("media/" + filename)
		if err != nil {
			continue // Skip files that fail
		}
		mediaFile.Write(mediaData)
	}

	zipWriter.Close()

	// Set headers for file download
	r.RequestCtx.SetContentType("application/zip")
	r.RequestCtx.Response.Header.Set("Content-Disposition", "attachment; filename=omnis-export.zip")

	// Return ZIP content directly
	r.RequestCtx.SetBody(zipBuffer.Bytes())
	return nil
}

// ImportOmnisRequest represents the import data structure
type ImportOmnisRequest struct {
	ReplaceExisting bool           `json:"replace_existing"`
	Data            OmniExportData `json:"data"`
}

// ImportOmnis imports omnis from JSON export data
func (a *App) ImportOmnis(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	var req ImportOmnisRequest
	if err := a.decodeRequest(r, &req); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid request body", nil, "")
	}

	if req.Data.Version == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid export data version", nil, "")
	}

	tx := a.DB.Begin()

	// Replace existing data if requested
	if req.ReplaceExisting {
		// Delete existing sequence steps
		tx.Exec("DELETE FROM omni_sequence_steps WHERE organization_id = ?", orgID)
		// Delete existing sequences
		tx.Where("organization_id = ?", orgID).Delete(&models.OmniSequence{})
		// Delete existing scripts
		tx.Where("organization_id = ?", orgID).Delete(&models.OmniScript{})
		// Delete existing categories
		tx.Where("organization_id = ?", orgID).Delete(&models.OmniCategory{})
	}

	// Create a map to track old category IDs to new ones
	categoryIDMap := make(map[string]uuid.UUID)
	importedCategories := 0
	importedScripts := 0
	importedSequences := 0
	importedSteps := 0

	// Import categories and scripts
	for _, catExport := range req.Data.Categories {
		oldCatID := catExport.ID
		newCatID := uuid.New()
		categoryIDMap[oldCatID] = newCatID

		category := models.OmniCategory{
			BaseModel:      models.BaseModel{ID: newCatID},
			OrganizationID: orgID,
			Name:           catExport.Name,
			Color:          catExport.Color,
		}

		if err := tx.Create(&category).Error; err != nil {
			tx.Rollback()
			return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create category", nil, "")
		}
		importedCategories++

		// Import scripts
		for _, scriptExport := range catExport.Scripts {
			content := scriptExport.Content

			// If media base64 is provided, save the file
			if scriptExport.MediaBase64 != "" && scriptExport.FileName != "" {
				mediaBytes, err := base64.StdEncoding.DecodeString(scriptExport.MediaBase64)
				if err == nil {
					mimeType := "application/octet-stream"
					if strings.HasSuffix(scriptExport.FileName, ".mp4") {
						mimeType = "video/mp4"
					} else if strings.HasSuffix(scriptExport.FileName, ".ogg") || strings.HasSuffix(scriptExport.FileName, ".mp3") {
						mimeType = "audio/mpeg"
					} else if strings.HasSuffix(scriptExport.FileName, ".png") {
						mimeType = "image/png"
					} else if strings.HasSuffix(scriptExport.FileName, ".jpg") || strings.HasSuffix(scriptExport.FileName, ".jpeg") {
						mimeType = "image/jpeg"
					}

					localPath, err := a.saveMediaLocally(mediaBytes, mimeType, scriptExport.FileName)
					if err == nil {
						content = localPath
					}
				}
			}

			script := models.OmniScript{
				BaseModel:      models.BaseModel{ID: uuid.New()},
				OrganizationID: orgID,
				CategoryID:     newCatID,
				Title:          scriptExport.Title,
				Content:        content,
				MediaType:      scriptExport.MediaType,
				DelayMs:        scriptExport.DelayMs,
			}

			if scriptExport.FileName != "" {
				script.FileName = &scriptExport.FileName
			}

			if err := tx.Create(&script).Error; err != nil {
				tx.Rollback()
				return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create script", nil, "")
			}
			importedScripts++
		}
	}

	// Import sequences
	for _, seqExport := range req.Data.Sequences {
		sequence := models.OmniSequence{
			BaseModel:      models.BaseModel{ID: uuid.New()},
			OrganizationID: orgID,
			Title:          seqExport.Title,
			Description:    seqExport.Description,
		}

		// Try to map to existing category if possible
		if len(req.Data.Categories) > 0 {
			// Use first category as default
			if firstCatID, ok := categoryIDMap[req.Data.Categories[0].ID]; ok {
				sequence.CategoryID = firstCatID
			} else {
				sequence.CategoryID = uuid.New()
			}
		} else {
			sequence.CategoryID = uuid.New()
		}

		if err := tx.Create(&sequence).Error; err != nil {
			tx.Rollback()
			return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create sequence", nil, "")
		}
		importedSequences++

		// Import steps
		for i, stepExport := range seqExport.Steps {
			content := stepExport.Content

			// If media base64 is provided, save the file
			if stepExport.MediaBase64 != "" && stepExport.FileName != "" {
				mediaBytes, err := base64.StdEncoding.DecodeString(stepExport.MediaBase64)
				if err == nil {
					mimeType := "application/octet-stream"
					if strings.HasSuffix(stepExport.FileName, ".mp4") {
						mimeType = "video/mp4"
					} else if strings.HasSuffix(stepExport.FileName, ".ogg") || strings.HasSuffix(stepExport.FileName, ".mp3") {
						mimeType = "audio/mpeg"
					} else if strings.HasSuffix(stepExport.FileName, ".png") {
						mimeType = "image/png"
					} else if strings.HasSuffix(stepExport.FileName, ".jpg") || strings.HasSuffix(stepExport.FileName, ".jpeg") {
						mimeType = "image/jpeg"
					}

					localPath, err := a.saveMediaLocally(mediaBytes, mimeType, stepExport.FileName)
					if err == nil {
						content = localPath
					}
				}
			}

			step := models.OmniSequenceStep{
				BaseModel:      models.BaseModel{ID: uuid.New()},
				OrganizationID: orgID,
				SequenceID:     sequence.ID,
				StepOrder:      i + 1,
				Title:          stepExport.Title,
				Content:        content,
				MediaType:      stepExport.MediaType,
				DelayMs:        stepExport.DelayMs,
			}

			if stepExport.FileName != "" {
				step.FileName = &stepExport.FileName
			}

			if err := tx.Create(&step).Error; err != nil {
				tx.Rollback()
				return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create sequence step", nil, "")
			}
			importedSteps++
		}
	}

	tx.Commit()

	return r.SendEnvelope(map[string]interface{}{
		"status":              "success",
		"imported_categories": importedCategories,
		"imported_scripts":    importedScripts,
		"imported_sequences":  importedSequences,
		"imported_steps":      importedSteps,
	})
}

// ImportOmnisFromZIP imports omnis from a ZIP file (similar to existing ImportOmnisZIP but supports sequences)
func (a *App) ImportOmnisFromZIP(r *fastglue.Request) error {
	orgID, _, err := a.getOrgAndUserID(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}

	fileHeader, err := r.RequestCtx.FormFile("file")
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "No file provided", nil, "")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to open uploaded file", nil, "")
	}
	defer file.Close()

	zipBytes, err := io.ReadAll(file)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to read zip bytes", nil, "")
	}

	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid zip archive", nil, "")
	}

	var configJsonFile *zip.File
	mediaFiles := make(map[string][]byte)

	for _, zf := range zipReader.File {
		if strings.Contains(zf.Name, "__MACOSX") || strings.Contains(zf.Name, ".DS_Store") {
			continue
		}

		baseName := filepath.Base(zf.Name)
		if baseName == "config.json" && configJsonFile == nil {
			configJsonFile = zf
		} else if strings.HasPrefix(zf.Name, "media/") {
			// Read media file
			if rc, err := zf.Open(); err == nil {
				if data, err := io.ReadAll(rc); err == nil {
					filename := strings.TrimPrefix(zf.Name, "media/")
					mediaFiles[filename] = data
				}
				rc.Close()
			}
		}
	}

	if configJsonFile == nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "config.json missing in zip", nil, "")
	}

	rc, err := configJsonFile.Open()
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to read config.json", nil, "")
	}
	defer rc.Close()

	configBytes, err := io.ReadAll(rc)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to read config.json contents", nil, "")
	}

	var importData OmniExportData
	if err := json.Unmarshal(configBytes, &importData); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid config.json format", nil, "")
	}

	// Use the ImportOmnis logic with the parsed data
	tx := a.DB.Begin()

	// Delete existing data
	tx.Unscoped().Where("organization_id = ?", orgID).Delete(&models.OmniSequenceStep{})
	tx.Unscoped().Where("organization_id = ?", orgID).Delete(&models.OmniSequence{})
	tx.Unscoped().Where("organization_id = ?", orgID).Delete(&models.OmniScript{})
	tx.Unscoped().Where("organization_id = ?", orgID).Delete(&models.OmniCategory{})

	importedCategories := 0
	importedScripts := 0
	importedSequences := 0
	importedSteps := 0

	// Import categories and scripts
	for _, catExport := range importData.Categories {
		category := models.OmniCategory{
			BaseModel:      models.BaseModel{ID: uuid.New()},
			OrganizationID: orgID,
			Name:           catExport.Name,
			Color:          catExport.Color,
		}

		if err := tx.Create(&category).Error; err != nil {
			tx.Rollback()
			return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create category", nil, "")
		}
		importedCategories++

		// Import scripts
		for _, scriptExport := range catExport.Scripts {
			content := scriptExport.Content

			// If media file exists in ZIP, save it
			if scriptExport.FileName != "" {
				if mediaData, ok := mediaFiles[scriptExport.FileName]; ok {
					mimeType := "application/octet-stream"
					if strings.HasSuffix(scriptExport.FileName, ".mp4") {
						mimeType = "video/mp4"
					} else if strings.HasSuffix(scriptExport.FileName, ".ogg") || strings.HasSuffix(scriptExport.FileName, ".mp3") {
						mimeType = "audio/mpeg"
					} else if strings.HasSuffix(scriptExport.FileName, ".png") {
						mimeType = "image/png"
					} else if strings.HasSuffix(scriptExport.FileName, ".jpg") || strings.HasSuffix(scriptExport.FileName, ".jpeg") {
						mimeType = "image/jpeg"
					}

					localPath, err := a.saveMediaLocally(mediaData, mimeType, scriptExport.FileName)
					if err == nil {
						content = localPath
					}
				}
			}

			script := models.OmniScript{
				BaseModel:      models.BaseModel{ID: uuid.New()},
				OrganizationID: orgID,
				CategoryID:     category.ID,
				Title:          scriptExport.Title,
				Content:        content,
				MediaType:      scriptExport.MediaType,
				DelayMs:        scriptExport.DelayMs,
			}

			if scriptExport.FileName != "" {
				script.FileName = &scriptExport.FileName
			}

			if err := tx.Create(&script).Error; err != nil {
				tx.Rollback()
				return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create script", nil, "")
			}
			importedScripts++
		}
	}

	// Import sequences
	for _, seqExport := range importData.Sequences {
		sequence := models.OmniSequence{
			BaseModel:      models.BaseModel{ID: uuid.New()},
			OrganizationID: orgID,
			Title:          seqExport.Title,
			Description:    seqExport.Description,
			CategoryID:     uuid.New(), // Create a default category reference
		}

		// Use first category if available
		if len(importData.Categories) > 0 {
			var firstCat models.OmniCategory
			if err := tx.Where("organization_id = ? AND name = ?", orgID, importData.Categories[0].Name).First(&firstCat).Error; err == nil {
				sequence.CategoryID = firstCat.ID
			}
		}

		if err := tx.Create(&sequence).Error; err != nil {
			tx.Rollback()
			return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create sequence", nil, "")
		}
		importedSequences++

		// Import steps
		for i, stepExport := range seqExport.Steps {
			content := stepExport.Content

			// If media file exists in ZIP, save it
			if stepExport.FileName != "" {
				if mediaData, ok := mediaFiles[stepExport.FileName]; ok {
					mimeType := "application/octet-stream"
					if strings.HasSuffix(stepExport.FileName, ".mp4") {
						mimeType = "video/mp4"
					} else if strings.HasSuffix(stepExport.FileName, ".ogg") || strings.HasSuffix(stepExport.FileName, ".mp3") {
						mimeType = "audio/mpeg"
					} else if strings.HasSuffix(stepExport.FileName, ".png") {
						mimeType = "image/png"
					} else if strings.HasSuffix(stepExport.FileName, ".jpg") || strings.HasSuffix(stepExport.FileName, ".jpeg") {
						mimeType = "image/jpeg"
					}

					localPath, err := a.saveMediaLocally(mediaData, mimeType, stepExport.FileName)
					if err == nil {
						content = localPath
					}
				}
			}

			step := models.OmniSequenceStep{
				BaseModel:      models.BaseModel{ID: uuid.New()},
				OrganizationID: orgID,
				SequenceID:     sequence.ID,
				StepOrder:      i + 1,
				Title:          stepExport.Title,
				Content:        content,
				MediaType:      stepExport.MediaType,
				DelayMs:        stepExport.DelayMs,
			}

			if stepExport.FileName != "" {
				step.FileName = &stepExport.FileName
			}

			if err := tx.Create(&step).Error; err != nil {
				tx.Rollback()
				return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create sequence step", nil, "")
			}
			importedSteps++
		}
	}

	tx.Commit()

	return r.SendEnvelope(map[string]interface{}{
		"status":              "success",
		"imported_categories": importedCategories,
		"imported_scripts":    importedScripts,
		"imported_sequences":  importedSequences,
		"imported_steps":      importedSteps,
	})
}
