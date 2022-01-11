package helper

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/danangkonang/oauth2-golang/config"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
)

// NewClientStore create client store
func NewClientStore(db *config.DB) *ClientStore {
	return &ClientStore{
		db:   db.Db,
		data: make(map[string]oauth2.ClientInfo),
	}
}

// ClientStore client information store
type ClientStore struct {
	db *sql.DB
	sync.RWMutex
	data map[string]oauth2.ClientInfo
}

type ClientStoreItem struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	Domain string `json:"domain"`
	Data   string `json:"data"`
}

func (cs *ClientStore) toClientInfo(data string) (oauth2.ClientInfo, error) {
	var cm models.Client
	err := json.Unmarshal([]byte(data), &cm)
	return &cm, err
}

// GetByID according to the ID for the client information
func (cs *ClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	if id == "" {
		return nil, nil
	}

	var item ClientStoreItem
	err := cs.db.QueryRow(fmt.Sprintf("SELECT data FROM %s WHERE id = ?", "oauth_clients"), id).Scan(&item.Data)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}

	return cs.toClientInfo(item.Data)
}

// Set set client information
func (cs *ClientStore) Set(id string, cli oauth2.ClientInfo) (err error) {
	data, err := json.Marshal(cli)
	if err != nil {
		return err
	}

	_, err = cs.db.Exec(fmt.Sprintf("INSERT INTO %s (id, secret, domain, data) VALUES (?,?,?,?)", "oauth_clients"),
		cli.GetID(),
		cli.GetSecret(),
		cli.GetDomain(),
		string(data),
	)
	if err != nil {
		return err
	}
	return nil
}
