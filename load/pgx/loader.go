package pgx

import (
	"database/sql"
	"fmt"
	"gitlab.com/dipper-iot/shared/cli"
	"gitlab.com/dipper-iot/shared/service"
	"os"
	"strings"
)

type PG struct {
	db     *sql.DB
	prefix string
}

const DbRead string = "REPL"
const DbWrite string = ""

func NewPG(prefix string) *PG {
	if len(prefix) > 0 {
		prefix = fmt.Sprintf("%s_", strings.ToUpper(prefix))
	}
	return &PG{prefix: prefix}
}

func (p *PG) Name() string {
	return "pg"
}

func (p *PG) Priority() int {
	return 1
}

func (p PG) Flags() []cli.Flag {
	return nil
}

func (p *PG) DB() *sql.DB {
	return p.db
}

func (p *PG) Start(o *service.Options, c *cli.Context) error {
	var err error

	username := os.Getenv(fmt.Sprintf("%sDB_USER", p.prefix))
	password := os.Getenv(fmt.Sprintf("%sDB_PASSWORD", p.prefix))
	dbName := os.Getenv(fmt.Sprintf("%sDB_NAME", p.prefix))
	dbHost := os.Getenv(fmt.Sprintf("%sDB_HOST", p.prefix))
	dbPort := os.Getenv(fmt.Sprintf("%sDB_PORT", p.prefix))
	schema := os.Getenv(fmt.Sprintf("%sDB_SCHEMA", p.prefix))

	p.db, err = NewDatabase(&OptionPg{
		Host:     dbHost,
		Port:     dbPort,
		Password: password,
		Database: dbName,
		Username: username,
		Schema:   schema,
	})

	return err
}

func (p *PG) Stop() error {
	if p.db == nil {
		return nil
	}
	return p.db.Close()
}
