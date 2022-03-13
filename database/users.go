package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (db *DB) CreateUser(input *model.NewUser) (*model.User, error) {
	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": input.Email}
	val := model.User{}
	err := collection.FindOne(ctx, filter).Decode(&val)
	if err == mongo.ErrNoDocuments {
		if len(input.Password) > 0 {
			hash, _ := HashPassword(input.Password)
			input.Password = hash
		} else {
			hash, _ := HashPassword(RandStringRunes(6))
			input.Password = hash
		}

		res, err := collection.InsertOne(ctx, input)
		if err != nil {
			log.Fatal(err)
		}
		return &model.User{
			ID:    res.InsertedID.(primitive.ObjectID).Hex(),
			Email: input.Email,
			Role:  input.Role,
		}, nil
	} else {
		return nil, errors.New("existing email")
	}
}

func (db *DB) GetUser(id string) *model.User {
	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ObjectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ObjectID}
	res := collection.FindOne(ctx, filter)

	user := model.User{}
	res.Decode(&user)
	return &user
}

func (db *DB) UpdateUser(input *model.UpdateUser) *model.User {
	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ObjectID, _ := primitive.ObjectIDFromHex(input.ID)
	filter := bson.M{"_id": ObjectID}
	update := bson.D{{"$set",
		bson.D{
			{"Role", input.Role},
			{"Email", input.Email},
		},
	}}
	_, err := collection.UpdateOne(ctx, filter, update)
	user := model.User{}

	jsonbody, err := json.Marshal(*input)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(jsonbody, &user); err != nil {
		fmt.Println(err)
	}
	user.ID = input.ID

	return &user
}

func (db *DB) DeleteUser(id string) *mongo.DeleteResult {
	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ObjectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ObjectID}
	res, _ := collection.DeleteOne(ctx, filter)
	return res
}
