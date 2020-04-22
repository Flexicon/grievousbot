package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/turnage/graw/reddit"
)

const helloTherePattern = "(?i)^(hello there)[!]?$"

// GrievousBot is the main handler for reddit events
type GrievousBot struct {
	bot         reddit.Bot
	username    string
	kenobiCount map[string]int
}

func newGrievousBot(bot reddit.Bot, username string) *GrievousBot {
	return &GrievousBot{
		bot:         bot,
		username:    username,
		kenobiCount: make(map[string]int),
	}
}

// Comment captures the event that a comment was made in a watched subreddit
func (b *GrievousBot) Comment(c *reddit.Comment) error {
	if c.Author == b.username {
		return nil
	}
	log.Printf("Received comment with ID [%s] by [%s] - Link: https://reddit.com%s", c.ID, c.Author, c.Permalink)

	r, _ := regexp.Compile(helloTherePattern)
	if !r.MatchString(c.Body) {
		log.Printf("Comment [%s] did not match pattern, moving on", c.ID)
		return nil
	}

	b.kenobiCount[c.Author]++
	count := b.kenobiCount[c.Author]
	msg := "General Kenobi"
	if count > 1 {
		msg += fmt.Sprintf("\n\nWe meet again /u/%s... (%d times now)", c.Author, count)
	}

	log.Printf("Comment with ID [%s] matched pattern, sending reply", c.ID)
	reply, err := b.bot.GetReply(c.Name, msg)
	if err != nil {
		log.Printf("Reply to [%s] sent successfully - Link: https://reddit.com%s", c.ID, reply.URL)
	}

	return err
}
