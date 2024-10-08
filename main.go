package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Review struct {
	Author       string `json:"author"`
	Date         string `json:"date"`
	Text         string `json:"text"`
	Rating       int    `json:"rating"`
	ProfileImage string `json:"profile_image"`
}

const cacheFilePath = "./cache/reviews_cache.json"

func main() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Println("Failed to get working directory:", err)
	} else {
		log.Println("Current working directory:", workingDir)
	}

	http.HandleFunc("/reviews", reviewsHandler)
	http.HandleFunc("/feedback", feedbackHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

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

	err = os.MkdirAll("./cache", 0755)
	if err != nil {
		log.Println("Failed to create cache directory:", err)
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
	// Определяем путь к браузеру
	browserPath := getBrowserPath()

	l := launcher.New().Headless(true).Bin(browserPath).MustLaunch()
	browser := rod.New().ControlURL(l).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://yandex.ru/maps/org/108710443862/reviews/").MustWaitLoad()

	for {
		previousHeight := page.MustEval(`() => document.body.scrollHeight`).Int()
		page.MustEval(`() => window.scrollTo(0, document.body.scrollHeight)`)
		time.Sleep(2 * time.Second)
		currentHeight := page.MustEval(`() => document.body.scrollHeight`).Int()

		if currentHeight == previousHeight {
			break
		}
	}

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

		profileImageElement, err := reviewElement.Element(".business-review-view__user-icon .user-icon-view__icon")
		profileImageURL := ""
		if err == nil {
			styleAttribute, _ := profileImageElement.Attribute("style")
			if styleAttribute != nil {
				profileImageURL = extractURLFromStyle(*styleAttribute)
			}
		}

		reviews = append(reviews, Review{
			Author:       authorText,
			Date:         dateText,
			Text:         text,
			Rating:       rating,
			ProfileImage: profileImageURL,
		})
	}

	return reviews, nil
}

func extractURLFromStyle(style string) string {
	start := strings.Index(style, `url("`) + 5
	end := strings.Index(style, `")`)
	if start > 4 && end > start {
		return style[start:end]
	}
	return ""
}

func getBrowserPath() string {
	// Пытаемся найти браузер в системе
	possiblePaths := []string{
		"/usr/bin/chromium-browser", // Default for many Linux distributions
		"/usr/bin/google-chrome",    // Default for Google Chrome on Linux
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", // Default for macOS
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Если браузер не найден, пробуем встроенный
	return launcher.NewBrowser().MustGet()
}
