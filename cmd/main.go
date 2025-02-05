// cmd/main.go
package main

import (
	"fmt"
	"os"

	"github.com/FabricSoul/auto-resume/internal/models"
	"github.com/FabricSoul/auto-resume/internal/types"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	pm, err := types.NewPrejectManager()
	if err != nil {
		fmt.Errorf("Failed to init pm: %v \n", err)
		os.Exit(1)
	}
	p := tea.NewProgram(models.NewMainModel(pm), tea.WithMouseCellMotion(),
		tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
}
