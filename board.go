package main

import (
	"fmt"
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

var oppositeIdxOffsets [14]uint8 = [14]uint8{7, 13, 12, 11, 10, 9, 8, 0, 6, 5, 4, 3, 2, 1}

type Board struct {
	board         [14]uint8
	lastPlacedPit uint8
}

func NewBoard() Board {
	return Board{
		board: [14]uint8{0, 4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 4},
	}
}

func NewBoardWithState(state [14]uint8) Board {
	return Board{board: state}
}

func NewBoardWithOverrideState(state map[uint8]uint8) Board {
	newState := [14]uint8{0, 4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 4}
	for idx, val := range state {
		newState[idx] = val
	}
	return NewBoardWithState(newState)
}

// Get returns the numver of stones at the pit index
func (b Board) Get(pit uint8) uint8 {
	return b.board[pit]
}

// MoveLeft moves the selected pit to the left, skipping empty pits
func (b Board) MoveLeft(selectedPit uint8, currentPlayer Player) uint8 {
	lBound, uBound := getPlayerBounds(currentPlayer)
	for range 6 {
		if currentPlayer == P2 {
			if selectedPit < uBound {
				selectedPit++
			} else {
				selectedPit = lBound
			}
		} else {
			if selectedPit > lBound {
				selectedPit--
			} else {
				selectedPit = uBound
			}
		}
		if b.board[selectedPit] != 0 {
			break
		}
	}
	return selectedPit
}

// MoveRight moves the selected pit to the right, skipping empty pits
func (m *Board) MoveRight(selectedPit uint8, currentPlayer Player) uint8 {
	lBound, uBound := getPlayerBounds(currentPlayer)
	for range 6 {
		if currentPlayer == P2 {
			if selectedPit > lBound {
				selectedPit--
			} else {
				selectedPit = uBound
			}
		} else {
			if selectedPit < uBound {
				selectedPit++
			} else {
				selectedPit = lBound
			}
		}
		if m.board[selectedPit] != 0 {
			break
		}
	}
	return selectedPit
}

// SelectPit removes all stones from a pit and returns the number of stones
func (b *Board) SelectPit(pit uint8) uint8 {
	numInPit := b.board[pit]

	if numInPit == 0 {
		panic("can't select an empty pit")
	}

	b.board[pit] = 0
	return numInPit
}

// MoveFromHandToPit moves a stone from the hand to the pit, returning number left in
// the hand and the last pit placed in
func (b *Board) MoveFromHandToPit(inHand, lastPlacedPit uint8, currentPlayer Player) (uint8, uint8) {
	if inHand == 0 {
		return 0, lastPlacedPit
	}
	nextPitIndex := b.getNextPit(lastPlacedPit, currentPlayer)
	inHand--
	b.board[nextPitIndex]++
	return inHand, nextPitIndex
}

// getNextPit returns the next pit to place a stone in based on the current pit and
// player
func (b Board) getNextPit(currentPit uint8, player Player) uint8 {
	var nextPit uint8
	if currentPit == 0 {
		nextPit = uint8(len(b.board) - 1)
	} else {
		nextPit = currentPit - 1
	}

	otherPlayer := Player((player + 1) % 2)
	otherStore := b.getStoreIndex(otherPlayer)

	if nextPit == otherStore {
		if nextPit == 0 {
			nextPit = uint8(len(b.board) - 1)
		} else {
			nextPit--
		}
	}
	return nextPit

}

func (b *Board) Steal(currentPlayer Player, lastPlacedPit uint8) {
	// TODO: Validate lastPlacedPit has 1 stone
	oppositePit := b.GetOppositePit(lastPlacedPit)
	toAdd := b.board[oppositePit] + b.board[lastPlacedPit]
	b.addToStore(currentPlayer, toAdd)
	b.board[oppositePit] = 0
	b.board[lastPlacedPit] = 0
}

func getPlayerBounds(p Player) (lBound, uBound uint8) {
	lBound = uint8(p)*7 + 1
	uBound = lBound + 5
	return
}

func (s Board) GetOppositePit(pit uint8) uint8 {
	return oppositeIdxOffsets[pit]
}

func (s Board) getStoreIndex(player Player) uint8 {
	return (uint8(player) * 7)
}

func (b Board) GetPlayerForPit(pit uint8) Player {
	if pit < 7 {
		return P1
	}
	return P2
}

func (b Board) IsPlayersStore(pit uint8, player Player) bool {
	return uint8(player)*7 == pit
}

func (b Board) GetFirstNonEmptyPit(player Player) uint8 {
	startPit := uint8(1)
	sign := 1
	if player == P2 {
		startPit = uint8(len(b.board) - 1)
		sign = -1
	}
	for i := range 6 {
		if pit := startPit + uint8((i * sign)); b.board[pit] != 0 {
			return pit
		}
	}
	panic("no stones in any of the player's pits")
}

func (b *Board) addToStore(player Player, numStones uint8) {
	b.board[b.getStoreIndex(player)] += numStones
}

func (b Board) GetNumInStore(player Player) uint8 {
	return b.board[b.getStoreIndex(player)]
}

// CollectRemainder moves remaining stones to that players side
func (b *Board) CollectRemainder() {
	p1, p2 := uint8(0), uint8(0)

	for i := range 6 {
		p1Index, p2Index := uint8(i+1), uint8(i+8)
		p1 += b.board[p1Index]
		b.board[p1Index] = 0

		p2 += b.board[p2Index]
		b.board[p2Index] = 0
	}

	if p1 == 0 {
		b.board[b.getStoreIndex(P2)] += p2
	} else if p2 == 0 {
		b.board[b.getStoreIndex(P1)] += p1
	} else {
		panic(fmt.Sprintf("Both players have stones left. p1: %d, p2: %d", p1, p2))
	}
}
