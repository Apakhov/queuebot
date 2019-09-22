package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var mainModKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("SKIP"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("DOWN"),
	),
)

var skipModKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("INSERT"),
	),
)

func createKeyboard(prefix string, delim int, strs ...string) *tgbotapi.ReplyKeyboardMarkup {
	btrows := make([][]tgbotapi.KeyboardButton, 0)

	for i, str := range strs {
		text := str
		if delim >= 0 && i >= delim {
			text = prefix + text
		}
		btrows = append(btrows,
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(text),
			),
		)
	}

	k := tgbotapi.NewReplyKeyboard(btrows...)
	return &k
}

func createCommandKeyboard(cmds ...command) *tgbotapi.ReplyKeyboardMarkup {
	btrows := make([][]tgbotapi.KeyboardButton, 0)

	for _, cmd := range cmds {
		btrows = append(btrows,
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(cmd.text+cmd.addText),
			),
		)
	}

	k := tgbotapi.NewReplyKeyboard(btrows...)
	return &k
}

var queueChooseQueueKeyboard *tgbotapi.ReplyKeyboardMarkup
var inQueueKeyboard *tgbotapi.ReplyKeyboardMarkup

func init() {
	queueChooseQueueKeyboard = createKeyboard(CmdGoToQueue.text, 0, classes...)
	inQueueKeyboard = createCommandKeyboard(CmdCheck, CmdSkip, CmdOut)
}
