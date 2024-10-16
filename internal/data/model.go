package data

import (
	"database/sql"
)

type Model struct {
	Expenses ExpenseModel
	Users    UserModel
}

func New(db *sql.DB) Model {
	return Model{
		Expenses: ExpenseModel{DB: db},
		Users:    UserModel{DB: db},
	}
}
