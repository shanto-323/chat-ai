package dto

import "github.com/go-playground/validator"

type CreateUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (u *CreateUser) Validate() error {
	return validator.New().Struct(u)
}

type UserAuthResponse struct {
	AccessToken string `json:"accessToken"`
}
