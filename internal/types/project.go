package types

import (
	// "github.com/pelletier/go-toml/v2"
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
	baseDir    string
	Projects   []Project
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

	// Initialize with empty projects slice
	return &ProjectManager{
		baseDir:  baseDir,
		Projects: make([]Project, 0),
	}, nil
}

func (pm *ProjectManager) AddProject(name string) error {
	project := Project{
		Name:      name,
		CreatedAt: time.Now(),
		Path:      filepath.Join(pm.baseDir, name),
	}
	
	pm.Projects = append(pm.Projects, project)
	return nil
}

func (pm *ProjectManager) LoadProjects() error {
	// TODO: Implement loading projects from disk
	return nil
}
