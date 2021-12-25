package migration

import (
	"fmt"
	"os"
)

func (m *Migration) UpOauth_clients() {
	query := `
		CREATE TABLE oauth_clients(
			id VARCHAR(255) NOT NULL PRIMARY KEY,
			secret VARCHAR(255) NOT NULL,
			domain VARCHAR(255) NOT NULL,
			data TEXT NOT NULL
		)
	`
	_, err := Connection().Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fmt.Println(string(Green), "success", string(Reset), "up 0001_migration_oauth_clients.go")
}

func (m *Migration) DownOauth_clients() {
	query := `DROP TABLE oauth_clients`
	_, err := Connection().Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fmt.Println(string(Green), "success", string(Reset), "down 0001_migration_oauth_clients.go")
}
