package main

import (
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/net/html"
)

type Data struct {
	Title       string `json:"titleText"`
	Description string `json:"description"`
}

func findInfo(n *html.Node, data *Data, foundTitle *bool) {
	if n.Type == html.ElementNode && n.Data == "meta" {
		var hasOgTitle bool
		var content string
		for _, attr := range n.Attr {
			if attr.Key == "property" && attr.Val == "og:title" {
				hasOgTitle = true
			}
			if attr.Key == "content" {
				content = attr.Val
			}
		}
		if hasOgTitle {
			data.Title = content
			*foundTitle = true
		}
	}

	if n.Type == html.ElementNode && n.Data == "span" {
		var hasDescriptionClass, hasDescriptionId bool
		for _, attr := range n.Attr {
			if attr.Key == "class" && attr.Val == "sc-55855a9b-0 dAbouZ" {
				hasDescriptionClass = true
			}
			if attr.Key == "data-testid" && attr.Val == "plot-xs_to_m" {
				hasDescriptionId = true
			}
		}
		if hasDescriptionClass && hasDescriptionId && n.FirstChild != nil {
			data.Description = n.FirstChild.Data
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findInfo(c, data, foundTitle)
	}
}

func GetMovieData(url string, data chan<- Data, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error fetching:", err)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	var foundTitle bool
	var movieData Data
	findInfo(doc, &movieData, &foundTitle)
	data <- movieData
}

func main() {
	var wg sync.WaitGroup

	urls := []string{
		"https://www.imdb.com/title/tt0111161/?ref_=chttp_t_1",
		"https://www.imdb.com/title/tt0068646/?ref_=chttp_tt_2",
		"https://www.imdb.com/title/tt0071562/?ref_=chttp_tt_3",
		"https://www.imdb.com/title/tt0468569/?ref_=chttp_tt_4",
		"https://www.imdb.com/title/tt0050083/?ref_=chttp_tt_6",
	}

	data := make(chan Data)

	go func() {
		wg.Wait()
		close(data)
	}()

	for _, url := range urls {
		wg.Add(1)
		go GetMovieData(url, data, &wg)
	}

	for movieData := range data {
		fmt.Printf("Title: %s\nDescription: %s\n\n", movieData.Title, movieData.Description)
	}
}
