package entity

import "github.com/shanto-323/chat-ai/model"

type User struct {
	model.Base

	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
}
