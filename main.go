package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

var kenobiCount = make(map[string]int)

const helloTherePattern = "(?i)^(hello there)[!]?$"

// GrievousBot is the main handler for reddit events
type GrievousBot struct {
	bot reddit.Bot
}

// Comment captures the event that a comment was made in a watched subreddit
func (b *GrievousBot) Comment(c *reddit.Comment) error {
	log.Printf("Received comment with ID [%v] by %vn", c.ID, c.Author)

	r, _ := regexp.Compile(helloTherePattern)
	if !r.MatchString(c.Body) {
		log.Printf("Comment did not match pattern, moving on\n")
		return nil
	}

	kenobiCount[c.Author]++
	count := kenobiCount[c.Author]
	msg := "General Kenobi"

	if count > 1 {
		msg += fmt.Sprintf("\n\n^(We meet again... score: %d)", count)
	}

	log.Printf("Comment with ID [%v] matched pattern, sending reply\n", c.ID)
	return b.bot.Reply(c.Name, msg)
}

func main() {
	bot, err := reddit.NewBotFromAgentFile("bot.agent", 0)
	if err != nil {
		log.Fatalln("Failed to create bot handle: ", err)
	}

	cfg := graw.Config{SubredditComments: []string{"flexicondev"}}
	handler := &GrievousBot{bot: bot}

	_, wait, err := graw.Run(handler, bot, cfg)
	if err != nil {
		log.Fatalln("Failed to start graw run: ", err)
	}

	log.Println("General Grievous standing by...")
	log.Fatalln(wait())
}
