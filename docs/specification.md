## Table of Contents

- [1. Project Overview](#1-project-overview)
- [2. Project Structure](#2-project-structure)
- [3. Configuration Management](#3-configuration-management)
- [4. State Management](#4-state-management)
- [5. Models and Components](#5-models-and-components)
- [6. User Interface](#6-user-interface)
- [7. Data Persistence](#7-data-persistence)
- [8. Error Handling](#8-error-handling)
- [9. Testing Strategy](#9-testing-strategy)
- [10. Development Guidelines](#10-development-guidelines)

## 1. Project Overview

### 1.1 Purpose

A TUI app uses LLM to automate the process of generating resumes in LaTex
format, and then compiles the LaTex file to PDF format, targeting different
jobs.

### 1.2 Features

- Generate targeted resume with support of LLM.
- Able to use LLM from both local LLM via Ollama or use API with gollm.
- Can have multiple projects to store multiple resumes.

### 1.3 Target Users

Anyone who wants to automate their targeted resume process can use this app to
generate the resume.

## 2. Project Structure

```
auto-resume/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── defaults.go
│   ├── models/
│   │   ├── main_model.go
│   │   ├── overview.go
│   │   └── [feature].go
│   ├── state/
│   │   ├── store.go
│   │   └── actions.go
│   └── utils/
│       └── helpers.go
├── pkg/
│   └── [reusable packages]
├── docs/
│   ├── specifications.md
├── tests/
│   ├── unit/
│   └── integration/
├── .gitignore
├── go.mod
├── go.sum
├── README.md
└── Makefile
```

## 3. Configuration Management

### 3.1 Configuration File Locations

- App-specific: `$HOME/.local/share/auto-resume/config.toml`
- User-specific: `$HOME/.config/auto-resume/config.toml`
- Project-specific: `./project.toml`

### 3.2 Configuration Structure

App-specific

```toml
user_config_path = "$HOME/.config/auto-resume"

[[projects]]
name = "First Project" # need to be unique
created_at = "2025-02-05T08:36:26-05:00"
last_opened = "2025-02-05T08:37:50-05:00"
path = "path/to/project.toml"

[[projects]]
name = "Second Project"
created_at = "2025-02-05T08:36:26-05:00"
last_opened = "2025-02-05T08:37:50-05:00"
path = "path/to/project.toml"

[[models]]
name = "Model 0" # need to be unique
provider = "openai"
model = "gpt-4o-mini"
api_key = "your api key here"

[[models]]
name = "Model 1"
provider = "ollama"
model = "deepseek-r1:8b"
api_key = "" # Ollama models does not need api key
```

User-specific:

```toml
TBD
```

Project-specific:

```toml
name = "Second Project"
model = "Model 1"
resume_input = """
the latex content of the original resume
"""

[[outputs]]
name = "the name of this targeted resume"
job_description = """
the job description of this job
"""
output = """
the new generated resume
"""
```

### 3.3 Configuration Loading Priority

1. Command-line flags
2. Environment variables
3. Project-specific config
4. User-specific config
5. Default config

## 4. State Management

### 4.1 Global State Structure

```go
type Project struct {
	Name       string    `toml:"name"`
	Path       string    `toml: "path"`
	CreatedAt  time.Time `toml: "created_at"`
	LastOpened time.Time `toml: "last_opened"`
}
```

### 4.2 State Updates

- Use message passing for state updates
- Implement command pattern for side effects
- Use tea.Cmd for asynchronous operations

### 4.3 State Persistence

- Save state on exit
- Load state on startup
- Auto-save intervals for critical data

## 5. Models and Components

### 5.1 Core Models

1. Main Model
   - Manages overall application state
   - Handles global keyboard shortcuts
   - Routes messages to child models

2. Overview Model
   - Displays and manages the projects read form App-specific `config.toml`
   - Creates new project

3. Project Model
   - Displays specific project info

4. Error Model
   - Displays errors in a popup window

5. LLM Manager Model
   - List and manages the LLM configs from App-specific `config.toml`

### 5.2 Component Hierarchy

```
Main Model
├── Overview Model
├── Project Model
├── Error Model
└── LLM Manager Model
```

### 5.3 Message Flow

[Diagram or description of how messages flow between models]

## 6. User Interface

### 6.1 Layout Structure

[Describe your layout structure and components]

### 6.2 Theme Management

- Color schemes
- Style definitions
- Component-specific styles

### 6.3 Keyboard Shortcuts

```
Global:
- Ctrl+C: Exit

Navigation:
- Tab: Next item
- Shift+Tab: Previous item
- Enter: Select
- Up: up, k
- Down: down, j

Feature-specific:
[List feature-specific shortcuts]
```

## 7. Data Persistence

### 7.1 Storage Locations

- Configuration: `$HOME/.config/auto-resume/`
- Data: `$HOME/.local/share/auto-resume/`
- Cache: `$HOME/.cache/auto-resume/`
- Logs: `$HOME/.local/state/auto-resume/logs/`

### 7.2 Data Formats

- Configuration: TOML
- User data: TOML
- Cache: Binary
- Logs: Plain text

### 7.3 Backup Strategy

[Describe backup and recovery procedures]

## 8. Error Handling

### 8.1 Error Types

- UI errors
- Configuration errors
- Data errors
- Network errors

### 8.2 Error Reporting

- Log levels
- User notifications
- Error recovery procedures

## 9. Testing Strategy

### 9.1 Unit Tests

- Model testing
- Component testing
- Utility function testing

### 9.2 Integration Tests

- Model interaction testing
- State management testing
- Configuration testing

### 9.3 UI Tests

- Component rendering tests
- User interaction tests
- Layout tests

## 10. Development Guidelines

### 10.1 Code Style

- Follow Go best practices
- Use consistent naming conventions
- Document public APIs

### 10.2 Git Workflow

- Branch naming convention
- Commit message format
- PR review process

### 10.3 Documentation

- Keep documentation in sync with code
- Update changelog
- Maintain examples

```
## Supporting Documentation

### Additional Specification Files

1. `docs/architecture.md`:
   - Detailed system architecture
   - Component interactions
   - Design decisions

2. `docs/models.md`:
   - Model specifications
   - State management details
   - Message handling

3. `docs/state.md`:
   - State structure
   - Update patterns
   - Persistence strategy

4. `docs/components.md`:
   - UI component library
   - Styling guidelines
   - Usage examples
```
