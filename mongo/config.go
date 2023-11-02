package mongo

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type config struct {
	URI                 string
	DBName              string
	CertificatePath     string
	IDType              IDType
	AutoPreload         bool
	ClearEmbeddedFields bool
	Driver              *driver
}

type driver struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func newDefaultConfig() *config {
	return &config{
		IDType:              ObjectID,
		AutoPreload:         true,
		ClearEmbeddedFields: true,
	}
}

func NewConfig(uri string, dbName string, certificatePath string) *config {
	cfg := newDefaultConfig()
	cfg.URI = uri
	cfg.DBName = dbName
	cfg.CertificatePath = certificatePath
	return cfg
}

func NewConfigString(uri string, dbName string, certificatePath string) *config {
	cfg := newDefaultConfig()
	cfg.URI = uri
	cfg.DBName = dbName
	cfg.CertificatePath = certificatePath
	cfg.IDType = String
	return cfg
}

func NewConfigWithClient(client *mongo.Client, database *mongo.Database) *config {
	cfg := newDefaultConfig()
	cfg.Driver = &driver{Client: client, Database: database}
	return cfg
}

func (c *config) validate() error {
	if c == nil {
		return errors.New("config is required")
	}

	if c.Driver == nil {
		if c.URI == "" || c.DBName == "" {
			return errors.New("uri/dbname is required")
		}

		return nil
	}

	if c.Driver.Client == nil || c.Driver.Database == nil {
		return errors.New("invalid mongo driver")
	}

	return nil
}

func (c *config) SetAutoPreload(value bool) {
	c.AutoPreload = value
}

func (c *config) SetIDType(t IDType) {
	if t.isValid() {
		c.IDType = t
	}
}

func (c *config) SetClearEmbeddedFields(value bool) {
	c.ClearEmbeddedFields = value
}

func (c *config) SetDriver(client *mongo.Client, database *mongo.Database) {
	if client != nil && database != nil {
		c.Driver = &driver{Client: client, Database: database}
	}
}
