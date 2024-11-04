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

    cfg := config.Load()
    dbClient, err := db.Connect(cfg.MongoDBURI)
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }

    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        if err := dbClient.Disconnect(ctx); err != nil {
            log.Fatalf("Error disconnecting from the database: %v", err)
        }
    }()

    videoCollection := dbClient.Database(cfg.DatabaseName).Collection(cfg.VideoCollectionName)

    go service.StartFetcher(videoCollection, cfg)
    
    api.StartServer(dbClient, cfg)
}
