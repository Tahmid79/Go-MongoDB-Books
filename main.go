package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	Title     string
	Author    string
	ISBN      string
	Publisher string
	Copies    int
}

func main() {
	URI := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(URI)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	checkError(err)

	err = client.Ping(context.TODO(), nil)

	checkError(err)

	fmt.Println("Connected to MongoDB")

	booksCollection := client.Database("testdb").Collection("books")

	// Insert One
	book1 := Book{"Animal Farm", "George Orwell", "0451526341", "Signet Classics", 100}

	insertResult, err := booksCollection.InsertOne(context.TODO(), book1)

	checkError(err)

	fmt.Println("Inserted a single document ", insertResult.InsertedID)

	//Insert Many

	book2 := Book{"Super Freakonomics", "Steven D. Levitt", "0062312871", "HARPER COLLINS USA", 100}
	book3 := Book{"The Alchemist", "Paulo Coelho", "0062315005", "HarperOne", 100}
	multipleBooks := []interface{}{book2, book3}

	insertManyResult, err := booksCollection.InsertMany(context.TODO(), multipleBooks)

	checkError(err)

	fmt.Println("Inserted multiple documents ", insertManyResult.InsertedIDs)

	// Update One Document

	filter := bson.D{{"isbn", "0451526341"}}

	update := bson.D{
		{"$inc", bson.D{
			{"copies", 10},
		}},
	}

	updateResult, err := booksCollection.UpdateOne(context.TODO(), filter, update)

	checkError(err)

	fmt.Printf("Matched %v documents and updated %v documents\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	//Find Documents

	//Find One document

	// A variable in which the result will be decoded
	var result Book

	err = booksCollection.FindOne(context.TODO(), filter).Decode(&result)

	checkError(err)

	fmt.Printf("Found a single document %+v\n", result)

	// Find multiple documents

	cursor, err := booksCollection.Find(context.TODO(), bson.D{{}})
	checkError(err)

	var books []Book

	err = cursor.All(context.TODO(), &books)
	checkError(err)

	fmt.Printf("Found multiple documents: %+v", books)

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
