package lib

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	tgAPIToken = "6722486938:AAHgayUeypJJNWmcjT-auj1Fva76DHEeqYc"
)

func StartBot() {
	bot := initTelegramBot()

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			go handleMessage(bot, update.Message)
		} else if update.CallbackQuery != nil {
			go handleCallbackQuery(bot, update.CallbackQuery)
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf("新消息 %s", message.From.UserName)
	// text := message.Text
	chatID := message.Chat.ID

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("开奖历史", "betting_history"),
			tgbotapi.NewInlineKeyboardButtonData("开奖历史1", "betting_history1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("登陆", "auth_login"),
			tgbotapi.NewInlineKeyboardButtonData("注册", "auth_signin"),
		),
	)

	file, _ := os.Open("../test_image.webp")
	reader := tgbotapi.FileReader{Name: "image.jpg", Reader: file}

	photo := tgbotapi.NewPhoto(chatID, reader)
	photo.Caption = "test caption"

	// replyMsg := tgbotapi.NewMessage(chatID, text+" "+message.Chat.UserName+" "+message.From.UserName+" "+message.From.FirstName+" "+"<a href='https://static0.gamerantimages.com/wordpress/wp-content/uploads/2023/10/collage-maker-11-oct-2023-07-51-pm-2801.jpg'>&#8205;</a>")
	// replyMsg := tgbotapi.NewDiceWithEmoji(chatID, "🎲")
	photo.ReplyMarkup = keyboard
	_, _ = bot.Send(photo)

	// rollDice(bot, chatID, 3)
}

func rollDice(bot *tgbotapi.BotAPI, chatID int64, numDice int) ([]int, error) {
	diceValues := make([]int, numDice)
	diceConfig := tgbotapi.NewDiceWithEmoji(chatID, "🎲")

	for i := 0; i < numDice; i++ {
		diceMsg, err := bot.Send(diceConfig)
		if err != nil {
			log.Println("发送骰子消息异常:", err)
			return nil, err
		}
		diceValues[i] = diceMsg.Dice.Value
	}

	return diceValues, nil
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	log.Printf("Key: %s", callbackQuery.Data)
}

func initTelegramBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(tgAPIToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = true

	log.Printf("已授权帐户 %s", bot.Self.UserName)

	return bot
}
