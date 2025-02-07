package tg

import (
	"librelink-up-tg/config"
	"librelink-up-tg/internal/clients/libre"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/sync/errgroup"
)

type Client struct {
	config *config.Config
}

func NewClient(config *config.Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) SendToFriends(libreData *libre.GraphData) error {
	bot, err := tgbotapi.NewBotAPI(c.config.BotFatherToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	eg := errgroup.Group{}
	libreMessage := libreData.String()
	log.Printf("Current state:\n%s", libreMessage)

	for _, id := range c.config.ChatIDsToNotify {
		id := id

		eg.Go(func() error {
			msg := tgbotapi.NewMessage(id, libreMessage)

			_, err = bot.Send(msg)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
