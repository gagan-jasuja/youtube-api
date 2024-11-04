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
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var ytResponse YouTubeResponse
    if err := json.Unmarshal(body, &ytResponse); err != nil {
        return nil, err
    }

    results := []models.Video{}
    for _, item := range ytResponse.Items {
        results = append(results, models.Video{
            ID:          item.ID.VideoID,
            Title:       item.Snippet.Title,
            Description: item.Snippet.Description,
            PublishDate: primitive.NewDateTimeFromTime(time.Now()),
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

