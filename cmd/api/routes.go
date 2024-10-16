package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	route := httprouter.New()

	route.HandlerFunc(http.MethodPost, "/expenses", app.createExpenses)
	route.HandlerFunc(http.MethodGet, "/expenses/:id", app.readExpenses)
	route.HandlerFunc(http.MethodGet, "/expenses", app.readExpenseAll)

	route.HandlerFunc(http.MethodPost, "/user", app.registrationHandler)
	route.HandlerFunc(http.MethodPost, "/login", app.loginHandler)
	return app.protected(route)
}
