// build for heroku:
// env GOOS=linux GOARCH=amd64 go build -o bin/go_build_tg_dao_bot_src -v ./src/
package main

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	godotenv.Load()

	rand.Seed(time.Now().Unix())
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Получить наставление Дао")))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	quoteFile, err := os.Open("Dao.json")
	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(quoteFile)
	quoteStruct := struct {
		Quotes []string `json:"quotes"`
	}{}
	dec.Decode(&quoteStruct)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyMarkup = keyboard
		switch update.Message.Text {
		case "/start":
			msgText := "Этот бот присылает наставления из Дао Де Цзинь"
			msg.Text = msgText
		case "Получить наставление Дао":
			quoteNum := rand.Int() % 82
			msgText := strconv.Itoa(quoteNum) + ":\n" + quoteStruct.Quotes[quoteNum]
			msg.Text = msgText
		}

		bot.Send(msg)
	}
}
