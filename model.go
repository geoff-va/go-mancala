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
	DoneMoving
	SwitchPlayer
	Stealing
	GameOver
	CollectRemainder
)

func (s State) String() string {
	switch s {
	case SelectingPit:
		return "SelectingPit"
	case MovingFromHandToPit:
		return "MovingFromHandToPit"
	case DoneMoving:
		return "DoneMoving"
	case SwitchPlayer:
		return "SwitchPlayer"
	case Stealing:
		return "Stealing"
	case GameOver:
		return "GameOver"
	case CollectRemainder:
		return "CollectRemainder"
	}
	return ""
}

type Model struct {
	board         Board
	selectedPit   uint8
	lastPlacedPit uint8
	currentPlayer Player
	inHand        uint8
	state         State
	winner        Player
}

func NewModel() Model {
	return Model{
		board:         NewBoard(),
		currentPlayer: P1,
		selectedPit:   1,
	}
}
func NewModelWithState(boardState [14]uint8) Model {
	return Model{
		board:         NewBoardWithState(boardState),
		currentPlayer: P1,
		selectedPit:   1,
	}
}

func (s Model) Init() tea.Cmd {
	return nil
}
