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
	board         Board
	selectedPit   uint8
	lastPlacedPit uint8
	currentPlayer Player
	inHand        uint8
	state         State
	winner        Player
	isWinner      bool
}

func NewState() Model {
	return Model{
		board:         NewBoard(),
		currentPlayer: P1,
		selectedPit:   1,
	}
}

func (s Model) Init() tea.Cmd {
	return nil
}
