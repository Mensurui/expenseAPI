package data

import (
	"database/sql"
	"fmt"
	"time"
)

type Expense struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type ExpenseModel struct {
	DB *sql.DB
}

func (e ExpenseModel) Insert(expense *Expense) error {
	query := `
    INSERT INTO expense(name, price, created_at)
    VALUES($1, $2, $3)
    RETURNING id, name, price, created_at
    `
	args := []interface{}{expense.Name, expense.Price, expense.CreatedAt}

	// Update the Scan method to receive all four values
	err := e.DB.QueryRow(query, args...).Scan(&expense.ID, &expense.Name, &expense.Price, &expense.CreatedAt)
	if err != nil {
		// Log the error with a more descriptive message
		return fmt.Errorf("failed to insert expense: %w", err)
	}
	return nil
}

func (e ExpenseModel) Get(id int64) (*Expense, error) {
	query := `
	SELECT id, name, price, created_at
	FROM expense 
	WHERE id=$1
	`
	var expense Expense

	err := e.DB.QueryRow(query, id).Scan(
		&expense.ID,
		&expense.Name,
		&expense.Price,
		&expense.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &expense, nil
}

func (e ExpenseModel) GetAll() ([]*Expense, error) {
	query := `
	SELECT id, name, price, created_at
	FROM expense
	`

	rows, err := e.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := []*Expense{}

	for rows.Next() {
		var expense Expense
		err := rows.Scan(
			&expense.ID,
			&expense.Name,
			&expense.Price,
			&expense.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, &expense)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}
