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

func forwardMessageToChannel(bot *tgbotapi.BotAPI, chatID int64, update tgbotapi.Update) {
	message := fmt.Sprintf(
		"@%s (%s %s) is requesting an airdrop to <code>%s</code>",
		update.Message.From.UserName,
		update.Message.From.FirstName,
		update.Message.From.LastName,
		update.Message.Text,
	)
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := bot.Send(msg)
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

func Main() {
	bot, err := tgbotapi.NewBotAPI(TELEGRAM_BOT_TOKEN)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// reply in code format
		msg.ParseMode = tgbotapi.ModeMarkdownV2

		if update.Message.Chat.Type != "private" { // ignore any group chat Message updates
			continue
		}

		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			log.Println(pretty.JSONString(update))
			if validateAddress(update.Message.Text) {
				msg.Text = airdrop(update.Message.Text)
				bot.Send(msg)
				forwardMessageToChannel(bot, TELEGRAM_CHANNEL_ID, update)
				continue
			}
			msg.Text = fmt.Sprintf("Error: invalid SS58 address: %s", pretty.JSONString(update.Message.Text))
			bot.Send(msg)
		}

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /airdrop /sayhi and /status."
		case "sayhi":
			msg.Text = "Hi :)"
		case "status":
			msg.Text = "I'm ok."
		case "validator":
			msg.Text = "https://hackmd.io/@gearvara/validator"
		case "airdrop":
			log.Println(pretty.YAMLString(update.Message.Text))
			log.Println(pretty.YAMLString(update.Message.From))
			msg.Text = "Your request has been submitted to @GearvaraBotAirdropQueue. It should be approved by @btwiuse shortly. If you didn't receive testnet tokens within 24 hours, please leave a message in @GearvaraBotDiscussion. Thank you!"
			forwardMessageToChannel(bot, TELEGRAM_CHANNEL_ID, update)
		default:
			msg.Text = "Please enter your SS58 address to receive the airdrop on [Vara testnet](https://polkadot.js.org/apps/?rpc=wss://vara.gear.rs), for example: `5CtLwzLdsTZnyA3TN7FUV58FV4NZ1tUuTDM9yjwRuvt6ac1i`\n\nThe testnet tokens are not transferrable, but you can stake them and become a nominator or /validator on Vara testnet."
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
