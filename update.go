package main

import (
	"log"
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
	var cmd tea.Cmd = nil

	log.Printf("State: %d", s.state)
	switch s.state {
	case SelectingPit:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q":
				return s, tea.Quit
			case "h":
				s.MoveLeft()
			case "l":
				s.MoveRight()
			case "enter", " ":
				s.SelectPit()
				cmd = doTick(1)
			}
		}
	case MovingFromHandToPit:
		isDone := s.moveFromHandToPit()
		if !isDone {
			cmd = doTick(1000)
			log.Printf("Not done; more in hand")
		} else {
			cmd = doTick(1)
		}
	case IsWinner:
		// TODO: Evaluate winner
		s.state = SwitchPlayer
		cmd = doTick(1)
		log.Println("Switching player")
	case SwitchPlayer:
		s.SwitchPlayer()
		s.state = SelectingPit
		cmd = doTick(1)
	case GameOver:
	}

	return s, cmd
}

func (s *Model) MoveRight() {
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
}

func (s *Model) MoveLeft() {
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
}

// getPlayerBounds returns the lower and upper movable bounds of the pits for a given player.
func getPlayerBounds(p Player) (lBound, uBound uint8) {
	lBound = uint8(p)*7 + 1
	uBound = lBound + 5
	return
}

func (s *Model) SelectPit() {
	numInPit := s.board[s.selectedPit]

	if numInPit == 0 {
		panic("can't select an empty pit")
	}

	s.board[s.selectedPit] = 0
	s.inHand = numInPit
	s.selectedNum = numInPit
	s.state = MovingFromHandToPit
	s.lastSelectedPit[s.currentPlayer] = s.selectedPit
}

func (s *Model) moveFromHandToPit() bool {
	otherPlayer := Player((s.currentPlayer + 1) % 2)
	otherStore := s.getStoreIndex(otherPlayer)
	if s.inHand > 0 {
		pitIndex := (s.selectedPit + 1 + s.selectedNum - s.inHand)
		if pitIndex == otherStore {
			pitIndex++
		}
		pitIndex %= 14
		s.board[pitIndex]++
		s.inHand--
		return false
	}

	s.state = IsWinner
	return true

}

func (s *Model) SwitchPlayer() {
	if s.currentPlayer == P1 {
		s.currentPlayer = P2
		s.selectedPit = s.lastSelectedPit[P2]
	} else {
		s.currentPlayer = P1
		s.selectedPit = s.lastSelectedPit[P1]
	}
}

func (s Model) getStoreIndex(player Player) uint8 {
	return uint8(player)*7 + 6
}
