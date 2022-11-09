package mongodb

import (
	"context"
	"fmt"
	"gitlab.com/dipper-iot/shared/cli"
	"gitlab.com/dipper-iot/shared/service"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"strings"
)

type Mongo struct {
	Config *MongoConfig
	db     *mongo.Database
	client *mongo.Client
	prefix string
}

const DbRead string = "REPL"
const DbWrite string = ""

func (r *Mongo) Flags() []cli.Flag {
	return nil
}

func (r *Mongo) Priority() int {
	return 1
}

func (r *Mongo) Start(o *service.Options, c *cli.Context) error {
	var (
		err error
	)

	uri := os.Getenv(fmt.Sprintf("%sMONGO_URI", r.prefix))
	database := os.Getenv(fmt.Sprintf("%sMONGO_DB", r.prefix))

	r.Config = &MongoConfig{
		Database: database,
		Uri:      uri,
	}

	err, r.db, r.client = ConnectDatabaseError(r.Config)
	return err
}

func NewMongo(prefix string) *Mongo {
	if len(prefix) > 0 {
		prefix = fmt.Sprintf("%s_", strings.ToUpper(prefix))
	}
	return &Mongo{
		Config: &MongoConfig{},
		prefix: prefix,
	}
}

func (r *Mongo) Name() string {
	return "mongodb"
}

func (r *Mongo) Stop() error {
	return r.client.Disconnect(context.TODO())
}

func (r *Mongo) Client() *mongo.Client {
	return r.client
}

func (r *Mongo) Database() *mongo.Database {
	return r.db
}
