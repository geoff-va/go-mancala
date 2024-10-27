//	1  2  3  4  5  6
//
// 0                    7
//
//	13 12 11 10 9  8
package main

import (
	"log"
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
	case Stealing:
		s.steal()
		cmd = doTick(1000)

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
	s.state = MovingFromHandToPit
	s.lastSelectedPit[s.currentPlayer] = s.selectedPit
	s.lastPlacedPit = s.selectedPit
}

func (s *Model) moveFromHandToPit() bool {
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
	}

	// Stones Left
	if s.inHand > 0 {
		return false
	}

	// You get another turn if you end in your store
	if pitIndex == s.getStoreIndex(s.currentPlayer) {
		s.state = SelectingPit
	} else {
		s.state = IsWinner
	}

	if !onPlayersSide(pitIndex, s.currentPlayer) ||
		pitIndex == s.getStoreIndex(s.currentPlayer) ||
		s.board[pitIndex] != 1 ||
		s.board[GetOppositePit(pitIndex)] == 0 {
		return false
	}

	s.state = Stealing
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
	return (uint8(player) * 7)
}

func (s *Model) steal() {
	oppositePit := GetOppositePit(s.lastPlacedPit)
	s.board[s.getStoreIndex(s.currentPlayer)] += s.board[oppositePit] + 1
	s.board[oppositePit] = 0
	s.board[s.lastPlacedPit] = 0
	s.state = IsWinner
}

func GetOppositePit(pit uint8) uint8 {
	return oppositeIdxOffsets[pit]
}

func onPlayersSide(pit uint8, player Player) bool {
	lBound, uBound := getPlayerBounds(player)
	return pit >= lBound && pit <= uBound
}
