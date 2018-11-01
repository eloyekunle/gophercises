package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/eloyekunle/gophercises/13/hn"
)

var (
	cache                           cacheItem
	port, numStories, cacheDuration int
)

func main() {
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.IntVar(&cacheDuration, "cache_duration", 5, "duration in seconds to cache content")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))
	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

var storyMutex sync.Mutex

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		stories, err := getTopStories()

		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func getTopStories() ([]Story, error) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, err
	}
	// We're getting slightly more than 'numStories' to account for filtering.
	hedgedNum := numStories * 5 / 4
	seen := 0
	c := make(chan Story)

	for i := 0; i < hedgedNum; i++ {
		go func(id int) {
			hnItem, _ := client.GetItem(id)
			item := parseHNItem(hnItem)
			if isStoryLink(item) {
				c <- item
			}

			storyMutex.Lock()
			defer storyMutex.Unlock()
			seen++
			if seen == hedgedNum {
				close(c)
			}
		}(ids[i])
	}

	storiesMap := make(map[int]Story, numStories)
	for item := range c {
		storiesMap[item.Item.ID] = item
	}

	var stories []Story
	for i := 0; len(stories) < numStories; i++ {
		item, ok := storiesMap[ids[i]]

		if ok {
			stories = append(stories, item)
		}
	}

	return stories, nil
}

func isStoryLink(item Story) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) Story {
	ret := Story{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// Story is the same as the hn.Item, but adds the Host field
type Story struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []Story
	Time    time.Duration
}

// Represents an item in our in-memory cache.
type cacheItem struct {
	Content    []byte
	Expiration int64
	Mutex      sync.Mutex
}
