package main

import (
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

func (b *Board) Set(idx, val uint8) {
	b.board[idx] = val
}

func (b *Board) Get(idx uint8) uint8 {
	return b.board[idx]
}

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

func (b *Board) SelectPit(pit uint8) uint8 {
	numInPit := b.board[pit]

	if numInPit == 0 {
		panic("can't select an empty pit")
	}

	b.board[pit] = 0
	return numInPit
}

func (b *Board) MoveFromHandToPit(inHand, lastPlacedPit uint8, currentPlayer Player) (uint8, uint8) {
	// get next pit we're going to place a stone in
	var pitIndex uint8
	logger.Info("lastPlacedPit", lastPlacedPit)
	if lastPlacedPit == 0 {
		pitIndex = uint8(len(b.board) - 1)
	} else {
		pitIndex = lastPlacedPit - 1
	}
	logger.Info("pitIndex", pitIndex)

	otherPlayer := Player((currentPlayer + 1) % 2)
	otherStore := b.getStoreIndex(otherPlayer)

	if inHand > 0 {
		// Rule: We don't place stones in the other players store
		// TODO: Can I refactor this to make it more simple?
		if pitIndex == otherStore {
			if pitIndex == 0 {
				pitIndex = uint8(len(b.board) - 1)
			} else {
				pitIndex--
			}
		}
		inHand--
		b.board[pitIndex]++
		lastPlacedPit = pitIndex
	}
	return inHand, lastPlacedPit
}

func (b *Board) Steal(currentPlayer Player, lastPlacedPit uint8) {
	oppositePit := b.GetOppositePit(lastPlacedPit)
	b.board[b.getStoreIndex(currentPlayer)] += b.board[oppositePit] + 1
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

func (s Model) getNextNonEmptyPit() string {
	return "hello"
}

func (b Board) GetLastPlacedPit() uint8 {
	return b.lastPlacedPit
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
