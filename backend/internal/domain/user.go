package domain

import (
	"time"
	"uuid"
)

type User struct {
	UserId    uuid.UUID
	Balance   int
	CreatedAt time.Time
}

type UserAccountResponse struct {
	Balance int            `json:"balance"`
	Games   []GameResponse `json:"games"`
}
