package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/shridarpatil/whatomate/internal/models"
	"github.com/valyala/fasthttp"
	"github.com/zerodha/fastglue"
)

type ConfigData struct {
	Name string            `json:"name"`
	Rows []json.RawMessage `json:"rows"`
}

type ConfigFile struct {
	Data []ConfigData `json:"data"`
}

type GroupRow struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
	Items []struct {
		ID int `json:"id"`
	} `json:"items"`
}

type ScriptRow struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Items []struct {
		ID    int `json:"id"`
		Delay int `json:"delay"`
	} `json:"items"`
}

type ItemRow struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	FileName string `json:"filename"`
}

// ImportOmnisZIP accepts a multipart form with a .zip file, parses config.json and saves media
func (a *App) ImportOmnisZIP(r *fastglue.Request) error {
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
	fileMap := make(map[string]*zip.File)

	for _, zf := range zipReader.File {
		if strings.Contains(zf.Name, "__MACOSX") || strings.Contains(zf.Name, ".DS_Store") {
			continue
		}
		
		baseName := filepath.Base(zf.Name)
		if baseName == "config.json" && configJsonFile == nil {
			configJsonFile = zf
		}
		fileMap[baseName] = zf
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

	var parsed ConfigFile
	if err := json.Unmarshal(configBytes, &parsed); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid config.json format", nil, "")
	}

	groups := []GroupRow{}
	scripts := map[int]ScriptRow{}
	items := map[int]ItemRow{}

	for _, d := range parsed.Data {
		if d.Name == "groups" {
			for _, row := range d.Rows {
				var grp GroupRow
				_ = json.Unmarshal(row, &grp)
				groups = append(groups, grp)
			}
		} else if d.Name == "scripts" {
			for _, row := range d.Rows {
				var scr ScriptRow
				_ = json.Unmarshal(row, &scr)
				scripts[scr.ID] = scr
			}
		} else if d.Name == "items" {
			for _, row := range d.Rows {
				var it ItemRow
				_ = json.Unmarshal(row, &it)
				items[it.ID] = it
			}
		}
	}

	tx := a.DB.Begin()
	tx.Unscoped().Where("organization_id = ?", orgID).Delete(&models.OmniScript{})
	tx.Unscoped().Where("organization_id = ?", orgID).Delete(&models.OmniCategory{})

	for _, grp := range groups {
		catID := uuid.New()
		cat := models.OmniCategory{
			BaseModel:      models.BaseModel{ID: catID},
			OrganizationID: orgID,
			Name:           grp.Title,
			Color:          grp.Color,
		}
		if err := tx.Create(&cat).Error; err != nil {
			tx.Rollback()
			return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create category", nil, "")
		}

		for _, grpItem := range grp.Items {
			scr, ok := scripts[grpItem.ID]
			if !ok {
				continue
			}

			for _, scrItem := range scr.Items {
				it, ok := items[scrItem.ID]
				if !ok {
					continue
				}

				omniScr := models.OmniScript{
					BaseModel:      models.BaseModel{ID: uuid.New()},
					OrganizationID: orgID,
					CategoryID:     catID,
					Title:          scr.Title,
					Content:        it.Content,
					MediaType:      it.Type,
					DelayMs:        scrItem.Delay,
				}

				if it.FileName != "" {
					fName := it.FileName
					omniScr.FileName = &fName
					
					if zFile, exists := fileMap[fName]; exists {
						mf, _ := zFile.Open()
						mediaBytes, _ := io.ReadAll(mf)
						mf.Close()

						mimeType := "application/octet-stream"
						if strings.HasSuffix(fName, ".mp4") { mimeType = "video/mp4" }
						if strings.HasSuffix(fName, ".ogg") { mimeType = "audio/ogg" }
						if strings.HasSuffix(fName, ".png") { mimeType = "image/png" }

						localPath, err := a.saveMediaLocally(mediaBytes, mimeType, fName)
						if err == nil {
							omniScr.Content = localPath
						}
					}
				}

				if err := tx.Create(&omniScr).Error; err != nil {
					tx.Rollback()
					return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create script", nil, "")
				}
			}
		}
	}

	tx.Commit()
	return r.SendEnvelope(map[string]string{"status": "success"})
}
