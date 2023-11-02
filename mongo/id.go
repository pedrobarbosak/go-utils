package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IDType uint

const (
	String IDType = iota
	ObjectID
)

func (id IDType) isValid() bool {
	return id == String || id == ObjectID
}

func (repo *repository) getIDFilter(id string) (bson.D, error) {
	if repo.config.IDType == String {
		return bson.D{{Key: "_id", Value: id}}, nil
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return bson.D{{Key: "_id", Value: objectID}}, nil
}

func (repo *repository) getInsertedID(result *mongo.InsertOneResult) string {
	if repo.config.IDType == String {
		return result.InsertedID.(string)
	}

	return result.InsertedID.(primitive.ObjectID).Hex()
}
