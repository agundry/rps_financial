package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/agundry/rps_financial/util"
	"github.com/agundry/rps_financial/app/models"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	Router   *mux.Router
	Database *sql.DB
}

func (app *App) Initialize(user, password, serverName, dbName string) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, serverName, dbName)

	var err error
	app.Database, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	app.Router = mux.NewRouter().StrictSlash(true)

	app.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (app *App) initializeRoutes() {
	app.Router.Methods("GET").Path("/expenses").HandlerFunc(app.indexExpenses)
	app.Router.Methods("POST").Path("/expenses").HandlerFunc(app.createExpense)
	app.Router.Methods("GET").Path("/expenses/{id:[0-9]+}").HandlerFunc(app.getExpense)
	app.Router.Methods("DELETE").Path("/expenses/{id:[0-9]+}").HandlerFunc(app.deleteExpense)
}

func (app *App) indexExpenses(w http.ResponseWriter, r *http.Request) {
	expenses, err := models.GetExpenses(app.Database)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, expenses)
}

func (app *App) createExpense(w http.ResponseWriter, r *http.Request) {
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
	var expense = models.ConstructExpense(austinHand, samHand, expenseRequest.Cost)

	if err := expense.InsertExpense(app.Database); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, expense)
}

func (app *App) getExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid expense ID")
		return
	}

	expense := models.Expense{Id: id}
	if err := expense.GetExpense(app.Database); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Expense not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, expense)
}

func (app *App) deleteExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Expense ID")
		return
	}

	e := models.Expense{Id: id}
	if err := e.DeleteExpense(app.Database); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}