package main

import tea "github.com/charmbracelet/bubbletea"

func (s State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Flow:
	// State: SelectPit
	//   - MoveLeft, MoveRight, SelectPit
	// State: SelectPit -> MoveFromHandToPit
	//   - While inHand -> MoveFromHandToPit
	// State: MoveFromHandToPit -> IsWinner
	// State IsWinner -> SelectPit (Switch player)
	switch s.state {
	case SelectingPit:
	case MovingFromHandToPit:
	case IsWinner:
	case SwitchPlayer:
	case GameOver:
	}

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
		}

	}
	return s, nil
}

func (s *State) MoveRight() {
	// TODO: handle moving past empty pits
	lBound, uBound := getPlayerBounds(s.currentPlayer)
	if s.selectedPit < uBound {
		s.selectedPit++
	} else {
		s.selectedPit = lBound
	}
}

func (s *State) MoveLeft() {
	// TODO: handle moving past empty pits
	lBound, uBound := getPlayerBounds(s.currentPlayer)
	if s.selectedPit > lBound {
		s.selectedPit--
	} else {
		s.selectedPit = uBound
	}
}

// getPlayerBounds returns the lower and upper movable bounds of the pits for a given player.
func getPlayerBounds(p Player) (lBound, uBound uint8) {
	lBound = uint8(p) * 7
	uBound = lBound + 5
	return
}

func (s *State) SelectPit() {
	numInPit := s.board[s.selectedPit]

	if numInPit == 0 {
		panic("can't select an empty pit")
	}

	s.board[s.selectedPit] = 0
	s.inHand = numInPit
}

func (s *State) moveFromHandToPit() {

}
