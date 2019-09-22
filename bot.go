package main

import (
	"log"
	"sync"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type queueBot struct {
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
	stCtrl  *StateController

	stopCh chan struct{}
	wg     sync.WaitGroup
}

func newQueueBot(token string, path string) (*queueBot, error) {
	q := queueBot{}
	var err error
	q.bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	q.bot.Debug = true
	log.Printf("Authorized on account %s", q.bot.Self.UserName)

	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	q.updates, err = q.bot.GetUpdatesChan(ucfg)
	if err != nil {
		return nil, err
	}

	q.stCtrl, err = NewController()
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (q *queueBot) Save() {

}

func (q *queueBot) Stop() {
	q.stCtrl.Save()

	q.stopCh <- struct{}{}
}

func (q *queueBot) Wait() {
	q.stopCh = make(chan struct{}, 0)
	go func() {
		q.wg.Wait()
		q.Stop()
	}()

	<-q.stopCh
}

func (q *queueBot) ListenUsers() {
	q.wg.Add(1)
	go func() {
		for {
			select {
			case update, ok := <-q.updates:
				if !ok {
					q.wg.Done()
					return
				}
				UserName := update.Message.From.UserName
				ChatID := update.Message.Chat.ID
				Text := update.Message.Text

				log.Printf("[%s] %d %s", UserName, ChatID, Text)

				kb, reply := q.stCtrl.Process(Text, UserName, ChatID)
				msg := tgbotapi.NewMessage(ChatID, reply)
				if kb != nil {
					msg.ReplyMarkup = *kb
				}

				q.bot.Send(msg)
			}
		}
	}()
}
