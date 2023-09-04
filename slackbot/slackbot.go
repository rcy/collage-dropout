package slackbot

import (
	"log"
	"os"
	"regexp"

	"github.com/slack-go/slack"
)

func Fetch() {
	// Set up a Slack API client
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	// Get the list of channels the bot is a member of
	channels, _, err := api.GetConversations(&slack.GetConversationsParameters{
		Types: []string{"public_channel", "private_channel"},
	})
	if err != nil {
		log.Fatalf("Error getting channel list: %s", err)
	}

	// Loop through each channel and retrieve the message history
	for _, channel := range channels {
		history, err := api.GetConversationHistory(&slack.GetConversationHistoryParameters{
			ChannelID: channel.ID,
		})
		if err != nil {
			log.Printf("Error getting message history for channel %s: %s", channel.Name, err)
			continue
		}

		// Loop through each message and extract any links
		for _, message := range history.Messages {
			// log.Println()
			// log.Println()
			//log.Printf("message: %v\n", message)
			for _, f := range message.Files {
				//log.Printf("\nfile: %v", f)
				log.Printf("private: %s", f.URLPrivateDownload)
			}
			for _, a := range message.Attachments {
				//log.Printf("\nattachment: %v", a)
				if a.ImageURL != "" {
					log.Printf("imageURL: %s", a.ImageURL)
				}
			}
			// links := extractLinks(message.Text)
			// for _, link := range links {
			// 	log.Println(link)
			// }
		}
	}
}

func extractLinks(text string) []string {
	log.Printf("text: %s\n", text)

	// Use a regular expression to find links in the message text
	re := regexp.MustCompile(`(?P<url>https?://[^\s]+)`)
	matches := re.FindAllStringSubmatch(text, -1)

	// Extract the URLs from the matches
	var urls []string
	for _, match := range matches {
		urls = append(urls, match[1])
	}

	return urls
}
