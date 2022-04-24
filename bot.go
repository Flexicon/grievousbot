package main

import (
	"log"
	"regexp"

	"github.com/turnage/graw/reddit"
)

const helloTherePattern = "(?i)^(hello there)[!]?$"

// GrievousBot is the main handler for reddit events
type GrievousBot struct {
	bot      reddit.Bot
	username string
}

func newGrievousBot(bot reddit.Bot, username string) *GrievousBot {
	return &GrievousBot{
		bot:      bot,
		username: username,
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

	log.Printf("Comment with ID [%s] matched pattern, sending reply", c.ID)

	msg := "General Kenobi. You are a bold one."
	reply, err := b.bot.GetReply(c.Name, msg)
	if err != nil {
		log.Printf("Reply to [%s] sent successfully - Link: https://reddit.com%s", c.ID, reply.URL)
	}

	return err
}
