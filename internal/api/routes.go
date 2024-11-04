package api

import (
    "net/http"
    "youtube-api/config"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/gorilla/mux"
)

// StartServer starts the HTTP server
func StartServer(dbConn *mongo.Client, cfg *config.Config) {
    r := mux.NewRouter()

    // Initialize video collection
    videoCollection := dbConn.Database(cfg.DatabaseName).Collection(cfg.VideoCollectionName)
    
    // Initialize the handlers
    handlers := NewHandlers(videoCollection)

    r.HandleFunc("/videos", handlers.GetVideosHandler).Methods("GET")
    r.HandleFunc("/search", handlers.SearchVideosHandler).Methods("GET")

    http.ListenAndServe(":8080", r)
}

