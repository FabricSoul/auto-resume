package types

import "errors"

var (
	ErrEmptyProjectName = errors.New("project name cannot be empty")
	ErrEmptyModelName   = errors.New("model name cannot be empty")
)
