package main

import (
	"encoding/json"
	"net/http"

	"github.com/Mensurui/expenseAPI/internal/data"
)

func (app *application) registrationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"-"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.l.Println("Error decoding: ", err)
		return
	}

	user := data.User{
		Username: input.Username,
		Email:    input.Email,
	}

	err = user.Password.Set(input.Password)

	if err != nil {
		app.l.Println("Error: ", err)
	}

	err = app.models.Users.Insert(&user)
	if err != nil {
		app.l.Println("Error inserting: ", err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, user, nil)
	if err != nil {
		app.l.Println("Error writing: ", err)
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var input struct {
		Email    string `json:"email"`
		Password string `json:"-"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.l.Println("Error Decoding: ", err)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email, input.Password)
	if err != nil {
		app.l.Println("Error getting by email: ", err)
	}
	matches, err := user.Password.Matches(input.Password)

	if err != nil {
		app.l.Println("Error getting password: ", err)
		return
	}
	if !matches {
		app.l.Println(w, "Password and email mismatch")
		return
	}

	tokenString, err := app.createToken(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		app.l.Fatal("No username found")
	}
	err = app.writeJSON(w, http.StatusOK, tokenString, nil)
	if err != nil {
		app.l.Println("Error: ", err)
		return
	}
}
