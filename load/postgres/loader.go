package postgres

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"gitlab.com/dipper-iot/shared/cli"
	"gitlab.com/dipper-iot/shared/service"
	"os"
	"strings"
)

type Postgres struct {
	db     *pg.DB
	Config *Config
	prefix string
}

func (p *Postgres) Flags() []cli.Flag {
	return nil
}

func (p *Postgres) Priority() int {
	return 1
}

const DbRead string = "REPL"
const DbWrite string = ""

func (p *Postgres) Start(o *service.Options, c *cli.Context) error {

	dbHost := os.Getenv(fmt.Sprintf("%sDB_HOST", p.prefix))
	dbPort := os.Getenv(fmt.Sprintf("%sDB_PORT", p.prefix))

	p.Config.User = os.Getenv(fmt.Sprintf("%sDB_USER", p.prefix))
	p.Config.Password = os.Getenv(fmt.Sprintf("%sDB_PASSWORD", p.prefix))
	p.Config.Database = os.Getenv(fmt.Sprintf("%sDB_NAME", p.prefix))
	p.Config.Addr = fmt.Sprintf("%s:%s", dbHost, dbPort)
	p.Config.Schema = os.Getenv(fmt.Sprintf("%sDB_SCHEMA", p.prefix))

	var err error
	p.db = ConnectPostgres(p.Config, o.Name)
	err = p.db.Ping(context.TODO())
	if err != nil {
		return err
	}

	return nil
}

func NewPostgres(prefix string) *Postgres {
	if len(prefix) > 0 {
		prefix = fmt.Sprintf("%s_", strings.ToUpper(prefix))
	}
	return &Postgres{
		prefix: prefix,
		Config: &Config{},
	}
}

func (p *Postgres) Database() *pg.DB {
	return p.db
}

func (p Postgres) Name() string {
	return "postgres"
}

func (p *Postgres) Stop() error {
	return p.db.Close()
}
