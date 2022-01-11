package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/danangkonang/oauth2-golang/config"
	"github.com/danangkonang/oauth2-golang/model"
)

type UserService interface {
	Register(m *model.UserRegister) error
	Login(user string) (*model.Autorization, error)
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

func (c *connUser) Login(user string) (*model.Autorization, error) {
	usr := new(model.Autorization)
	query := `
		SELECT id, password from users WHERE user_name = ?
	`
	row := c.Psql.QueryRow(query, user)
	err := row.Scan(&usr.Id, &usr.Password)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("invalid username")
	case err != nil:
		return nil, errors.New("INTERNAL_SERVER_ERROR")
	}
	return usr, nil
}
