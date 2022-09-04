package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

func main() {
	bot, username, err := newRedditBot()
	if err != nil {
		log.Fatalln("Failed to create bot handle: ", err)
	}

	handler := newGrievousBot(bot, username)
	cfg := graw.Config{
		SubredditComments: []string{"flexicondev", os.Getenv("SUBREDDITS")},
		CommentReplies:    true,
	}

	_, wait, err := graw.Run(handler, bot, cfg)
	if err != nil {
		log.Fatalln("Failed to start graw run: ", err)
	}

	go func() {
		if err := runHttpServer(); err != nil {
			log.Fatalln(err)
		}
	}()

	log.Println("General Grievous standing by...")
	if err := wait(); err != nil {
		log.Fatalln(err)
	}
}

func newRedditBot() (reddit.Bot, string, error) {
	username := os.Getenv("CLIENT_USERNAME")

	bot, err := reddit.NewBot(reddit.BotConfig{
		Agent: os.Getenv("USER_AGENT"),
		App: reddit.App{
			ID:       os.Getenv("CLIENT_ID"),
			Secret:   os.Getenv("CLIENT_SECRET"),
			Username: username,
			Password: os.Getenv("CLIENT_PASSWORD"),
		},
		Rate: 0,
	})

	return bot, username, err
}

func runHttpServer() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "General Grievous Bot - https://www.reddit.com/user/gen_grievous_bot\n%s", os.Getenv("USER_AGENT"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("PORT not set, not starting http server")
		return nil
	}

	log.Printf("Grievous http server started on [::]:%s", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
