package model

import "time"

type UserRegister struct {
	ClientId  string    `json:"client_id"`
	UserName  string    `json:"user_name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"crated_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Autorization struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}
