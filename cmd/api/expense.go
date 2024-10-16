package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Mensurui/expenseAPI/internal/data"
)

func (app *application) readExpenses(w http.ResponseWriter, r *http.Request) {
	id, err := app.readID(r)
	if err != nil {
		app.l.Println("Error reading id: ", err)
		return
	}
	data, err := app.models.Expenses.Get(id)

	if err != nil {
		app.l.Println("Error fetching data: ", err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.l.Println("Error writing data: ", err)
		return
	}
}

func (app *application) readExpenseAll(w http.ResponseWriter, r *http.Request) {
	expenses, err := app.models.Expenses.GetAll()

	if err != nil {
		app.l.Println("Error getting the data: ", err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, expenses, nil)
	if err != nil {
		app.l.Println("Error writting the json", err)
		return
	}
}

func (app *application) createExpenses(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		app.l.Println("Error Decoding")
		return
	}

	expense := &data.Expense{
		Name:      input.Name,
		Price:     input.Price,
		CreatedAt: time.Now(),
	}

	err = app.models.Expenses.Insert(expense)
	if err != nil {
		app.l.Println("Error inserting data:", err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, expense, nil)
	if err != nil {
		app.l.Println("Error writing json")
	}
}
