package model

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Account contains all information regarding user accounts.
type Account struct {
	ID          int
	Username    string
	Password    string
	Created     time.Time
	LastUpdated time.Time
}

func (a *Account) String() string {
	return fmt.Sprintf("ID=%d, Username=%s, Password=%s, Created=%s, LastUpdated=%s",
		a.ID, a.Username, a.Password, a.Created.String(), a.LastUpdated.String())
}

// CreateAccount creates and stores a new account to the database.
func (a *Account) CreateAccount(dbpool *pgxpool.Pool) error {
	query := "INSERT INTO account (username, password, created, last_updated) VALUES ($1, $2, $3, $3)"
	_, err := dbpool.Exec(context.Background(), query, a.Username, a.Password, time.Now())
	if err != nil {
		// Assume username already taken
		log.Println(err)
		return err
	}
	return nil
}
