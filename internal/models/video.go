package models

import (
	"context"
	"log"
    "time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Video represents a video document in MongoDB
type Video struct {
	ID           string             `json:"id" bson:"_id,omitempty"`
	Title        string             `json:"title" bson:"title"`
	Description  string             `json:"description" bson:"description"`
	PublishDate  primitive.DateTime `json:"publish_date" bson:"publish_date"`
	ThumbnailURL string             `json:"thumbnail_url" bson:"thumbnail_url"`
}

// InsertVideo inserts a new video into the collection
func InsertVideo(collection *mongo.Collection, video Video) error {
    if publishDate, ok := interface{}(video.PublishDate).(time.Time); ok {
		video.PublishDate = primitive.NewDateTimeFromTime(publishDate)
	}
	_, err := collection.InsertOne(context.TODO(), video)
	if err != nil {
		log.Printf("Error inserting video: %v", err)
	} else {
		log.Printf("Inserted video: %s", video.Title)
	}
	return err
}

// GetPaginatedVideos fetches paginated video results from the database
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

// SearchVideos searches for videos by title and description using regex in JSON format
func SearchVideos(collection *mongo.Collection, query string) ([]Video, error) {
	filter := map[string]interface{}{
		"$or": []interface{}{
			map[string]interface{}{"title": map[string]interface{}{"$regex": query, "$options": "i"}},
			map[string]interface{}{"description": map[string]interface{}{"$regex": query, "$options": "i"}},
		},
	}
	log.Println("Reached Search Videos *****************")
	cursor, err := collection.Find(context.TODO(), filter)
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
