package models_test

import (
	"../../util"
	"../models"
	"testing"
)

func TestConstructExpense(t *testing.T) {
	tables := []struct {
		a util.Hand
		s util.Hand
		c int
		w string
	}{
		{util.ROCK, util.PAPER, 1234, util.SAM},
		{util.SCISSORS, util.PAPER, 2345, util.AUSTIN},
		{util.ROCK, util.SCISSORS, 3456, util.AUSTIN},
	}

	for _, table := range tables {
		expense := models.ConstructExpense(table.a, table.s, table.c)
		if expense.Winner != table.w {
			t.Errorf("Winner of (%s+%s) was incorrect, got: %s, want: %s.", table.a, table.s, expense.Winner, table.w)
		}
	}
}