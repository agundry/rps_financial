package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"

	"../db"
	"../util"
	"./models"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	Router   *mux.Router
	Database *sql.DB
}

func (app *App) SetupRouter() {
	app.Router.Methods("POST").Path("/expenses").HandlerFunc(CreateExpense)
}


func CreateExpense(w http.ResponseWriter, r *http.Request) {
	var expenseRequest ExpenseRequest
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "New expenses need a hand for both players and a cost")
	}

	// TODO validate request inputs
	// Parse request into expense
	json.Unmarshal(reqBody, &expenseRequest)
	austinHand, err := util.HandFromString(expenseRequest.AustinThrow)
	samHand, err := util.HandFromString(expenseRequest.AustinThrow)
	var expense = models.NewExpense(austinHand, samHand, expenseRequest.Cost)

	database, err := db.InitDbConnection()
	if err != nil {
		log.Fatal("Database connection failed:", err.Error())
	}

	stmt, err := database.Prepare("INSERT INTO `expenses` (austin_throw, sam_throw, winner, cost, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Database INSERT failed:", err.Error())
	}
	defer stmt.Close()

	res, err := stmt.Exec(expense.AustinThrow, expense.SamThrow, expense.Winner, expense.Cost, expense.CreatedAt)
	if err != nil {
		log.Fatal("Database INSERT failed:", err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	expense.Id = int(id)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(expense)
}
