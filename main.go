package main

import (
	"fmt"
	// "log"
	// "strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	// Запуск браузера с Rod
	l := launcher.New().Headless(true).MustLaunch()
	browser := rod.New().ControlURL(l).MustConnect()
	defer browser.MustClose()

	// Открытие страницы с отзывами
	page := browser.MustPage("https://yandex.ru/maps/org/108710443862/reviews/").MustWaitLoad()

	// Извлечение списка отзывов
	reviews := page.MustElements(".business-review-view")

	for _, review := range reviews {
		// Извлечение имени автора
		author, err := review.Element(".business-review-view__author-name span")
		if err != nil {
			fmt.Println("Author not found")
		} else {
			fmt.Println("Author:", author.MustText())
		}

		// Извлечение даты
		date, err := review.Element("meta[itemprop='datePublished']")
		if err != nil {
			fmt.Println("Date not found")
		} else {
			fmt.Println("Date:", date.MustProperty("content"))
		}

		// Извлечение текста отзыва
		text, err := review.Element(".business-review-view__body-text")
		if err != nil {
			fmt.Println("Text not found")
		} else {
			fmt.Println("Text:", text.MustText())
		}

		// Извлечение рейтинга
		ratingElements := review.MustElements(".business-rating-badge-view__star._full")
		rating := len(ratingElements)
		if rating == 0 {
			fmt.Println("Rating attribute not found")
		} else {
			fmt.Printf("Rating: %d/5\n", rating)
		}

		fmt.Println("---")
	}
}