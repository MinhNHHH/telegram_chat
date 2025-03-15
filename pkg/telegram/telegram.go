package telegram

import (
	"log"
	"strings"

	"github.com/MinhNHHH/telegram-bot/pkg/dictionary"
	"github.com/MinhNHHH/telegram-bot/pkg/llm"
	tgbotai "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	Bot         *tgbotai.BotAPI
	ChatGroupID int64
	Dictionary  *dictionary.Dictionary
}

func NewTelegram(botToken string, chatGroupId int64, dict *dictionary.Dictionary) *Telegram {
	teleBot, err := tgbotai.NewBotAPI(botToken)
	if err != nil {
		return nil
	}
	return &Telegram{
		Bot:         teleBot,
		ChatGroupID: chatGroupId,
		Dictionary:  dict,
	}
}

func (t *Telegram) ListenMessage() {
	u := tgbotai.NewUpdate(0)
	u.Timeout = 60

	channels := t.Bot.GetUpdatesChan(u)

	for channel := range channels {
		if channel.Message == nil {
			continue
		}

		// user := channel.Message.From.UserName
		messageText := channel.Message.Text
		reply, err := t.handleMessage(messageText)
		if err != nil {
			log.Fatal(err)
		}
		t.SendMessage(reply)
	}
}

func (t *Telegram) SendMessage(messageText string) {
	reply := tgbotai.NewMessage(t.ChatGroupID, messageText)
	t.Bot.Send(reply)
}

func (t *Telegram) handleMessage(message string) (string, error) {
	messageSplited := strings.Split(message, " ")
	command, text := messageSplited[0], ""
	if len(messageSplited) > 1 {
		text = strings.Join(messageSplited[1:], " ")
	}
	isCommand := strings.HasPrefix(command, "/")
	reply := "Unknown command"
	if isCommand {
		switch command {
		case "/help":
			reply = t.helpCommand()
		case "/search":
			reply = t.searchWordCommand(text)
		case "/fixgrammar":
			reply = t.fixGrammarCommand(text)
		}
	}
	return reply, nil
}

// The help command should return a explain features in  tool
func (t *Telegram) helpCommand() string {
	template := `
		You can control me by sending these commands:
		/search: Search meaning words
		/fixgrammar: Fix english grammar
	`
	return template
}

func (t *Telegram) searchWordCommand(word string) string {
	result, err := t.Dictionary.Search(word)
	reply := "Unknown word"
	if err != nil {
		return reply
	}
	reply = t.Dictionary.FormatDefinition(result)
	return reply
}

func (t *Telegram) fixGrammarCommand(context string) string {
	result, err := llm.CallLLM(context)
	reply := ""
	if err != nil {
		reply = err.Error()
	}
	reply = result
	return reply
}
