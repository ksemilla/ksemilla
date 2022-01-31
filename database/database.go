package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	"ksemilla/config"
	"ksemilla/graph/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"path/filepath"

	"github.com/joho/godotenv"
)

type DB struct {
	client *mongo.Client
}

func Connect() *DB {

	godotenv.Load(filepath.Join(".", ".env"))
	config := config.Config()

	client, err := mongo.NewClient(options.Client().ApplyURI(config.MONGODB_URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{
		client: client,
	}
}

func (db *DB) Save(input *model.NewInvoice) *model.Invoice {
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
	fmt.Println("asd", res)
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
	var perPage int64 = 3
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

func (db *DB) PaginatedInvoice(dateCreated string) []*model.Invoice {
	collection := db.client.Database("ksemilla").Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var invoices []*model.Invoice
	fmt.Println("SEARCHING FOR", dateCreated, reflect.TypeOf(dateCreated))
	// test, _ := primitive.ObjectIDFromHex("61e951588ca65b48275ca0e2")
	// filter := bson.M{"_id": test}
	filter := bson.M{"datecreated": "2020-01-01"}
	// filter = bson.M{}
	// filter := bson.D{{"created", bson.D{{"$eq", "2020-01-01"}}}}
	// filter := bson.M{"datecreated": bson.M{"$regex": primitive.Regex{Pattern: "2", Options: "i"}}}
	// filter = bson.D{}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (db *DB) CreateUser(input *model.NewUser) *model.User {
	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hash, _ := HashPassword(input.Password)
	input.Password = hash

	res, err := collection.InsertOne(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	return &model.User{
		ID:    res.InsertedID.(primitive.ObjectID).Hex(),
		Email: input.Email,
	}
}

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

func (db *DB) Login(input *model.Login) (*model.LoginReturn, error) {

	config := config.Config()

	collection := db.client.Database("ksemilla").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result *bson.M
	err := collection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&result)
	if err != nil {
		return nil, errors.New("cannot find user")
	}
	user_map := *result
	password := user_map["password"].(string)
	match := CheckPasswordHash(input.Password, password)
	_token_duration, _ := strconv.Atoi(os.Getenv("TOKEN_DURATION"))
	if match {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"foo":       "bar",
			"key":       "value",
			"ExpiresAt": int(time.Now().Add(time.Second * time.Duration(_token_duration)).Unix()),
			"userId":    user_map["_id"].(primitive.ObjectID).Hex(),
		})
		tokenString, _ := token.SignedString([]byte(config.APP_SECRET_KEY))

		jsonbody, _ := json.Marshal(user_map)
		user := model.User{}

		err = json.Unmarshal(jsonbody, &user)
		if err != nil {
			return nil, errors.New("json unmarshal error")
		}

		// return &model.LoginReturn{
		// 	User:  &model.User{},
		// 	Token: tokenString,
		// }, nil

		return &model.LoginReturn{
			User:  &user,
			Token: tokenString,
		}, nil
	} else {
		return nil, errors.New("wrong credentials")
	}
}

func (db *DB) VerifyToken(input *model.VerifyToken) (*model.User, error) {

	config := config.Config()

	token, err := jwt.Parse(input.Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New("unexpected signing method")
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(config.APP_SECRET_KEY), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		timeValue := int64(claims["ExpiresAt"].(float64)) - time.Now().Unix()
		if timeValue <= 0 {
			return nil, errors.New("expired token")
		}

		ObjectID, err := primitive.ObjectIDFromHex(claims["userId"].(string))
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
	} else {
		fmt.Println(err)
		return nil, errors.New("token unrecognized")
	}
}
