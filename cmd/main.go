package main

import (
    "context"
    "log"
    "time"
    "youtube-api/config"
    "youtube-api/internal/api"
    "youtube-api/internal/db"
    "youtube-api/internal/service"
)

func main() {
    // Load configuration
    cfg := config.Load()

    // Initialize the database client
    dbClient, err := db.Connect(cfg.MongoDBURI)
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }

    
    // Ensure to disconnect from the database when the main function returns
    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        if err := dbClient.Disconnect(ctx); err != nil {
            log.Fatalf("Error disconnecting from the database: %v", err)
        }
    }()

    // Get the collection for storing videos
    videoCollection := dbClient.Database(cfg.DatabaseName).Collection(cfg.VideoCollectionName)

    // Start the background fetcher
    go service.StartFetcher(videoCollection, cfg)

    // Start the server
    api.StartServer(dbClient, cfg)
}
