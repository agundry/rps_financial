package util

import (
	"errors"
	"fmt"
)

type Hand int
const (
	ROCK Hand = iota
	PAPER
	SCISSORS
)

type Outcome int
const (
	P1Win Outcome = iota
	P2Win
	DRAW
)

func (h Hand) String() string {
	return [...]string{"ROCK", "PAPER", "SCISSORS"}[h]
}

func HandFromString(in string) (Hand, error) {
	if in == "ROCK" {
		return ROCK, nil
	} else if in == "PAPER" {
		return PAPER, nil
	} else if in == "SCISSORS" {
		return SCISSORS, nil
	} else {
		return 0, errors.New(fmt.Sprintf("Invalid input for Hand: %s", in))
	}
}

func PlayHand(p1 Hand, p2 Hand) Outcome {
	if p1 == p2 {
		return DRAW
	} else if (p1 == ROCK && p2 == SCISSORS) || (p1 == SCISSORS && p2 == PAPER ) || (p1 == PAPER && p2 == ROCK) {
		return P1Win
	} else {
		return P2Win
	}
}
