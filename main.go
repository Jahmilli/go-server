package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type album struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title  string             `json:"title,omitempty" bson:"title,omitempty"`
	Artist string             `json:"artist,omitempty" bson:"artist,omitempty"`
	Price  float64            `json:"price,omitempty" bson:"price,omitempty"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: primitive.NewObjectID(), Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: primitive.NewObjectID(), Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: primitive.NewObjectID(), Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

var client *mongo.Client
var albumCollection *mongo.Collection

func main() {
	// Setup Database Connection
	authCredentials := options.Credential{
		// AuthMechanism: "PLAIN",
		Username: "admin",
		Password: "pass",
	}

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(authCredentials))
	if err != nil {
		panic(err)
	}

	// Ping the primary instance of Mongo to check liveness
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	albumCollection = client.Database("testing").Collection("albums")

	// HELP: Is it okay to just globally access the repositories like this? Should we be passing in repo via context, building a class or anything?
	// insertOne()
	// insertMany()
	// readOne()
	// readMany()

	// // Setup Router
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/album/:id", getAlbum)
	router.POST("/album", postAlbum)
	router.Run("localhost:8080")
}

func getAlbums(c *gin.Context) {

	var albums []album
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	cursor, err := albumCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
	}

	c.JSON(http.StatusOK, albums)
}

func getAlbum(c *gin.Context) {
	id := c.Param("id")
	// ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	// albumCollection.FindOne(ctx, {})
	for _, a := range albums {
		if a.ID.String() == id {
			c.JSON(http.StatusOK, a)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func postAlbum(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		fmt.Println("An error occurred", err)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := albumCollection.InsertOne(ctx, newAlbum)
	fmt.Println("result is ", result)
	// albums = append(albums, newAlbum)
	c.JSON(http.StatusOK, result)
}

func insertOne() {
	album := bson.D{{"id", "1"}, {"title", "Blue Train"}, {"artist", "John Coltrace"}, {"price", 50.00}}

	result, err := albumCollection.InsertOne(context.TODO(), album)

	if err != nil {
		panic(err)
	}

	fmt.Println(result.InsertedID)
}

func insertMany() {
	albums := []interface{}{
		bson.D{{"id", "2"}, {"title", "Jeru"}, {"artist", "Gerry Mulligan"}, {"price", 17.99}},
		bson.D{{"id", "3"}, {"title", "Sarah Vaughan and Clifford Brown"}, {"artist", "Sarah Vaughan"}, {"price", 39.99}}}

	res, error := albumCollection.InsertMany(context.TODO(), albums)

	if error != nil {
		panic(error)
	}
	for _, val := range res.InsertedIDs {
		fmt.Println("val is ", val)
	}
}

func readMany() {
	cursor, err := albumCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	// Convert result to bson
	var results []bson.M
	// check for errors in conversion
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	fmt.Println("Displaying all results in a collection")
	for _, result := range results {
		fmt.Println(result)
	}
}

func readOne() {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"title", bson.D{{"$eq", "Jeru"}}},
				},
			},
		},
	}
	// retrieve all the documents that match the filter
	cursor, err := albumCollection.Find(context.TODO(), filter)
	// check for errors in the finding
	if err != nil {
		panic(err)
	}

	// convert the cursor result to bson
	var results []bson.M
	// check for errors in the conversion
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	// display the documents retrieved
	fmt.Println("displaying all results from the search query")
	for _, result := range results {
		fmt.Println(result)
	}

	// retrieving the first document that match the filter
	var result bson.M
	// check for errors in the finding
	if err = albumCollection.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		panic(err)
	}

	// display the document retrieved
	fmt.Println("displaying the first result from the search filter")
	fmt.Println(result)
}

func deleteOne() {
	// delete single and multiple documents with a specified filter using DeleteOne() and DeleteMany()
	// create a search filer
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"age", bson.D{{"$gt", 25}}},
				},
			},
		},
	}

	// delete the first document that match the filter
	result, err := albumCollection.DeleteOne(context.TODO(), filter)
	// check for errors in the deleting
	if err != nil {
		panic(err)
	}
	// display the number of documents deleted
	fmt.Println("deleting the first result from the search filter")
	fmt.Println("Number of documents deleted:", result.DeletedCount)

}

func deleteMany() {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"age", bson.D{{"$gt", 25}}},
				},
			},
		},
	}
	// delete every document that match the filter
	results, err := albumCollection.DeleteMany(context.TODO(), filter)
	// check for errors in the deleting
	if err != nil {
		panic(err)
	}
	// display the number of documents deleted
	fmt.Println("deleting every result from the search filter")
	fmt.Println("Number of documents deleted:", results.DeletedCount)
}

/* ------------------- Context stuff here ------------------- */

// func main() {
// 	ctx := context.Background()

// 	ctx, cancel := context.WithCancel(ctx)

// 	time.AfterFunc(time.Second*2, cancel)
// 	sleepAndTalk(ctx, 1*time.Second, "hello")
// }

// func sleepAndTalk(context context.Context, waitTime time.Duration, message string) {

// 	select {
// 	case <-time.After(waitTime):
// 		fmt.Println(message)
// 	case <-context.Done():
// 		log.Print(context.Err())
// 	}
// 	// time.After(waitTime)
// 	// fmt.Println(message)
// }
