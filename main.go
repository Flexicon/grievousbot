package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

var (
	RequiredEnvVars = []string{"CLIENT_USERNAME", "CLIENT_SECRET", "CLIENT_ID", "CLIENT_PASSWORD", "USER_AGENT"}
)

func main() {
	ensureEnvironmentVariablesPresent(RequiredEnvVars)

	setupSentry()
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	bot, username, err := newRedditBot()
	if err != nil {
		log.Fatalln("Failed to create bot handle: ", err)
		sentry.CaptureException(err)
	}

	handler := newGrievousBot(bot, username)
	cfg := graw.Config{
		SubredditComments: []string{"flexicondev", os.Getenv("SUBREDDITS")},
		CommentReplies:    true,
	}

	_, wait, err := graw.Run(handler, bot, cfg)
	if err != nil {
		log.Fatalln("Failed to start graw run: ", err)
		sentry.CaptureException(err)
	}

	go func() {
		if err := runHttpServer(); err != nil {
			log.Fatalln(err)
			sentry.CaptureException(err)
		}
	}()

	log.Println("General Grievous standing by...")
	if err := wait(); err != nil {
		log.Fatalln(err)
		sentry.CaptureException(err)
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

func ensureEnvironmentVariablesPresent(vars []string) {
	for _, v := range vars {
		if _, ok := os.LookupEnv(v); !ok {
			log.Fatalf("Missing environment variable '%s'", v)
		}
	}
}

func setupSentry() {
	dsn, ok := os.LookupEnv("SENTRY_DSN")
	if !ok {
		log.Println("Skipping Sentry setup - SENTRY_DSN is not set")
		return
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		TracesSampleRate: 1.0,
		AttachStacktrace: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}
