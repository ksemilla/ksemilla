package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"ksemilla/graph/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) CreateInvoice(input *model.NewInvoice) *model.Invoice {
	collection := db.client.Database("ksemilla").Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, input)
	if err != nil {
		log.Fatal(err)
		// return &model.Invoice{}
	}
	return &model.Invoice{
		ID:          res.InsertedID.(primitive.ObjectID).Hex(),
		DateCreated: input.DateCreated,
		From:        input.From,
		Address:     input.Address,
		Amount:      input.Amount,
	}
}

func (db *DB) All() []*model.Invoice {
	collection := db.client.Database("ksemilla").Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "datecreated", Value: -1}})
	var page int64 = 1
	var perPage int64 = 3

	findOptions.SetSkip((page - 1) * perPage)
	findOptions.SetLimit(perPage)

	cur, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	var invoices []*model.Invoice
	for cur.Next(ctx) {
		var invoice *model.Invoice
		err := cur.Decode(&invoice)
		if err != nil {
			log.Fatal(err)
		}
		invoices = append(invoices, invoice)
	}
	return invoices
}

func (db *DB) GetInvoice(id string) *model.Invoice {
	collection := db.client.Database("ksemilla").Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ObjectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ObjectID}
	res := collection.FindOne(ctx, filter)
	invoice := model.Invoice{}
	res.Decode(&invoice)
	return &invoice
}

func (db *DB) InvoicesPaginated(page int64) ([]*model.Invoice, int64) {
	collection := db.client.Database("ksemilla").Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "datecreated", Value: -1}})
	var perPage int64 = 5
	total, _ := collection.CountDocuments(ctx, bson.M{})

	findOptions.SetSkip((page - 1) * perPage)
	findOptions.SetLimit(perPage)

	cur, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	var invoices []*model.Invoice
	for cur.Next(ctx) {
		var invoice *model.Invoice
		err := cur.Decode(&invoice)
		if err != nil {
			log.Fatal(err)
		}
		invoices = append(invoices, invoice)
	}
	return invoices, total
}

func (db *DB) UpdateInvoice(input *model.InvoiceInput) *model.Invoice {
	collection := db.client.Database("ksemilla").Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ObjectID, _ := primitive.ObjectIDFromHex(input.ID)
	filter := bson.M{"_id": ObjectID}
	update := bson.D{{"$set",
		bson.D{
			{"Address", input.Address},
			{"From", input.From},
			{"DateCreated", input.DateCreated},
			{"Amount", input.Amount},
		},
	}}
	res, err := collection.UpdateOne(ctx, filter, update)
	fmt.Println(res, err, reflect.TypeOf(res))
	invoice := model.Invoice{}
	// res.Decode(&invoice)

	jsonbody, err := json.Marshal(*input)
	fmt.Println(jsonbody)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(jsonbody, &invoice); err != nil {
		// do error check
		fmt.Println(err)
	}
	invoice.ID = input.ID

	return &invoice
}

func (db *DB) DeleteInvoice(id string) *mongo.DeleteResult {
	collection := db.client.Database("ksemilla").Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ObjectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ObjectID}
	res, _ := collection.DeleteOne(ctx, filter)
	return res
}
