package bbtmvbot

import (
	"bbtmvbot/config"
	"bbtmvbot/database"
	"bbtmvbot/website"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	telebot "gopkg.in/tucnak/telebot.v2"
)

var (
	db        *database.Database
	tb        *telebot.Bot
	allChatID int64
)

func Start(c *config.Config, dbPath *string) {
	// Open DB
	var err error
	db, err = database.Open(*dbPath)
	if err != nil {
		log.Fatalln(err)
	}

	allChatID = c.Telegram.AllChatID

	// Connect to Telegram
	poller := &telebot.LongPoller{Timeout: 10 * time.Second}
	middlewarePoller := telebot.NewMiddlewarePoller(poller, func(upd *telebot.Update) bool {
		if upd.Message == nil {
			return false
		}
		db.EnsureUserInDB(upd.Message.Chat.ID) // This ensures that user is always in DB
		return true
	})
	tb, err = telebot.NewBot(telebot.Settings{Token: c.Telegram.ApiKey, Poller: middlewarePoller})
	if err != nil {
		log.Fatalln(err)
	}
	initTelegramHandlers()

	// Start telegram bot
	go tb.Start()

	// Setup cronjob
	location, _ := time.LoadLocation("Europe/Vilnius")
	s := gocron.NewScheduler(location)
	s.Every("3m").Do(refreshWebsites) // Retrieve new posts, send to users
	s.Every("24h").Do(cleanup)        // Cleanup (remove posts that are not seen in the last 30 days)

	// Start cronjob and block execution
	s.StartBlocking()
}

func refreshWebsites() {
	for title, site := range website.Websites {

		go func(title string, site website.Website) {
			posts := site.Retrieve(db)
			for _, post := range posts {
				go processPost(post)
			}
		}(title, site)

	}
}

func processPost(post *website.Post) {
	excludable := post.IsExcludable()
	insertedPostID := db.AddPost(post.Link)
	msg := post.FormatTelegramMessage(insertedPostID)

	allMsg := msg + "\n» *Excluded:* "
	if excludable {
		allMsg += "Yes"
	} else {
		allMsg += "No"
	}
	sendTelegram(allChatID, allMsg)

	if excludable {
		return
	}

	telegramIDs := db.GetInterestedTelegramIDs(post.Price, post.Rooms, post.Year)
	for _, telegramID := range telegramIDs {
		sendTelegram(telegramID, msg)
	}
}

func cleanup() {
	db.DeleteOldPosts() // Older than 30 days
}
