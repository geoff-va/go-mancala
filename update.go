package main

import (
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
		if nextState == MovingFromHandToPit || nextState == Stealing {
			cmd = doTick(500)
		} else {
			cmd = doTick(1)
		}
	case Stealing:
		nextState = s.HandleSteal()
		cmd = doTick(1)
	case IsWinner:
		nextState = s.HandleIsWinner()
		cmd = doTick(1)
	case SwitchPlayer:
		nextState = s.HandleSwitchPlayer()
		cmd = doTick(1)
	case GameOver:
		return s, tea.Quit
	}

	s.state = nextState
	return s, cmd
}

func (s *Model) HandleMoveRight() State {
	// TODO: handle moving past empty pits
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
	// BUG: handle game ending in own store gets another turn but no stones left
	s.inHand, s.lastPlacedPit = s.board.MoveFromHandToPit(s.inHand, s.lastPlacedPit, s.currentPlayer)

	if s.inHand > 0 {
		return MovingFromHandToPit
	}

	// Rule: You get another turn if you end in your store
	if s.board.IsPlayersStore(s.lastPlacedPit, s.currentPlayer) {
		return SelectingPit
	}

	// Turn over
	if s.board.GetPlayerForPit(s.lastPlacedPit) != s.currentPlayer ||
		s.board.IsPlayersStore(s.lastPlacedPit, s.currentPlayer) ||
		s.board.Get(s.lastPlacedPit) != 1 ||
		s.board.Get(s.board.GetOppositePit(s.lastPlacedPit)) == 0 {
		return IsWinner
	}

	return Stealing

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
	return IsWinner
}

func (s Model) HandleIsWinner() State {
	p1wins := true
	p2wins := true

	for i := range 5 {
		p1wins = p1wins && s.board.Get(uint8(i+1)) == 0
		p2wins = p2wins && s.board.Get(uint8(i+7)) == 0
	}

	if p1wins {
		s.winner = P1
		s.isWinner = true

	}

	if p2wins {
		s.winner = P2
		s.isWinner = true
	}

	if s.isWinner {
		return GameOver
	}

	return SwitchPlayer

}
