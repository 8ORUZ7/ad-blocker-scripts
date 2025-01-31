package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// Configurable values
const (
	maxTries = 100
	delay    = 200 * time.Millisecond
)

// runAdBypass starts the Chrome session and executes the script
func runAdBypass(ctx context.Context) {
	var currentPageURL string
	var tries int

	// Capture current URL
	err := chromedp.Run(ctx,
		chromedp.Evaluate(`window.location.href`, &currentPageURL),
	)
	if err != nil {
		log.Println("Failed to capture current URL:", err)
	}

	// Watch for page navigation and refresh iframe
	for tries < maxTries {
		time.Sleep(delay)
		var newURL string

		// Fetch the updated URL
		err := chromedp.Run(ctx,
			chromedp.Evaluate(`window.location.href`, &newURL),
		)
		if err != nil {
			log.Println("Failed to fetch updated URL:", err)
			continue
		}

		// If user leaves the watch page, remove iframe
		if !strings.Contains(newURL, "watch") {
			removeIframe(ctx)
			break
		}

		// If on a video page, inject the iframe
		if videoID := extractVideoID(newURL); videoID != "" {
			createIframe(ctx, videoID)
		}

		tries++
	}
}

// extractVideoID extracts video ID from a YouTube URL
func extractVideoID(videoURL string) string {
	parsedURL, err := url.Parse(videoURL)
	if err != nil {
		log.Println("Failed to parse URL:", err)
		return ""
	}

	params := parsedURL.Query()
	videoID := params.Get("v")
	if videoID == "" {
		log.Println("No video ID found in URL")
	}
	return videoID
}

// createIframe injects an embedded YouTube iframe
func createIframe(ctx context.Context, videoID string) {
	embedURL := fmt.Sprintf(`https://www.youtube-nocookie.com/embed/%s?autoplay=1&modestbranding=1`, videoID)

	err := chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf(`
            var player = document.getElementById("youtube-iframe");
            if (!player) {
                var iframe = document.createElement("iframe");
                iframe.src = "%s";
                iframe.style = "height:100%%; width:100%%; border-radius:12px;";
                iframe.id = "youtube-iframe";
                document.body.appendChild(iframe);
            } else {
                player.src = "%s";
            }
        `, embedURL, embedURL), nil),
	)
	if err != nil {
		log.Println("Failed to create iframe:", err)
	}
}

// removeIframe removes the injected iframe when leaving the video page
func removeIframe(ctx context.Context) {
	err := chromedp.Run(ctx,
		chromedp.Evaluate(`
            var player = document.getElementById("youtube-iframe");
            if (player && player.parentNode) {
                player.parentNode.removeChild(player);
            }
        `, nil),
	)
	if err != nil {
		log.Println("Failed to remove iframe:", err)
	}
}

func main() {
	// Create a new browser context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Run the ad bypass script
	runAdBypass(ctx)
}

//
