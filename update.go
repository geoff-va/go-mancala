package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var oppositeIdxOffsets [14]uint8 = [14]uint8{7, 13, 12, 11, 10, 9, 8, 0, 6, 5, 4, 3, 2, 1}

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
		if nextState == MovingFromHandToPit {
			cmd = doTick(1000)
		} else {
			cmd = doTick(1)
		}
	case Stealing:
		nextState = s.HandleSteal()
		cmd = doTick(1)

	case IsWinner:
		// TODO: Evaluate winner
		nextState = SwitchPlayer
		cmd = doTick(1)
	case SwitchPlayer:
		nextState = s.HandleSwitchPlayer()
		cmd = doTick(1)
	case GameOver:
	}

	s.state = nextState
	return s, cmd
}

func (s *Model) HandleMoveRight() State {
	// TODO: handle moving past empty pits
	lBound, uBound := getPlayerBounds(s.currentPlayer)
	if s.currentPlayer == P2 {
		if s.selectedPit > lBound {
			s.selectedPit--
		} else {
			s.selectedPit = uBound
		}
	} else {
		if s.selectedPit < uBound {
			s.selectedPit++
		} else {
			s.selectedPit = lBound
		}
	}
	return SelectingPit
}

func (s *Model) HandleMoveLeft() State {
	// TODO: handle moving past empty pits
	lBound, uBound := getPlayerBounds(s.currentPlayer)
	if s.currentPlayer == P2 {
		if s.selectedPit < uBound {
			s.selectedPit++
		} else {
			s.selectedPit = lBound
		}
	} else {
		if s.selectedPit > lBound {
			s.selectedPit--
		} else {
			s.selectedPit = uBound
		}
	}
	return SelectingPit
}

// getPlayerBounds returns the lower and upper movable bounds of the pits for a given player.
func getPlayerBounds(p Player) (lBound, uBound uint8) {
	lBound = uint8(p)*7 + 1
	uBound = lBound + 5
	return
}

func (s *Model) HandleSelectPit() State {
	numInPit := s.board[s.selectedPit]

	if numInPit == 0 {
		panic("can't select an empty pit")
	}

	s.board[s.selectedPit] = 0
	s.inHand = numInPit
	s.lastSelectedPit[s.currentPlayer] = s.selectedPit
	s.lastPlacedPit = s.selectedPit
	return MovingFromHandToPit
}

func (s *Model) HandleMoveFromHandToPit() State {
	// If something in your hand, place it in next pit
	// After if you have nothing left, check if you can steal
	otherPlayer := Player((s.currentPlayer + 1) % 2)
	otherStore := s.getStoreIndex(otherPlayer)

	// get next pit we're going to place a stone in
	var pitIndex uint8
	if s.lastPlacedPit == 0 {
		pitIndex = 13
	} else {
		pitIndex = s.lastPlacedPit - 1
	}

	if s.inHand > 0 {
		// Skip the other player's store
		if pitIndex == otherStore {
			pitIndex--
		}
		s.inHand--
		s.board[pitIndex]++
		s.lastPlacedPit = pitIndex
		return MovingFromHandToPit
	}

	pitIndex = s.lastPlacedPit
	// You get another turn if you end in your store
	if pitIndex == s.getStoreIndex(s.currentPlayer) {
		return SelectingPit
	}

	// Turn over
	if !onPlayersSide(pitIndex, s.currentPlayer) ||
		pitIndex == s.getStoreIndex(s.currentPlayer) ||
		s.board[pitIndex] != 1 ||
		s.board[GetOppositePit(pitIndex)] == 0 {
		return IsWinner
	}

	return Stealing

}

func (s *Model) HandleSwitchPlayer() State {
	if s.currentPlayer == P1 {
		s.currentPlayer = P2
		s.selectedPit = s.lastSelectedPit[P2]
	} else {
		s.currentPlayer = P1
		s.selectedPit = s.lastSelectedPit[P1]
	}
	return SelectingPit
}

func (s Model) getStoreIndex(player Player) uint8 {
	return (uint8(player) * 7)
}

func (s *Model) HandleSteal() State {
	oppositePit := GetOppositePit(s.lastPlacedPit)
	s.board[s.getStoreIndex(s.currentPlayer)] += s.board[oppositePit] + 1
	s.board[oppositePit] = 0
	s.board[s.lastPlacedPit] = 0
	return IsWinner
}

func GetOppositePit(pit uint8) uint8 {
	return oppositeIdxOffsets[pit]
}

func onPlayersSide(pit uint8, player Player) bool {
	lBound, uBound := getPlayerBounds(player)
	return pit >= lBound && pit <= uBound
}
