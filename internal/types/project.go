package types

import (
	// "github.com/pelletier/go-toml/v2"
	"fmt"
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
	configPath string
	Projects   []Project
}

func NewPrejectManager() (*ProjectManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	baseDir := filepath.Join(homeDir, ".local", "share", "autoResume")
	configPath := filepath.Join(baseDir, "config.toml")

	// Create base directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}

	pm := &ProjectManager{
		baseDir:    baseDir,
		configPath: configPath,
		Projects:   make([]Project, 0),
	}

	// Check if config exists, if not create it
	if err := pm.initializeConfig(); err != nil {
		return nil, err
	}

	return pm, nil
}

func (pm *ProjectManager) initializeConfig() error {
	// Check if config file exists
	_, err := os.Stat(pm.configPath)
	if os.IsNotExist(err) {
		// Create default config content
		defaultConfig := `user_config_path = ""`

		// Write default config to file
		err = os.WriteFile(pm.configPath, []byte(defaultConfig), 0644)
		if err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check config file: %w", err)
	}

	return nil
}

func (pm *ProjectManager) AddProject(name string) error {
	// Check if project with same name exists
	for _, p := range pm.Projects {
		if p.Name == name {
			return fmt.Errorf("project with name '%s' already exists", name)
		}
	}

	project := Project{
		Name:       name,
		CreatedAt:  time.Now(),
		LastOpened: time.Now(),
		Path:      filepath.Join(pm.baseDir, name),
	}
	
	// Create project directory
	if err := os.MkdirAll(project.Path, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	pm.Projects = append(pm.Projects, project)
	return nil
}

func (pm *ProjectManager) LoadProjects() error {
	// TODO: Implement loading projects from disk
	return nil
}
