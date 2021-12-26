package migration

import (
	"fmt"
	"os"
)

func (m *Migration) UpUsers() {
	query := `
		CREATE TABLE users(
			id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			user_name VARCHAR (225) NOT NULL,
			password VARCHAR (225) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := Connection().Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fmt.Println(string(Green), "success", string(Reset), "up 0003_migration_users.go")
}

func (m *Migration) DownUsers() {
	query := `DROP TABLE users`
	_, err := Connection().Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fmt.Println(string(Green), "success", string(Reset), "down 0003_migration_users.go")
}
