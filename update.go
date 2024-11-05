package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg time.Time

func doTick(t int) tea.Cmd {
	return tea.Tick(time.Duration(t)*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (s Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Catch quit early
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return s, tea.Quit
		}
	}
	var cmd tea.Cmd = nil
	var nextState State

	switch s.state {
	case SelectingPit:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "h":
				nextState = s.HandleMoveLeft()
			case "l":
				nextState = s.HandleMoveRight()
			case "enter", " ":
				nextState = s.HandleSelectPit()
				cmd = doTick(1)
			}
		}
	case MovingFromHandToPit:
		nextState = s.HandleMoveFromHandToPit()
		cmd = doTick(500)
	case DoneMoving:
		nextState = s.HandleDoneMoving()
		if nextState == Stealing || nextState == CollectRemainder {
			cmd = doTick(500)
		} else {
			cmd = doTick(1)
		}
		cmd = doTick(1)
	case Stealing:
		nextState = s.HandleSteal()
		cmd = doTick(1)
	case SwitchPlayer:
		nextState = s.HandleSwitchPlayer()
		cmd = doTick(1)
	case CollectRemainder:
		nextState = s.HandleCollectRemainder()
		cmd = doTick(1)
	case GameOver:
		return s, tea.Quit
	}

	s.state = nextState
	return s, cmd
}

func (s *Model) HandleMoveRight() State {
	s.selectedPit = s.board.MoveRight(s.selectedPit, s.currentPlayer)
	return SelectingPit
}

func (m *Model) HandleMoveLeft() State {
	m.selectedPit = m.board.MoveLeft(m.selectedPit, m.currentPlayer)
	return SelectingPit
}

func (s *Model) HandleSelectPit() State {
	s.inHand = s.board.SelectPit(s.selectedPit)
	s.lastPlacedPit = s.selectedPit
	return MovingFromHandToPit
}

func (s *Model) HandleMoveFromHandToPit() State {
	s.inHand, s.lastPlacedPit = s.board.MoveFromHandToPit(s.inHand, s.lastPlacedPit, s.currentPlayer)

	if s.inHand > 0 {
		return MovingFromHandToPit
	}

	return DoneMoving
}

func (m *Model) HandleDoneMoving() State {
	if m.inHand > 0 {
		panic(fmt.Sprintf("inHand should be 0, but is %d", m.inHand))
	}

	gameOver := m.isGameOver()
	if gameOver {
		return CollectRemainder
	}

	// Rule: You get another turn if you end in your store
	if m.board.IsPlayersStore(m.lastPlacedPit, m.currentPlayer) {
		m.selectedPit = m.board.GetFirstNonEmptyPit(m.currentPlayer)
		return SelectingPit
	}

	// Rule: Steal if last stone landed in empty pit on your side and the the opposite
	// pit has stones in it
	if m.board.GetPlayerForPit(m.lastPlacedPit) == m.currentPlayer &&
		m.board.Get(m.lastPlacedPit) == 1 &&
		m.board.Get(m.board.GetOppositePit(m.lastPlacedPit)) > 0 {
		return Stealing
	}

	return SwitchPlayer
}

func (m *Model) HandleCollectRemainder() State {
	m.board.CollectRemainder()
	return GameOver
}

func (s *Model) HandleSwitchPlayer() State {
	if s.currentPlayer == P1 {
		s.currentPlayer = P2
		s.selectedPit = s.board.GetFirstNonEmptyPit(P2)
	} else {
		s.currentPlayer = P1
		s.selectedPit = s.board.GetFirstNonEmptyPit(P1)
	}
	return SelectingPit
}

func (s *Model) HandleSteal() State {
	s.board.Steal(s.currentPlayer, s.lastPlacedPit)
	return SwitchPlayer
}

func (s Model) isGameOver() bool {
	p1Empty := true
	p2Empty := true

	for i := range 6 {
		p1Empty = p1Empty && s.board.Get(uint8(i+1)) == 0
		p2Empty = p2Empty && s.board.Get(uint8(i+7)) == 0
	}

	return p1Empty || p2Empty
}
