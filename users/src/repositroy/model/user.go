package model

import (
	"time"

	"github.com/google/uuid"
)

const UsersTableName = "users"

type User struct {
	ID                   uuid.UUID `db:"id"`
	Name                 string    `db:"user_name"`
	FullName             string    `db:"full_name"`
	PasswordHash         []byte    `db:"password_hash"`
	Status               string    `db:"status_text"`
	RegitrationTimestamp time.Time `db:"timestamp"`
}

func (User) TableName() string {
	return UsersTableName
}