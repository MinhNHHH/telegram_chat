package main

import (
	"fmt"

	"github.com/MinhNHHH/telegram-bot/pkg/dictionary"
	"github.com/MinhNHHH/telegram-bot/pkg/settings"
	"github.com/MinhNHHH/telegram-bot/pkg/telegram"
)

func main() {
	configs := settings.LoadEnv()
	dict := dictionary.NewDictionary(fmt.Sprintf("%s/%s", configs.DICTIONARY_API_URL, "en"))
	telegram := telegram.NewTelegram(configs.BOT_TOKEN, configs.CHAT_GROUP_ID, dict)
	telegram.ListenMessage()
}
