package models

import (
	"context"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
    ID           string             `bson:"_id,omitempty" json:"id"`
    Title        string             `bson:"title" json:"title"`
    Description  string             `bson:"description" json:"description"`
    PublishDate  primitive.DateTime `bson:"publish_date" json:"publish_date"`
    ThumbnailURL string             `bson:"thumbnail_url" json:"thumbnail_url"`
}

func InsertVideo(collection *mongo.Collection, video Video) error {
	_, err := collection.InsertOne(context.TODO(), video)
	if err != nil {
		log.Printf("Error inserting video: %v", err)
	} else {
		log.Printf("Inserted video: %s", video.Title)
	}
	return err
}

func GetPaginatedVideos(collection *mongo.Collection, page, limit int) ([]Video, error) {
	skip := (page - 1) * limit
	findOptions := options.Find()
	findOptions.SetSort(map[string]int{"publish_date": -1}).SetLimit(int64(limit)).SetSkip(int64(skip))

	cursor, err := collection.Find(context.TODO(), map[string]interface{}{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var videos []Video
	for cursor.Next(context.TODO()) {
		var video Video
		if err := cursor.Decode(&video); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}