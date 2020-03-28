package app_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/agundry/rps_financial/app"
)

var application app.App

func TestMain(m *testing.M) {

	application.Initialize(
		"root",
		"my-secret-pw",
		"localhost:13306",
		"rps")

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := application.Database.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	application.Database.Exec("DELETE FROM expenses")
	application.Database.Exec("ALTER TABLE expenses AUTO_INCREMENT = 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS expenses 
(
  id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  austin_throw varchar(10) NOT NULL,
  sam_throw varchar(10) NOT NULL,
  winner varchar(10) NOT NULL,
  cost int(11) NOT NULL COMMENT 'value in cents',
  created_at int(11) NOT NULL COMMENT 'time in epoch seconds',
  PRIMARY KEY (id)
)`

func TestIndexExpensesEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/expenses", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentExpense(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/expenses/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Expense not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Expense not found'. Got '%s'", m["error"])
	}
}

func TestCreateExpense(t *testing.T) {

	clearTable()

	var jsonStr = []byte(`{"austin_throw":"ROCK", "sam_throw":"PAPER", "cost": 1122}`)
	req, _ := http.NewRequest("POST", "/expenses", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["winner"] != "SAM" {
		t.Errorf("Expected winner to be 'Sam'. Got '%v'", m["winner"])
	}

	if m["cost"] != 1122.0 {
		t.Errorf("Expected expense cost to be '1122'. Got '%v'", m["cost"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected expense ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetExpense(t *testing.T) {
	clearTable()
	addExpenses(1)

	req, _ := http.NewRequest("GET", "/expenses/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addExpenses(1)

	req, _ := http.NewRequest("GET", "/expenses/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/expenses/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/expenses/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func addExpenses(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		application.Database.Exec("INSERT INTO expenses(austin_throw, sam_throw, winner, cost, created_at) VALUES(?, ?, ?, ?, ?)",
			"ROCK", "PAPER", "SAM", i, (i+1)*10)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	application.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}