package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/robertprincipe/graphql-server/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DATABASE   = "graphql"
	COLLECTION = "videos"
)

type VideoRepository interface {
	Save(video *model.Video)
	FindAll() []*model.Video
}

type database struct {
	client *mongo.Client
}

func New() VideoRepository {

	// MONGODB := os.Getenv("MONGODB")

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/graphql")

	clientOptions = clientOptions.SetMaxPoolSize(50)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)

	dbClient, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB database is connected!")

	return &database{
		client: dbClient,
	}
}

func (d *database) Save(video *model.Video) {
	collection := d.client.Database(DATABASE).Collection(COLLECTION)
	_, err := collection.InsertOne(context.TODO(), video)

	if err != nil {
		log.Fatal(err)
	}
}

func (d *database) FindAll() []*model.Video {
	collection := d.client.Database(DATABASE).Collection(COLLECTION)
	cursor, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	var videos []*model.Video

	for cursor.Next(context.TODO()) {
		var v *model.Video
		err := cursor.Decode(&v)
		if err != nil {
			log.Fatal(err)
		}
		videos = append(videos, v)
	}

	return videos
}
