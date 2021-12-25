package migration

import (
	"fmt"
	"os"
)

func (m *Migration) UpOauth_tokens() {
	query := `
		CREATE TABLE oauth_tokens(
			id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			code VARCHAR(255),
			access VARCHAR(255) NOT NULL,
			refresh VARCHAR(255) NOT NULL,
			data TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			expired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			KEY access_k(access),
			KEY refresh_k (refresh),
			KEY expired_at_k (expired_at),
			KEY code_k (code)
		)
	`
	_, err := Connection().Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fmt.Println(string(Green), "success", string(Reset), "up 0002_migration_oauth_tokens.go")
}

func (m *Migration) DownOauth_tokens() {
	query := `DROP TABLE oauth_tokens`
	_, err := Connection().Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fmt.Println(string(Green), "success", string(Reset), "down 0002_migration_oauth_tokens.go")
}
