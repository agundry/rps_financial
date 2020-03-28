package models

import (
	"../../util"
	"database/sql"
	"log"
)

type Expense struct {
	Id          int 	`json:"id"`
	AustinThrow	string 	`json:"austin_throw"`
	SamThrow 	string 	`json:"sam_throw"`
	Winner		string  `json:"winner"`
	Cost		int		`json:"cost"`
	CreatedAt	int64		`json:"created_at"`
}

func ConstructExpense(austinThrow util.Hand, samThrow util.Hand, cost int) Expense {
	var winner util.Name
	var outcome = util.PlayHand(austinThrow, samThrow)
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

func (e *Expense) GetExpense(db *sql.DB) error {
	return db.QueryRow("SELECT austin_throw, sam_throw, winner, cost, created_at from expenses where id = ?",
		e.Id).Scan(&e.AustinThrow, &e.SamThrow, &e.Winner, &e.Cost, &e.CreatedAt)
}

func (e *Expense) DeleteExpense(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM expenses WHERE id=?", e.Id)

	return err
}

func (e *Expense) InsertExpense(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO `expenses` (austin_throw, sam_throw, winner, cost, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Database INSERT failed:", err.Error())
	}
	defer stmt.Close()

	res, err := stmt.Exec(e.AustinThrow, e.SamThrow, e.Winner, e.Cost, e.CreatedAt)
	if err != nil {
		log.Fatal("Database INSERT failed:", err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	e.Id = int(id)

	return err
}

func GetExpenses(db *sql.DB) ([]Expense, error) {
	rows, err := db.Query("SELECT id, austin_throw, sam_throw, winner, cost, created_at from expenses")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data := []Expense{}
	for rows.Next() {
		var e Expense
		err = rows.Scan(&e.Id, &e.AustinThrow, &e.SamThrow, &e.Winner, &e.Cost, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		data = append(data, e)
	}

	return data, nil
}
