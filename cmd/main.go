package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"todo-orion-bot/internal/worker"
	"todo-orion-bot/storage"

	"github.com/joho/godotenv"
	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	ch := make(chan string, 3)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := storage.New(os.Getenv("STORAGE_PATH"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbDon, err := storage.New(os.Getenv("STORAGE_PATH_DON"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbPah, err := storage.New(os.Getenv("STORAGE_PATH_PAH"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	worker.Created(db, os.Getenv("MAX_SHEET"), os.Getenv("RANGE_SHEET"), os.Getenv("MAX_TELEGRAM"), ch)

	worker.Created(dbDon, os.Getenv("DON_SHEET"), os.Getenv("RANGE_SHEET"), os.Getenv("DON_TELEGRAM"), ch)

	worker.Created(dbPah, os.Getenv("PAH_SHEET"), os.Getenv("RANGE_SHEET"), os.Getenv("PAH_TELEGRAM"), ch)

	bot, err := tg.NewBot(os.Getenv("BOT_TOKEN"), tg.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)

	defer bot.StopLongPolling()

	chatId, _ := strconv.Atoi(os.Getenv("CHAT_ID"))
	for update := range updates {
		if update.Message != nil {
			for v := range ch {
				message := tu.Message(tu.ID(int64(chatId)), v)
				_, _ = bot.SendMessage(message)
			}

		} else {
			_, _ = bot.SendMessage(tu.Message(tu.ID(update.Message.Chat.ID), ""))
		}
	}

}
