// p1_0 p1_1 p1_2 p1_3 p1_4 p1_5 p1_s
// p2_0 p2_1 p2_2 p2_3 p2_4 p2_5 p2_s
package main

import tea "github.com/charmbracelet/bubbletea"

type Player uint8

const (
	// Players
	P1 Player = iota
	P2

	// Game States
	SelectingPit uint8 = iota
	MovingFromHandToPit
	IsWinner
	SwitchPlayer
	GameOver
)

type State struct {
	board            [14]uint8
	currentPlayer    Player
	selectedPit      uint8
	inHand           uint8
	lastSelectedCell map[Player]uint8
	state            uint8
}

func NewState() State {
	return State{
		board:            [14]uint8{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 4, 0},
		currentPlayer:    P1,
		lastSelectedCell: map[Player]uint8{P1: 0, P2: 7},
		state:            SelectingPit,
	}
}

func (s State) Init() tea.Cmd {
	return nil
}

func (s State) numInStore(p Player) uint8 {
	return s.board[p*7+6]
}
