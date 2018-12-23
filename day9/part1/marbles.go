package main

import (
	"container/list"
	"fmt"
)

const PlayersCount int = 465
const LastMarbleWorth int = 71940

type Marble int

type MarblesCircle struct {
	Current *list.Element
	Next    Marble
	Marbles *list.List
}

func NewMarblesCircle() *MarblesCircle {
	circle := &MarblesCircle{
		Marbles: list.New(),
	}

	circle.Marbles.PushFront(Marble(0))
	circle.Next = Marble(1)
	circle.Current = circle.Marbles.Front()

	return circle
}

func (circle *MarblesCircle) InsertClockWise(marble Marble) {
	// circle is NOT large enough to insert new marble between two others
	if circle.Marbles.Len() == 1 {
		circle.Marbles.PushBack(marble)
		circle.Current = circle.Current.Next()
		return
	}

	inBetween := circle.Current.Next()
	if inBetween == nil {
		// we arrived at the end of the list, let's circle back to the front
		inBetween = circle.Marbles.Front()
	}

	circle.Marbles.InsertAfter(marble, inBetween)
	circle.Current = inBetween.Next()
}

func (circle *MarblesCircle) RemoveSeventhCounterClockwise() Marble {
	var removedMarble Marble

	current := circle.Current
	for i := 7; i > 0; i-- {
		if current.Prev() == nil {
			current = circle.Marbles.Back()
			continue
		}

		current = current.Prev()
	}

	removedMarble = current.Value.(Marble)

	if current.Next() != nil {
		circle.Current = current.Next()
	} else {
		circle.Current = circle.Marbles.Front()
	}

	circle.Marbles.Remove(current)

	return removedMarble
}

func (circle *MarblesCircle) Print() {
	for marble := circle.Marbles.Front(); marble != nil; marble = marble.Next() {
		if circle.Current == marble {
			fmt.Printf(" (%d)", marble.Value)
		} else {
			fmt.Printf(" %d", marble.Value)
		}
	}

	fmt.Println()
}

func main() {
	// setup the scores
	scores := make(map[int]int)
	for i := 1; i <= PlayersCount; i++ {
		scores[i] = 0
	}

	// setup the marbles circle
	circle := NewMarblesCircle()

	// play!
	player := 1

	circle.Print()
	for int(circle.Next) <= LastMarbleWorth {
		toInsert := Marble(circle.Next)

		//fmt.Printf("Current player is %d\n", player)
		//fmt.Printf("Marble to insert is %d\n", toInsert)

		// turn
		if circle.Next%23 == 0 {
			// weird stuff
			scores[player] += int(toInsert)

			removedMarble := circle.RemoveSeventhCounterClockwise()
			scores[player] += int(removedMarble)
		} else {
			// classical turn
			circle.InsertClockWise(toInsert)
		}
		//circle.Print()

		// setup the next turn
		circle.Next++
		player++

		if player > PlayersCount {
			player = 1
		}
	}

	fmt.Println("\nScores:")
	for player := range scores {
		fmt.Printf("→ Player %d: %d\n", player, scores[player])
	}

	winner := -1
	maxScore := -1
	for player := range scores {
		if scores[player] > maxScore {
			winner = player
			maxScore = scores[player]
		}
	}

	fmt.Printf("\nWinner: player %d with %d points\n", winner, maxScore)
}
