package model

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Websession contains all information regarding a single web session.
type Websession struct {
	Token          uuid.UUID
	AccountId    int
	ExpireAt time.Time
}

func (ws *Websession) String() string {
	return fmt.Sprintf("Token=%s, AccountId=%d, ExpireAt=%s",
		ws.Token, ws.AccountId, ws.ExpireAt.String())
}

// CreateWebsession creates and stores a new web session to the database.
func CreateWebsession(dbpool *pgxpool.Pool, accountId int, hoursLoggedIn uint32) (*Websession, error) {
	ws := Websession{uuid.New(), accountId, time.Now().Add(time.Hour * time.Duration(hoursLoggedIn))}
	query := "INSERT INTO websession (token, account_id, expire_at) VALUES ($1, $2, $3)"
	_, err := dbpool.Exec(context.Background(), query, ws.Token, ws.AccountId, ws.ExpireAt)
	if err != nil {
		return nil, err
	}
	return &ws, nil
}
