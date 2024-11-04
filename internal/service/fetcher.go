package service

import (
    "fmt"
    "log"
    "io"
    "net/http"
    "time"
    "encoding/json"
    "youtube-api/config"
    "youtube-api/internal/models"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type YouTubeResponse struct {
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


var currentQuery string

// SetQuery dynamically updates the predefined query
func SetQuery(query string) {
    currentQuery = query
}

// StartFetcher runs an asynchronous fetcher that fetches videos for a given query periodically
func StartFetcher(collection *mongo.Collection, cfg *config.Config) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        if currentQuery != "" {
            _ , err := FetchYouTubeVideos(currentQuery, collection, cfg)
            if err != nil {
                log.Printf("Error fetching videos: %v", err)
                continue
            }
        }
    }
}



func FetchYouTubeVideos(query string, collection *mongo.Collection, cfg *config.Config) ([]models.Video, error) {
    apiKey := cfg.YouTubeAPIKey
    requestURL := fmt.Sprintf(
        "https://www.googleapis.com/youtube/v3/search?part=snippet&type=video&order=date&q=%s&key=%s",
        query, apiKey,
    )

    resp, err := http.Get(requestURL)
    if err != nil || resp.StatusCode != http.StatusOK {
        log.Printf("Error fetching data from YouTube API: %v", err)
        // http.Error(w, "Failed to fetch data from YouTube API", http.StatusInternalServerError)
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        // http.Error(w, "Failed to read response from YouTube API", http.StatusInternalServerError)
        return nil, err
    }

    var ytResponse YouTubeResponse
    if err := json.Unmarshal(body, &ytResponse); err != nil {
        // http.Error(w, "Failed to parse response from YouTube API", http.StatusInternalServerError)
        return nil, err
    }

    // Convert YouTube API response to internal Video structure
    results := []models.Video{}
    for _, item := range ytResponse.Items {
        // Parse and set PublishDate explicitly to UTC
        publishDate, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
        if err != nil {
            log.Printf("Error parsing publish date for video %s: %v", item.Snippet.Title, err)
            continue
        }
        results = append(results, models.Video{
            ID:          item.ID.VideoID,
            Title:       item.Snippet.Title,
            Description: item.Snippet.Description,
            PublishDate: primitive.NewDateTimeFromTime(publishDate),
            ThumbnailURL:   item.Snippet.Thumbnails.Default.URL,
        })
    }

    log.Printf("Fetched %d videos for query '%s'", len(results), query)
    
    for _, video := range results {
        err := models.InsertVideo(collection, video)
        if err != nil {
            log.Printf("Error inserting video %s: %v", video.Title, err)
        } else {
            log.Printf("Inserted video: %s", video.Title)
        }
    }
    return results, nil
}

