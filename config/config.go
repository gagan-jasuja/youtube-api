package config

import (
    "os"
    "time"
    "log"
    "github.com/joho/godotenv"
)

type Config struct {
    YouTubeAPIKey string
    MongoDBURI    string
    DatabaseName  string
    VideoCollectionName  string
    FetchInterval time.Duration
}

func Load() *Config {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file")
    }
    return &Config{
        YouTubeAPIKey: os.Getenv("YOUTUBE_API_KEY"),
        MongoDBURI:    os.Getenv("MONGODB_URI"),
        DatabaseName:       os.Getenv("DATABASE_NAME"),
        VideoCollectionName: os.Getenv("VIDEO_COLLECTION_NAME"),
        FetchInterval: 10 * time.Second,
    }
}
