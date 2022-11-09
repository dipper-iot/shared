package postgres

import (
	"context"
	"github.com/go-pg/pg/v10"
)

func ConnectPostgres(postgres *Config, name string) *pg.DB {

	schemaName := "public"
	if postgres.Schema != "" {
		schemaName = postgres.Schema
	}

	db := pg.Connect(&pg.Options{
		Addr:            postgres.Addr,
		User:            postgres.User,
		Password:        postgres.Password,
		Database:        postgres.Database,
		ApplicationName: name,
		OnConnect: func(c context.Context, conn *pg.Conn) error {
			_, err := conn.ExecContext(c, "set search_path=?", schemaName)
			if err != nil {
				return err
			}
			return nil
		},
	}).WithParam("search_path", schemaName)

	return db
}
