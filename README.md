# YouTube API Application

This is a simple Go application that interacts with the YouTube API to fetch and store video data in MongoDB. The application provides a RESTful API to retrieve videos and perform searches based on titles and descriptions.

## Features

- Fetch the latest videos from YouTube using a predefined search query.
- Store video details such as title, description, publish date, and thumbnail URL in MongoDB.
- Retrieve a paginated list of stored videos.
- Search for videos by title and description.

## Technologies Used

- Go (Golang)
- MongoDB (MongoDB Atlas or local instance)
- Gorilla Mux for routing
- Postman for API testing

## Getting Started

### Prerequisites

- Go installed on your machine (1.16+)
- MongoDB instance (Atlas or local)
- Postman (for testing APIs)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/gagan-jasuja/youtube-api.git
   cd your-repo-name
   ```

2. Set up your environment variables in a `.env` file:
   ```bash
   MONGODB_URI=your_mongodb_uri
   ```

3. Install the necessary dependencies:
   ```bash
   go mod tidy
   ```

### Running the Application

1. Start the application:
   ```bash
   go run ./cmd/main.go
   ```

2. The application will run on `http://localhost:8080`.


### API Endpoints

#### 1. Get Paginated Videos
- **URL:** `/videos?page=1&limit=10`
- **Method:** `GET`
- **Description:** Retrieves a paginated list of videos sorted by publish date in descending order.

#### 2. Search Videos
- **URL:** `/search?query=your_search_query`
- **Method:** `GET`
- **Description:** Searches for videos based on title or description.

### Testing with Postman

You can test the API using Postman by following these steps:

1. Open Postman.
2. Create a new request.
3. Set the HTTP method (GET) and enter the API endpoint URL.
4. Click "Send" to make the request and view the response.

### Example Requests

**Fetch Videos:**
```
GET http://localhost:8080/videos?page=1&limit=10
```

**Search Videos:**
```
GET http://localhost:8080/search?query=Movies
```

### License

Primary_Owner - Gagan Jasuja

## Contributing

Contributions are welcome! Please create a pull request or open an issue to discuss changes.

## Acknowledgements

- [YouTube Data API](https://developers.google.com/youtube/v3)
- [MongoDB](https://www.mongodb.com/)
