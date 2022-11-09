package pgx

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func NewDatabase(option *OptionPg) (*sql.DB, error) {
	if option.Schema == "" {
		option.Schema = "public"
	}
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s search_path=%s", option.Host, option.Username, option.Database, option.Password, option.Port, option.Schema) //Build connection string
	db, err := sql.Open("pgx", dbUri)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}
