package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Schedule struct {
	Day   string `json:"day"`
	Event string `json:"event"`
	Link  string `json:"link"`
}

func main() {
	fName := "schedule_jkt48.json"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
		return
	}
	defer file.Close()

	schedules := []Schedule{}
	uniqueData := make(map[string]bool)

	// Create a collector
	c := colly.NewCollector()

	// find and visit all links
	c.OnHTML(".entry-schedule__calendar", func(h *colly.HTMLElement) {
		links := h.ChildAttrs("a", "href")
		for i, link := range links {
			jkt48 := "https://jkt48.com"
			full_link := jkt48 + link
			day := h.DOM.Find(".second").Eq(i).Text()
			event := h.DOM.Find(".contents").Eq(i).Text()
			day_trim := strings.TrimSpace(day)
			event_trim := strings.TrimSpace(event)
			fmt.Println("Day:", day, "Event:", event)
			if !uniqueData[event] {
				uniqueData[event] = true
				schedules = append(schedules, Schedule{
					Day:   day_trim,
					Event: event_trim,
					Link:  full_link,
				})
			}

		}
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit("https://jkt48.com/calendar/list?lang=id")

	b, err := json.MarshalIndent(schedules, "", "  ")
	if err != nil {
		log.Fatalf("failed to encode as JSON: %s", err)
		return
	}
	if _, err := file.Write(b); err != nil {
		log.Fatalf("failed writing to file: %s", err)
		return
	}

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
