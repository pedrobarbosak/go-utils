package mongo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNoResults = errors.New("no documents in result")

type repository struct {
	client   *mongo.Client
	database *mongo.Database
	config   *config
}

func NewRepository(cfg *config) (Repository, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	if cfg.Driver == nil {
		if err := connectMongo(cfg); err != nil {
			return nil, err
		}
	}

	return &repository{client: cfg.Driver.Client, database: cfg.Driver.Database, config: cfg}, nil
}

func connectMongo(cfg *config) error {
	clientOptions := options.Client()
	clientOptions.ApplyURI(cfg.URI)
	clientOptions.SetMaxConnIdleTime(5 * time.Second)
	clientOptions.SetSocketTimeout(30 * time.Second)
	clientOptions.SetServerSelectionTimeout(15 * time.Second)

	if cfg.CertificatePath != "" {
		rootCerts := x509.NewCertPool()
		if ca, err := os.ReadFile(cfg.CertificatePath); err == nil {
			rootCerts.AppendCertsFromPEM(ca)
		}
		clientOptions.SetTLSConfig(&tls.Config{RootCAs: rootCerts, ClientCAs: rootCerts, InsecureSkipVerify: true})
	}

	ctx := context.Background()
	cl, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	if err = cl.Ping(ctx, nil); err != nil {
		return err
	}

	cfg.Driver = &driver{Client: cl, Database: cl.Database(cfg.DBName)}
	return nil
}

func (repo *repository) Create(ctx context.Context, object StorableObject) error {
	if repo.config.ClearEmbeddedFields {
		if err := repo.clear(object); err != nil {
			return err
		}
	}

	id, err := repo.database.Collection(object.GetCollection()).InsertOne(ctx, object)
	if err != nil {
		return err
	}

	object.SetID(repo.getInsertedID(id))
	return nil
}

func (repo *repository) GetByID(ctx context.Context, objectID string, object StorableObject) error {
	filter, err := repo.getIDFilter(objectID)
	if err != nil {
		return err
	}

	result := repo.database.Collection(object.GetCollection()).FindOne(ctx, filter)
	if err = result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoResults
		}
		return err
	}

	if err = result.Decode(object); err != nil {
		return err
	}

	if repo.config.AutoPreload {
		return repo.Preload(ctx, object)
	}

	return nil
}

func (repo *repository) GetBy(ctx context.Context, object StorableObject, filters ...Filter) error {
	filter := bson.M{}
	for _, property := range filters {
		filter[property.Key] = property.Value
	}

	result := repo.database.Collection(object.GetCollection()).FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoResults
		}
		return err
	}

	return result.Decode(object)
}

func (repo *repository) Fetch(ctx context.Context, object StorableObject, out interface{}, filters ...Filter) error {
	filter := bson.M{}
	for _, property := range filters {
		filter[property.Key] = property.Value
	}

	cursor, err := repo.database.Collection(object.GetCollection()).Find(ctx, filter)
	if err != nil {
		return err
	}

	outs := make([]Object, 0)
	for cursor.Next(ctx) {
		data := reflect.New(reflect.TypeOf(object).Elem()).Interface()
		if err = cursor.Decode(data); err != nil {
			return err
		}

		obj := data.(StorableObject)
		if repo.config.AutoPreload {
			if err = repo.Preload(ctx, obj); err != nil {
				return err
			}
		}

		outs = append(outs, obj)
	}

	if err = cursor.Err(); err != nil {
		return err
	}

	if err = cursor.Close(ctx); err != nil {
		return err
	}

	b, err := json.Marshal(outs)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, out)
}

func (repo *repository) Update(ctx context.Context, objectID string, object StorableObject) error {
	if repo.config.ClearEmbeddedFields {
		if err := repo.clear(object); err != nil {
			return err
		}
	}

	filter, err := repo.getIDFilter(objectID)
	if err != nil {
		return err
	}

	return repo.database.Collection(object.GetCollection()).FindOneAndUpdate(ctx, filter, bson.D{{Key: "$set", Value: object}}).Err()
}

func (repo *repository) WithTransaction(ctx context.Context, fn func(sc context.Context) error) error {
	session, err := repo.client.StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		return err
	}

	return mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err = fn(sc); err != nil {
			_ = session.AbortTransaction(sc)
			return err
		}

		return session.CommitTransaction(sc)
	})
}

func (repo *repository) Aggregate(ctx context.Context, object StorableObject, query string, out interface{}) error {
	opts := options.Aggregate()
	opts.SetCollation(&options.Collation{Locale: "en", Strength: 3})
	opts.SetAllowDiskUse(true)

	var pipeline interface{}
	if err := bson.UnmarshalExtJSON([]byte(query), true, &pipeline); err != nil {
		return err
	}

	cursor, err := repo.database.Collection(object.GetCollection()).Aggregate(ctx, pipeline, opts)
	if err != nil {
		return err
	}

	var objects []interface{}
	for cursor.Next(ctx) {
		var obj interface{}
		if err = cursor.Decode(&obj); err != nil {
			return err
		}

		objects = append(objects, obj)
	}

	if err = cursor.Err(); err != nil {
		return err
	}

	if err = cursor.Close(ctx); err != nil {
		return err
	}

	if len(objects) == 0 {
		return nil
	}

	b, err := bson.Marshal(objects[0])
	if err != nil {
		return err
	}

	if err = bson.Unmarshal(b, out); err != nil {
		return err
	}

	if repo.config.AutoPreload {
		return repo.search(ctx, out)
	}

	return nil
}

func (repo *repository) Count(ctx context.Context, object StorableObject, filter interface{}) (int64, error) {
	return repo.database.Collection(object.GetCollection()).CountDocuments(ctx, filter)
}

func (repo *repository) UpdateOne(ctx context.Context, object StorableObject, filter interface{}, update interface{}) (int64, error) {
	result, err := repo.database.Collection(object.GetCollection()).UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return result.MatchedCount, nil
}

func (repo *repository) DeleteAll(ctx context.Context, object StorableObject) error {
	_, err := repo.database.Collection(object.GetCollection()).DeleteMany(ctx, bson.D{})
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) Disconnect(ctx context.Context) error {
	return repo.client.Disconnect(ctx)
}

func (repo *repository) CreateMany(ctx context.Context, obj StorableObject, data []interface{}) error {
	if len(data) == 0 {
		return nil
	}

	_, err := repo.database.Collection(obj.GetCollection()).InsertMany(ctx, data)
	return err
}
