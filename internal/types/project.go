package types

import (
	"github.com/pelletier/go-toml/v2"
	"os"
	"path/filepath"
	"time"
)

type Project struct {
	Name       string    `toml:"name"`
	Path       string    `toml: "path"`
	CreatedAt  time.Time `toml: "created_at"`
	LastOpened time.Time `toml: "last_opened"`
}

type ProjectManager struct {
	baseDir string
}

func NewPrejectManager() (*ProjectManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	baseDir := filepath.Join(homeDir, ".local", "share", "autoResume", "projects")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}

	return &ProjectManager{baseDir: baseDir}, nil
}
