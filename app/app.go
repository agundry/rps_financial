package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/agundry/rps_financial/app/models"
	"github.com/agundry/rps_financial/util"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App struct {
	Router   *mux.Router
	Database *sql.DB
	Logger 	 *util.StandardLogger
}

func (app *App) Initialize(user, password, serverName, dbName string) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, serverName, dbName)

	var err error
	app.Database, err = sql.Open("mysql", connectionString)

	// the err value at this point is not enough to check for a valid connection
	// err would only be non-null if there was an issue with the driver
	err = app.Database.Ping()

	if err != nil {
		app.Logger.Fatal(err)
	}
	app.Logger.Infof("Successfully connected to mysql database at %s", serverName)

	app.Router = mux.NewRouter().StrictSlash(true)

	app.initializeRoutes()
	app.Logger.Infof("Initialized Routes")
}

func (app *App) WithLogging(logger *util.StandardLogger) {
	app.Logger = logger
}

func (app *App) Run(addr string) {
	// Start the app, logging all requests
	app.Logger.Fatal(http.ListenAndServe(addr, util.LogRequest(app.Logger, app.Router)))
}

func (app *App) initializeRoutes() {
	app.Logger.Infof("Initializing routes")

	// Access to prometheus metrics
	app.Router.Methods("GET").Path("/metrics").Handler(promhttp.Handler())

	// Healthcheck for kubernetes
	app.Router.Methods("GET").Path("/healthcheck").HandlerFunc(app.healthcheck)

	// Application routes
	app.Router.Methods("GET").Path("/expenses").HandlerFunc(app.indexExpenses)
	app.Router.Methods("POST").Path("/expenses").HandlerFunc(app.createExpense)
	app.Router.Methods("GET").Path("/expenses/{id:[0-9]+}").HandlerFunc(app.getExpense)
	app.Router.Methods("DELETE").Path("/expenses/{id:[0-9]+}").HandlerFunc(app.deleteExpense)

	app.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		met, _ := route.GetMethods()
		app.Logger.Infof("Route: %s %s", tpl, met)
		return nil
	})
}

func (app *App) healthcheck(w http.ResponseWriter, r *http.Request) {
	rows, err := app.Database.Query("SELECT 1")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, rows)
}

func (app *App) indexExpenses(w http.ResponseWriter, r *http.Request) {
	expenses, err := models.GetExpenses(app.Database)
	if err != nil {
		app.Logger.InternalServerError(err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, expenses)
}

func (app *App) createExpense(w http.ResponseWriter, r *http.Request) {
	var expenseRequest ExpenseRequest
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.Logger.BadRequestError("New expenses need a hand for both players and a cost")
	}

	// TODO validate request inputs
	// Parse request into expense
	json.Unmarshal(reqBody, &expenseRequest)
	austinHand, err := util.HandFromString(expenseRequest.AustinThrow)
	samHand, err := util.HandFromString(expenseRequest.SamThrow)
	var expense = models.ConstructExpense(austinHand, samHand, expenseRequest.Cost)

	if err := expense.InsertExpense(app.Database); err != nil {
		app.Logger.InternalServerError(err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, expense)
}

func (app *App) getExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.Logger.BadRequestError("Invalid expense ID")
		respondWithError(w, http.StatusBadRequest, "Invalid expense ID")
		return
	}

	expense := models.Expense{Id: id}
	if err := expense.GetExpense(app.Database); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Expense not found")
		default:
			app.Logger.InternalServerError(err.Error())
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
		app.Logger.BadRequestError("Invalid expense ID")
		respondWithError(w, http.StatusBadRequest, "Invalid Expense ID")
		return
	}

	e := models.Expense{Id: id}
	if err := e.DeleteExpense(app.Database); err != nil {
		app.Logger.InternalServerError(err.Error())
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