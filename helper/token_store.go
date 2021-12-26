package helper

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/danangkonang/oauth2/config"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
)

// NewMemoryTokenStore create a token store instance based on memory
func NewMysqlTokenStore(db *config.DB) (oauth2.TokenStore, error) {
	return NewFileTokenStore(db)
}

// NewFileTokenStore create a token store instance based on file
func NewFileTokenStore(db *config.DB) (oauth2.TokenStore, error) {
	return &TokenStore{db: db.Db}, nil
}

type TokenStore struct {
	db *sql.DB
}

type TokenStoreItem struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Code      string    `json:"code"`
	Access    string    `json:"access"`
	Refresh   string    `json:"refresh"`
	Data      string    `json:"data"`
}

// Create create and store the new token information
func (ts *TokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	buf, _ := json.Marshal(info)
	item := &TokenStoreItem{
		Data:      string(buf),
		CreatedAt: time.Now(),
	}

	if code := info.GetCode(); code != "" {
		item.Code = code
		item.ExpiredAt = info.GetCodeCreateAt().Add(info.GetCodeExpiresIn())
	} else {
		item.Access = info.GetAccess()
		item.ExpiredAt = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn())

		if refresh := info.GetRefresh(); refresh != "" {
			item.Refresh = info.GetRefresh()
			item.ExpiredAt = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn())
		}
	}

	_, err := ts.db.Exec(
		fmt.Sprintf("INSERT INTO %s (created_at, expired_at, code, access, refresh, data) VALUES (?,?,?,?,?,?)", "oauth_tokens"),
		item.CreatedAt,
		item.ExpiredAt,
		item.Code,
		item.Access,
		item.Refresh,
		item.Data)
	if err != nil {
		return err
	}
	return nil
}

// RemoveByCode use the authorization code to delete the token information
func (ts *TokenStore) RemoveByCode(ctx context.Context, code string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE code=? LIMIT 1", "oauth_tokens")
	_, err := ts.db.Exec(query, code)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	return err
}

// RemoveByAccess use the access token to delete the token information
func (ts *TokenStore) RemoveByAccess(ctx context.Context, access string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE access=? LIMIT 1", "oauth_tokens")
	_, err := ts.db.Exec(query, access)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	return err
}

// RemoveByRefresh use the refresh token to delete the token information
func (ts *TokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE refresh=? LIMIT 1", "oauth_tokens")
	_, err := ts.db.Exec(query, refresh)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (ts *TokenStore) getData(data string) (oauth2.TokenInfo, error) {
	var ti oauth2.TokenInfo
	var tm models.Token
	json.Unmarshal([]byte(data), &tm)
	ti = &tm
	return ti, nil
}

// GetByCode use the authorization code for token information data
func (ts *TokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	if code == "" {
		return nil, nil
	}

	query := fmt.Sprintf("SELECT data FROM %s WHERE code=? LIMIT 1", "oauth_tokens")
	var item TokenStoreItem
	err := ts.db.QueryRow(query, code).Scan(&item.Data)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}

	return ts.getData(item.Data)
}

// GetByAccess use the access token for token information data
func (ts *TokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	if access == "" {
		return nil, nil
	}

	query := fmt.Sprintf("SELECT data FROM %s WHERE access=? LIMIT 1", "oauth_tokens")
	var item TokenStoreItem
	err := ts.db.QueryRow(query, access).Scan(&item.Data)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}
	return ts.getData(item.Data)
}

// GetByRefresh use the refresh token for token information data
func (ts *TokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	if refresh == "" {
		return nil, nil
	}

	query := fmt.Sprintf("SELECT data FROM %s WHERE refresh=? LIMIT 1", "oauth_tokens")
	var item TokenStoreItem
	err := ts.db.QueryRow(query, refresh).Scan(&item.Data)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}
	return ts.getData(item.Data)
}
