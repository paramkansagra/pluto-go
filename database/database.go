package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"pluto-go/models"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// first we need to connect and shit
const databaseName string = "pluto-api"
const userCollectionName string = "users"

var UserCollection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI(getConnectionString())

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatalf("Mongo client connection failed , err -> %v \n", err)
	}

	fmt.Println("Mongo Connection success")

	UserCollection = client.Database(databaseName).Collection(userCollectionName)

	fmt.Println("Connection ready with mongo db")
}

func getConnectionString() string {
	err := godotenv.Load("./.env")

	if err != nil {
		log.Fatalf("error loading .env file -> %v \n", err.Error())
	}

	var mongodb_username = os.Getenv("MONGODB_USERNAME")
	var mongodb_password = os.Getenv("MONGODB_PASSWORD")
	var mongodb_url = os.Getenv("MONGODB_URL")

	return fmt.Sprintf("mongodb+srv://%v:%v@%v", mongodb_username, mongodb_password, mongodb_url)
}

func checkSameEmailExists(email string) error {
	emailFilter := bson.D{{Key: "email", Value: email}}

	var user models.User

	err := UserCollection.FindOne(context.TODO(), emailFilter).Decode(&user)
	if err == nil {
		return errors.New("account with same email found")
	}

	return nil
}

func SignUp(user *models.User) error {
	err := models.CheckRequiredFeilds(user)

	if err != nil {
		return err
	}

	err = checkSameEmailExists(user.Email)
	if err != nil {
		return err
	}

	// email check karenge aage jaa ke
	err = models.HashUserPassword(user)

	if err != nil {
		return err
	}

	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	user.EmailVerfied = false
	inserted, err := UserCollection.InsertOne(context.TODO(), user)

	if err != nil {
		return err
	}

	user.ID = inserted.InsertedID.(primitive.ObjectID)

	return nil
}

func SignIn(user *models.User) error {

	return nil
}
