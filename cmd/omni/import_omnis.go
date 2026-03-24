package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/google/uuid"
	"github.com/omni-platform/omni/internal/config"
	"github.com/omni-platform/omni/internal/database"
	"github.com/omni-platform/omni/internal/models"
	"github.com/zerodha/logf"
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

func runImportOmnis(args []string) {
	flags := flag.NewFlagSet("import-omnis", flag.ExitOnError)
	configPath := flags.String("config", "config.toml", "Path to config file")
	jsonPath := flags.String("json", "config.json", "Path to config.json")
	_ = flags.Parse(args)

	lo := logf.New(logf.Opts{
		EnableColor:     true,
		Level:           logf.DebugLevel,
		EnableCaller:    true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	lo.Info("Starting Omni Import...", "jsonFile", *jsonPath)

	cfg, err := config.Load(*configPath)
	if err != nil {
		lo.Fatal("Failed to load config", "error", err)
	}

	db, err := database.NewPostgres(&cfg.Database, false)
	if err != nil {
		lo.Fatal("Failed to connect to DB", "error", err)
	}

	var org models.Organization
	if err := db.First(&org).Error; err != nil {
		lo.Fatal("No organization found to attach omnis to", "error", err)
	}

	fileBytes, err := os.ReadFile(*jsonPath)
	if err != nil {
		lo.Fatal("Failed to read JSON file", "error", err)
	}

	var parsed ConfigFile
	if err := json.Unmarshal(fileBytes, &parsed); err != nil {
		lo.Fatal("Failed to parse JSON file", "error", err)
	}

	groups := []GroupRow{}
	scripts := map[int]ScriptRow{}
	items := map[int]ItemRow{}

	for _, d := range parsed.Data {
		if d.Name == "groups" {
			for _, r := range d.Rows {
				var grp GroupRow
				_ = json.Unmarshal(r, &grp)
				groups = append(groups, grp)
			}
		} else if d.Name == "scripts" {
			for _, r := range d.Rows {
				var scr ScriptRow
				_ = json.Unmarshal(r, &scr)
				scripts[scr.ID] = scr
			}
		} else if d.Name == "items" {
			for _, r := range d.Rows {
				var it ItemRow
				_ = json.Unmarshal(r, &it)
				items[it.ID] = it
			}
		}
	}

	lo.Info("Parsed JSON", "groups", len(groups), "scripts", len(scripts), "items", len(items))

	tx := db.Begin()

	// Hard delete previous omnis since we're reimporting
	tx.Unscoped().Where("organization_id = ?", org.ID).Delete(&models.OmniScript{})
	tx.Unscoped().Where("organization_id = ?", org.ID).Delete(&models.OmniCategory{})

	for _, grp := range groups {
		catID := uuid.New()
		cat := models.OmniCategory{
			BaseModel:      models.BaseModel{ID: catID},
			OrganizationID: org.ID,
			Name:           grp.Title,
			Color:          grp.Color,
		}
		if err := tx.Create(&cat).Error; err != nil {
			tx.Rollback()
			lo.Fatal("Failed to create category", "err", err)
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
					OrganizationID: org.ID,
					CategoryID:     catID,
					Title:          scr.Title,
					Content:        it.Content,
					MediaType:      it.Type,
					DelayMs:        scrItem.Delay,
				}
				if it.FileName != "" {
                    fname := it.FileName
					omniScr.FileName = &fname
				}

				if err := tx.Create(&omniScr).Error; err != nil {
					tx.Rollback()
					lo.Fatal("Failed to create script", "err", err)
				}
			}
		}
	}

	tx.Commit()
	lo.Info("Omni import completed successfully.")
}
