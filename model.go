package main

import tea "github.com/charmbracelet/bubbletea"

type Board struct {
	cells [6]uint8
	store uint8
}

func (b Board) Init() tea.Cmd {
	return nil
}

func NewBoard() Board {
	return Board{
		cells: [6]uint8{4, 4, 4, 4, 4, 4},
		store: 0,
	}
}

type State struct {
	p1 Board
	p2 Board
}

func (s State) Init() tea.Cmd {
	return nil
}
