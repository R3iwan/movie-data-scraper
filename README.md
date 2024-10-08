# Movie Data Scraper using Golang

This Go application scrapes movie data (title and description) from IMDb using a list of movie URLs. It sends concurrent HTTP requests to fetch and parse the HTML content of the IMDb pages, extracting metadata like the title and description.

## Features
- Fetches movie data (title and description) from IMDb URLs.
- Uses Go's concurrency model (`sync.WaitGroup` and Goroutines) to handle multiple requests efficiently.
- Parses HTML using the `golang.org/x/net/html` package.

## Installation Instructions

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/movie-data-scraper.git
