package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"

	"./api"
	"./models"
	"./util"
)

var allExpenses []models.Expense

func createExpense(w http.ResponseWriter, r *http.Request) {
	var expenseRequest api.ExpenseRequest
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "New expenses need a hand for both players and a cost")
	}

	// TODO validate request inputs
	// Parse request into expense
	json.Unmarshal(reqBody, &expenseRequest)
	austinHand, err := util.HandFromString(expenseRequest.AustinThrow)
	samHand, err := util.HandFromString(expenseRequest.AustinThrow)
	// TODO handle errors
	var expense models.Expense = models.NewExpense(austinHand, samHand, expenseRequest.Cost)
	allExpenses = append(allExpenses, expense)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(allExpenses)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/expenses", createExpense)
	log.Fatal(http.ListenAndServe(":8080", router))
}
