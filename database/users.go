package database

import (
	"context"
	"errors"
	"log"
	"time"

	"ksemilla/graph/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DB) AllUsers() []*model.User {
	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var users []*model.User
	for cur.Next(ctx) {
		var user *model.User
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	return users
}

func (db *DB) FindOneUser(_id string) (*model.User, error) {
	ObjectID, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		log.Fatal(err)
	}
	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user := model.User{}
	res := collection.FindOne(ctx, bson.M{"_id": ObjectID})
	err = res.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no user found")
		}
	}
	return &user, nil
}
