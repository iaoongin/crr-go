package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	templateFile = "template.yaml"
	backupDir    = "./backups"
)

type TemplateManager struct{}

func NewTemplateManager() *TemplateManager {
	return &TemplateManager{}
}

func (tm *TemplateManager) GetTemplate() (string, error) {
	content, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %v", err)
	}
	return string(content), nil
}

func (tm *TemplateManager) SaveTemplate(content []byte) error {
	// 创建备份
	backupFile := fmt.Sprintf("%s/template_%s.yaml", backupDir, time.Now().Format("20060102_150405"))
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}
	if err := ioutil.WriteFile(backupFile, content, 0644); err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	// 保存新模板
	if err := ioutil.WriteFile(templateFile, content, 0644); err != nil {
		return fmt.Errorf("failed to save template: %v", err)
	}

	return nil
}

func (tm *TemplateManager) RestoreTemplate() error {
	backups, err := tm.GetBackups()
	if err != nil {
		return fmt.Errorf("failed to get backups: %v", err)
	}

	if len(backups) == 0 {
		return fmt.Errorf("no backups available")
	}

	latestBackup := backups[len(backups)-1]
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", backupDir, latestBackup))
	if err != nil {
		return fmt.Errorf("failed to read backup: %v", err)
	}

	if err := ioutil.WriteFile(templateFile, content, 0644); err != nil {
		return fmt.Errorf("failed to restore template: %v", err)
	}

	return nil
}

func (tm *TemplateManager) GetBackups() ([]string, error) {
	files, err := ioutil.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %v", err)
	}

	var backups []string
	for _, file := range files {
		if !file.IsDir() {
			backups = append(backups, file.Name())
		}
	}

	return backups, nil
}

func (tm *TemplateManager) RollbackToBackup(backupFile string) error {
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", backupDir, backupFile))
	if err != nil {
		return fmt.Errorf("failed to read backup: %v", err)
	}

	if err := ioutil.WriteFile(templateFile, content, 0644); err != nil {
		return fmt.Errorf("failed to rollback template: %v", err)
	}

	return nil
}

func SetupTemplateRoutes(r *gin.Engine, tm *TemplateManager) {
	r.GET("/api/template", func(c *gin.Context) {
		content, err := tm.GetTemplate()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.String(200, content)
	})

	r.POST("/api/template", func(c *gin.Context) {
		content, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid data"})
			return
		}
		if err := tm.SaveTemplate(content); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Template saved successfully"})
	})

	r.POST("/api/template/restore", func(c *gin.Context) {
		if err := tm.RestoreTemplate(); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Template restored successfully"})
	})

	r.GET("/api/backups", func(c *gin.Context) {
		backups, err := tm.GetBackups()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"backups": backups})
	})

	r.POST("/api/backups/rollback", func(c *gin.Context) {
		backupFile := c.PostForm("file")
		if backupFile == "" {
			c.JSON(400, gin.H{"error": "Backup file not specified"})
			return
		}
		if err := tm.RollbackToBackup(backupFile); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Template rolled back successfully"})
	})
}