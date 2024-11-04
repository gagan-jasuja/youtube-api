package api

import (
	// "context"
	"encoding/json"
	"fmt"

	// "log"
	"net/http"
	// "os"
    // "io"
    "strings"
    "strconv"
	// "time"
    "youtube-api/internal/models"
    "youtube-api/config"
    "youtube-api/internal/service"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type Handlers struct {
	Collection *mongo.Collection
}

// NewHandlers initializes handlers with MongoDB collection
func NewHandlers(collection *mongo.Collection) *Handlers {
	return &Handlers{Collection: collection}
}

// YouTube API URL
const youtubeAPIURL = "https://www.googleapis.com/youtube/v3/search"


// GetVideosHandler returns a paginated list of videos, sorted by publish date in descending order
func (h *Handlers) GetVideosHandler(w http.ResponseWriter, r *http.Request) {
    pageStr := r.URL.Query().Get("page")
    limitStr := r.URL.Query().Get("limit")

    // Default pagination values if not provided
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }

    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 {
        limit = 10
    }

    // Fetch paginated videos from MongoDB
    videos, err := models.GetPaginatedVideos(h.Collection, page, limit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return videos in JSON format
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(videos)
}

// YouTubeAPIResponse represents the structure of the YouTube API response
type YouTubeAPIResponse struct {
    Items []struct {
        ID struct {
            VideoID string `json:"videoId"`
        } `json:"id"`
        Snippet struct {
            Title       string `json:"title"`
            Description string `json:"description"`
            PublishedAt string `json:"publishedAt"`
            Thumbnails  struct {
                Default struct {
                    URL string `json:"url"`
                } `json:"default"`
            } `json:"thumbnails"`
        } `json:"snippet"`
    } `json:"items"`
}

// SearchVideosHandler searches for videos from YouTube API based on a query
func (h *Handlers) SearchVideosHandler(w http.ResponseWriter, r *http.Request) {
        query := strings.TrimSpace(r.URL.Query().Get("query"))
    if query == "" {
        http.Error(w, "Query parameter is required", http.StatusBadRequest)
        return
    }
    str := ""
	for i := 0; i < len(query); i++ {
		if query[i] == ' ' {
			str = str + "+"
		} else {
			str = str + string(query[i])
		}
	}
	fmt.Println(str)
    cfg := config.Load()
    results, err := service.FetchYouTubeVideos(query, h.Collection, cfg)
    if err!=nil {
        fmt.Printf("Error fetching data from YouTube API: %v", err)
        http.Error(w, "Failed to fetch data from YouTube API", http.StatusInternalServerError)
        return
    }

    // Send response to client
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}
