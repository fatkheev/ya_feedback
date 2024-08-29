package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Review struct {
	Author string `json:"author"`
	Date   string `json:"date"`
	Text   string `json:"text"`
	Rating int    `json:"rating"`
}

const cacheFilePath = "./cache/reviews_cache.json"

func main() {
	http.HandleFunc("/reviews", reviewsHandler)
	http.HandleFunc("/feedback", feedbackHandler) // Новый маршрут для загрузки страницы
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Обновление кеша при запуске и каждые 60 минут
	go func() {
		for {
			updateCache()
			time.Sleep(60 * time.Minute)
		}
	}()

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func feedbackHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/feedback.html")
}

func reviewsHandler(w http.ResponseWriter, r *http.Request) {
	reviews, err := loadReviewsFromCache()
	if err != nil {
		http.Error(w, "Failed to load reviews", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

func updateCache() {
	reviews, err := fetchReviews()
	if err != nil {
		log.Println("Failed to update cache:", err)
		return
	}

	data, err := json.Marshal(reviews)
	if err != nil {
		log.Println("Failed to marshal reviews:", err)
		return
	}

	err = os.WriteFile(cacheFilePath, data, 0644)
	if err != nil {
		log.Println("Failed to write cache file:", err)
	}
}

func loadReviewsFromCache() ([]Review, error) {
	file, err := os.Open(cacheFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var reviews []Review
	err = json.NewDecoder(file).Decode(&reviews)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func fetchReviews() ([]Review, error) {
	l := launcher.New().Headless(true).MustLaunch()
	browser := rod.New().ControlURL(l).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://yandex.ru/maps/org/108710443862/reviews/").MustWaitLoad()
	reviewElements := page.MustElements(".business-review-view")

	var reviews []Review
	for _, reviewElement := range reviewElements {
		author, err := reviewElement.Element(".business-review-view__author-name span")
		authorText := "Author not found"
		if err == nil {
			authorText = author.MustText()
		}

		dateElement, err := reviewElement.Element("meta[itemprop='datePublished']")
		dateText := "Date not found"
		if err == nil {
			dateText = dateElement.MustProperty("content").String()
		}

		textElement, err := reviewElement.Element(".business-review-view__body-text")
		text := "Text not found"
		if err == nil {
			text = textElement.MustText()
		}

		ratingElements := reviewElement.MustElements(".business-rating-badge-view__star._full")
		rating := len(ratingElements)

		reviews = append(reviews, Review{
			Author: authorText,
			Date:   dateText,
			Text:   text,
			Rating: rating,
		})
	}

	return reviews, nil
}