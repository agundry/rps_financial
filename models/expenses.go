package models

import (
	"../util"
)

type Expense struct {
	Id          int 	`json:"id"`
	AustinThrow	string 	`json:"austin_throw"`
	SamThrow 	string 	`json:"sam_throw"`
	Winner		string  `json:"winner"`
	Cost		int		`json:"cost"`
	CreatedAt	int64		`json:"created_at"`
}

func NewExpense(austinThrow util.Hand, samThrow util.Hand, cost int) Expense {
	var winner util.Name
	var outcome util.Outcome = util.PlayHand(austinThrow, samThrow)
	if outcome == util.P1Win {
		winner = util.AUSTIN
	} else if outcome == util.P2Win {
		winner = util.SAM
	} else {
	//	TODO throw error
	}

	return Expense{
		AustinThrow: austinThrow.String(),
		SamThrow:    samThrow.String(),
		Winner:      winner,
		Cost:        cost,
		CreatedAt:   util.GetEpochSeconds(),
	}
}

