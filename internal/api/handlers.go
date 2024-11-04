package api

import (
	"encoding/json"
	"fmt"
	"net/http"
    "strings"
    "strconv"
    "youtube-api/internal/models"
    "youtube-api/config"
    "youtube-api/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handlers struct {
	Collection *mongo.Collection
}

func NewHandlers(collection *mongo.Collection) *Handlers {
	return &Handlers{Collection: collection}
}

func (h *Handlers) GetVideosHandler(w http.ResponseWriter, r *http.Request) {
    pageStr := r.URL.Query().Get("page")
    limitStr := r.URL.Query().Get("limit")

    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }

    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 {
        limit = 10
    }

    videos, err := models.GetPaginatedVideos(h.Collection, page, limit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(videos)
}

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

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}
