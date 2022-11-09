package gearvarabot

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/btwiuse/pretty"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rjman-ljm/go-substrate-crypto/ss58"
)

var TELEGRAM_BOT_TOKEN = os.Getenv("TELEGRAM_BOT_TOKEN")
var TELEGRAM_CHANNEL_ID int64 = -1001849103490

// validate ss58 address
func validateAddress(address string) bool {
	pub, err := ss58.DecodeToPub(address)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println(pretty.YAMLString(pub))
	return true
}

func forwardMessageToChannel(chatID int64, update tgbotapi.Update) {
	message := fmt.Sprintf(
		"@%s (%s %s) is requesting an airdrop to `%s`",
		update.Message.From.UserName,
		update.Message.From.FirstName,
		update.Message.From.LastName,
		update.Message.Text,
	)
	err := sendMarkdown(chatID, message)
	if err != nil {
		log.Println(err)
	}
}

func airdrop(address string) string {
	cmd := exec.Command("airdrop.ts", address)
	cmd.Stderr = os.Stderr
	// get command output
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return "An error occured on our side. Please try again later."
	}
	return string(output)
}

var (
	Bot *tgbotapi.BotAPI
)

func initBot() {
	var err error
	Bot, err = tgbotapi.NewBotAPI(TELEGRAM_BOT_TOKEN)
	if err != nil {
		log.Panic(err)
	}

	Bot.Debug = true

	log.Printf("Authorized on account %s", Bot.Self.UserName)
}

func Main() {
	initBot()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := Bot.GetUpdatesChan(u)

	for update := range updates {
		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		chatID := update.Message.Chat.ID
		chatType := update.Message.Chat.Type
		reply := ""

		if chatType != "private" { // ignore any group chat Message updates
			continue
		}

		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			log.Println(pretty.JSONString(update))
			if validateAddress(update.Message.Text) {
				reply = airdrop(update.Message.Text)
				err := sendMarkdown(chatID, reply)
				if err != nil {
					log.Println(err)
				}
				forwardMessageToChannel(TELEGRAM_CHANNEL_ID, update)
				continue
			}
			reply = fmt.Sprintf("Error: invalid SS58 address: %s", pretty.JSONString(update.Message.Text))
			err := sendMarkdown(chatID, reply)
			if err != nil {
				log.Println(err)
			}
		}

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			reply = "I understand /airdrop /sayhi and /status."
		case "sayhi":
			reply = "Hi :)"
		case "status":
			reply = "I'm ok."
		case "validator":
			reply = "https://hackmd.io/@gearvara/validator"
		case "airdrop":
			log.Println(pretty.YAMLString(update.Message.Text))
			log.Println(pretty.YAMLString(update.Message.From))
			reply = "TODO: unimplemented"
		default:
			reply = "Please enter your SS58 address to receive the airdrop on [Vara testnet](https://polkadot.js.org/apps/?rpc=wss://vara.gear.rs), for example: `5CtLwzLdsTZnyA3TN7FUV58FV4NZ1tUuTDM9yjwRuvt6ac1i`\n\nThe testnet tokens are not transferrable, but you can stake them and become a nominator or /validator on Vara testnet."
		}

		err := sendMarkdown(chatID, reply)
		if err != nil {
			log.Println(err)
		}
	}
}

func sendMarkdown(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	msg.Text = tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, msg.Text)
	_, err := Bot.Send(msg)
	return err
}
