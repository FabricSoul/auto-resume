package types

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type Project struct {
	Name       string    `toml:"name"`
	Path       string    `toml:"path"`
	CreatedAt  time.Time `toml:"created_at"`
	LastOpened time.Time `toml:"last_opened"`
}

type ProjectManager struct {
	baseDir    string
	configPath string
	Projects   []Project
}

// Add a new type to match the config structure
type Config struct {
	UserConfigPath string    `toml:"user_config_path"`
	Projects       []Project `toml:"projects"`
	Models         []AIModel `toml:"models"`
}

type AIModel struct {
	Name     string `toml:"name"`
	Provider string `toml:"provider"`
	Model    string `toml:"model"`
	APIKey   string `toml:"api_key"`
}

func NewPrejectManager() (*ProjectManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	baseDir := filepath.Join(homeDir, ".local", "share", "auto-resume")
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
		// Create default config with empty arrays
		defaultConfig := Config{
			UserConfigPath: "",
			Projects:       []Project{},
			Models:         []AIModel{},
		}

		data, err := toml.Marshal(defaultConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal default config: %w", err)
		}

		// Write default config to file
		if err := os.WriteFile(pm.configPath, data, 0644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check config file: %w", err)
	}

	// Load existing config
	data, err := os.ReadFile(pm.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Update ProjectManager with loaded projects
	pm.Projects = config.Projects

	return nil
}

func (pm *ProjectManager) SaveConfig() error {
	config := Config{
		UserConfigPath: "", // Default empty string as shown in spec
		Projects:       pm.Projects,
		Models:         []AIModel{}, // Empty for now, will be implemented later
	}

	data, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(pm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
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

	// Create projects directory if it doesn't exist
	projectsDir := filepath.Join(pm.baseDir, "projects")
	if err := os.MkdirAll(projectsDir, 0755); err != nil {
		return fmt.Errorf("failed to create projects directory: %w", err)
	}

	// Create project-specific directory
	projectDir := filepath.Join(projectsDir, name)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	project := Project{
		Name:       name,
		Path:       projectDir,
		CreatedAt:  time.Now(),
		LastOpened: time.Now(),
	}

	// Create project-specific config file
	projectConfig := filepath.Join(projectDir, "project.toml")
	defaultProjectConfig := fmt.Sprintf(`name = "%s"
model = ""
resume_input = ""`, name)

	if err := os.WriteFile(projectConfig, []byte(defaultProjectConfig), 0644); err != nil {
		return fmt.Errorf("failed to create project config: %w", err)
	}

	pm.Projects = append(pm.Projects, project)

	// Save the updated config
	if err := pm.SaveConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func (pm *ProjectManager) LoadProjects() error {
	// TODO: Implement loading projects from disk
	return nil
}

func (pm *ProjectManager) GetModels() []AIModel {
	data, err := os.ReadFile(pm.configPath)
	if err != nil {
		return []AIModel{}
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return []AIModel{}
	}

	return config.Models
}

func (pm *ProjectManager) SaveModels(models []AIModel) error {
	data, err := os.ReadFile(pm.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Update models in the config
	config.Models = models

	newData, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.WriteFile(pm.configPath, newData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
