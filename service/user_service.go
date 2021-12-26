package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/danangkonang/oauth2/config"
	"github.com/danangkonang/oauth2/model"
)

type UserService interface {
	Register(m *model.UserRegister) error
}

func NewServiceUser(Con *config.DB) UserService {
	return &connUser{
		Psql: Con.Db,
	}
}

type connUser struct {
	Psql *sql.DB
}

func (c *connUser) Register(m *model.UserRegister) error {
	query := `
		INSERT INTO users(client_id, user_name, password, created_at, updated_at) VALUES(?,?,?,?,?)
	`
	_, err := c.Psql.Exec(query, m.ClientId, m.UserName, m.Password, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("INTERNAL_SERVER_ERROR")
	}
	return nil
}
