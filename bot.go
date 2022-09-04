package main

import (
	"log"
	"math/rand"
	"regexp"
	"time"

	"github.com/turnage/graw/reddit"
)

const (
	helloTherePattern = "(?i)^(hello there)[!]*$"
	helloThereMsg     = "General Kenobi. You are a bold one."
)

var (
	replyQuotes = []string{
		"That wasn't much of a rescue.",
		"I will deal with this Jedi slime myself.",
		"Jedi slime - Your comment will make a fine addition to my collection!",
		"Time to abandon ship.",
		"Army or not, you must realize, you are doomed.",
		"Your comment will make a fine addition to my collection!",
		"Your lightsabers will make a fine addition to my collection!",
	}
)

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

	if !isHelloThereMessage(c.Body) {
		log.Printf("Comment [%s] did not match pattern, moving on", c.ID)
		return nil
	}

	log.Printf("Comment with ID [%s] matched pattern, sending reply", c.ID)

	reply, err := b.bot.GetReply(c.Name, helloThereMsg)
	if err == nil {
		log.Printf("Reply to [%s] sent successfully - Link: https://reddit.com%s", c.ID, reply.URL)
	}

	return err
}

// CommentReply captures the event that a comment reply was made to the bot
func (b *GrievousBot) CommentReply(r *reddit.Message) error {
	if r.Author == b.username || isHelloThereMessage(r.Body) {
		return nil
	}
	log.Printf("Received reply to comment with ID [%s] by [%s] - Link: https://reddit.com%s", r.ID, r.Author, r.Context)

	newReply, err := b.bot.GetReply(r.Name, randomReplyQuote())
	if err == nil {
		log.Printf("Reply to [%s] sent successfully - Link: https://reddit.com%s", r.ID, newReply.URL)
	}

	return err
}

func isHelloThereMessage(msg string) bool {
	r, _ := regexp.Compile(helloTherePattern)
	return r.MatchString(msg)
}

func randomReplyQuote() string {
	rand.Seed(time.Now().Unix())
	return replyQuotes[rand.Intn(len(replyQuotes))]
}
