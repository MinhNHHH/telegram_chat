package settings

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BOT_TOKEN          string
	CHAT_GROUP_ID      int64
	DICTIONARY_API_URL string
}

func LoadEnv() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("TELE_BOT_TOKEN")
	chatGroupId, err := strconv.Atoi(os.Getenv("TELE_GROUP_ID"))
	dictionaryApiUrl := os.Getenv("DICTIONARY_API_URL")
	if err != nil {
		chatGroupId = 0
		log.Fatal("Chat group id should be number")
	}

	return &Config{
		BOT_TOKEN:          botToken,
		CHAT_GROUP_ID:      int64(chatGroupId),
		DICTIONARY_API_URL: dictionaryApiUrl,
	}

}
