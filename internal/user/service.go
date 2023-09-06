package user

import (
	"context"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "users"
const timeOut = 1 * time.Minute

func FindAllUsers() ([]User, error) {
	log := config.InitializeLogger().LogHandler

	coll := db.GetConnection().Database(config.EnvDbName()).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer cur.Close(ctx)

	results := make([]User, 0)

	for cur.Next(ctx) {
		doc := &User{}
		err := cur.Decode(&doc)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		results = append(results, *doc)
	}

	if err := cur.Err(); err != nil {
		log.Error(err)
		return nil, err
	}

	return results, nil
}

func FindUserByID(id string) (*User, error) {
	coll := db.GetConnection().Database(config.EnvDbName()).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	result := &User{}

	objectId, err1 := primitive.ObjectIDFromHex(id)
	if err1 != nil {
		return nil, err1
	}

	// Check for errors when executing FindOne
	err := coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func SaveOne(data *User) (interface{}, error) {

	coll := db.GetConnection().Database(config.EnvDbName()).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	result, err := coll.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil

}

func DeleteOne(id string) bool {

	coll := db.GetConnection().Database(config.EnvDbName()).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	objectId, err1 := primitive.ObjectIDFromHex(id)
	if err1 != nil {
		return false
	}

	filter := bson.M{"_id": objectId}

	result, err := coll.DeleteOne(ctx, filter)

	if err != nil || result.DeletedCount == 0 {
		return false
	}

	return true
}

func UpdateOneUser(user User) error {

	coll := db.GetConnection().Database(config.EnvDbName()).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	filter := bson.M{"_id": user.ID}
	opt := options.Update().SetUpsert(false)

	log := config.InitializeLogger().LogHandler

	update := bson.M{
		"$set": bson.M{
			"fullName": user.FullName,
			"age":      user.Age,
		},
	}

	log.Debug(user)

	_, err := coll.UpdateOne(ctx, filter, update, opt)

	if err != nil {
		return err
	}

	return nil
}
