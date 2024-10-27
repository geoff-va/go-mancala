// p1_0 p1_1 p1_2 p1_3 p1_4 p1_5 p1_s
// p2_0 p2_1 p2_2 p2_3 p2_4 p2_5 p2_s
package main

import tea "github.com/charmbracelet/bubbletea"

type Player uint8

// Players
const (
	P1 Player = iota
	P2
)

type State uint8

// Game States
const (
	SelectingPit State = iota
	MovingFromHandToPit
	IsWinner
	SwitchPlayer
	Stealing
	GameOver
)

type Model struct {
	board           [14]uint8
	currentPlayer   Player
	selectedPit     uint8
	lastPlacedPit   uint8
	inHand          uint8
	lastSelectedPit map[Player]uint8
	state           State
}

func NewState() Model {
	return Model{
		board:           [14]uint8{0, 4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 4},
		currentPlayer:   P1,
		lastSelectedPit: map[Player]uint8{P1: 1, P2: 13},
		state:           SelectingPit,
		selectedPit:     1,
	}
}

func (s Model) Init() tea.Cmd {
	return nil
}

func (s Model) numInStore(p Player) uint8 {
	return s.board[p*7+6]
}
