package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/lipgloss"
)

func (s Model) View() string {
	board := s.renderPitRow(P1)
	board += fmt.Sprintf("%d                      %d\n", s.board.board[0], s.board.board[7])
	board += s.renderPitRow(P2)
	return board
}

func (s Model) renderPitRow(p Player) string {
	pits := ""
	lBound, uBound := getPlayerBounds(p)
	log.Printf("Current Player: %d", s.currentPlayer)
	log.Printf("View selectedPit: %d, lBound: %d, uBound: %d", s.selectedPit, lBound, uBound)

	if p == P1 {
		for i := uint8(lBound); i <= uBound; i++ {
			if s.selectedPit == i {
				style := lipgloss.NewStyle().Background(lipgloss.Color("0")).Foreground(lipgloss.Color("12"))
				pits += style.Render(fmt.Sprintf("%d", s.board.board[i]))
				pits += " "
			} else {
				pits += fmt.Sprintf("%d ", s.board.board[i])
			}
		}
	} else {
		for i := uint8(uBound); i >= lBound; i-- {
			if s.selectedPit == i {
				style := lipgloss.NewStyle().Background(lipgloss.Color("0")).Foreground(lipgloss.Color("12"))
				pits += style.Render(fmt.Sprintf("%d", s.board.board[i]))
				pits += " "
			} else {
				pits += fmt.Sprintf("%d ", s.board.board[i])
			}
		}
	}
	return "    " + pits + "\n"
}
