package data

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("used email")
	ErrRecordNotFound = errors.New("email or password not correct")
)

type User struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password password `json:"-"`
}

type password struct {
	hash      []byte
	plainText *string
}

func (p *password) Set(plainTextPassword string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return nil
	}

	p.plainText = &plainTextPassword
	p.hash = hashed

	return nil
}

func (p *password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}

	}
	return true, nil
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users(username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id, username, email
	`
	args := []interface{}{user.Username, user.Email, user.Password.hash}
	err := u.DB.QueryRow(query, args...).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (u *UserModel) GetByEmail(email string, password string) (*User, error) {
	query := `
	SELECT id, username, email, password
	FROM users
	WHERE email = $1 
	`
	var user User

	err := u.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
