package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var version = "0.0.1"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println(version)
		os.Exit(0)
	}

	p := tea.NewProgram(NewBoard())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting program: %v", err)
		os.Exit(1)
	}

}
