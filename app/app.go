package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"

	"../util"
	"./models"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	Router   *mux.Router
	Database *sql.DB
}

func (app *App) SetupRouter() {
	app.Router.Methods("GET").Path("/expenses").HandlerFunc(app.IndexExpenses)
	app.Router.Methods("POST").Path("/expenses").HandlerFunc(app.CreateExpense)
	app.Router.Methods("GET").Path("/expenses/{expenseId}").HandlerFunc(app.GetExpense)
}

func (app *App) IndexExpenses(w http.ResponseWriter, r *http.Request) {
	rows, err := app.Database.Query("SELECT id, austin_throw, sam_throw, winner, cost, created_at from expenses")
	if err != nil {
		log.Fatal("Query Failed:", err.Error())
	}

	data := []models.Expense{}
	for rows.Next() {
		var e models.Expense
		err = rows.Scan(&e.Id, &e.AustinThrow, &e.SamThrow, &e.Winner, &e.Cost, &e.CreatedAt)
		if err != nil {
			log.Fatal("Query result parsing failed:", err.Error())
		}
		data = append(data, e)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (app *App) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var expenseRequest ExpenseRequest
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "New expenses need a hand for both players and a cost")
	}

	// TODO validate request inputs
	// Parse request into expense
	json.Unmarshal(reqBody, &expenseRequest)
	austinHand, err := util.HandFromString(expenseRequest.AustinThrow)
	samHand, err := util.HandFromString(expenseRequest.SamThrow)
	var expense = models.NewExpense(austinHand, samHand, expenseRequest.Cost)

	stmt, err := app.Database.Prepare("INSERT INTO `expenses` (austin_throw, sam_throw, winner, cost, created_at) VALUES (?, ?, ?, ?, ?)")
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

func (app *App) GetExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	expenseId := vars["expenseId"]

	expense := &models.Expense{}
	err := app.Database.QueryRow("SELECT id, austin_throw, sam_throw, winner, cost, created_at from expenses where id = ?", expenseId).Scan(&expense.Id, &expense.AustinThrow, &expense.SamThrow, &expense.Winner, &expense.Cost, &expense.CreatedAt)
	if err != nil {
		log.Fatal("Query Failed:", err.Error())
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expense)
}
