package main

import (
	"fmt"
	"queuebot/queue"
	"strconv"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	CmdIDGoToQueue = iota
	CmdIDSkip
	CmdIDCheck
	CmdIDInsert
	CmdIDOut
	CmdIDDone

	CmdUnknown
)

const (
	CmdTextGoToQueue = "Мне бы сдать "
	CmdTextSkip      = "Отдохнуть от безумия этого мира"
	CmdTextCheck     = "Что по очереди?"
	CmdTextInsert    = "Занять за "
	CmdTextOut       = "Выйти"
	CmdTextDone      = "Сдал!"

	CmdTextUnknown = ""
)

const (
	StateChooseQueue = iota
	StateInQueue
)

type command struct {
	commandID int
	text      string
	addText   string
}

var CmdGoToQueue = command{
	commandID: CmdIDGoToQueue,
	text:      CmdTextGoToQueue,
}
var CmdSkip = command{
	commandID: CmdIDSkip,
	text:      CmdTextSkip,
}
var CmdCheck = command{
	commandID: CmdIDCheck,
	text:      CmdTextCheck,
}
var CmdOut = command{
	commandID: CmdIDOut,
	text:      CmdTextOut,
}
var CmdDone = command{
	commandID: CmdIDDone,
	text:      CmdTextDone,
}

var commands = []command{
	CmdGoToQueue,
	CmdSkip,
	CmdCheck,
	CmdOut,
}

type User struct {
	TgNick string
	ChatID int64
}

type StateController struct {
	userStateMap      map[string]int
	userClassMap      map[string]string
	classNameQueueMap map[string]*queue.Queue
}

func NewController() (*StateController, error) {
	sc := &StateController{
		userStateMap:      make(map[string]int),
		userClassMap:      make(map[string]string),
		classNameQueueMap: make(map[string]*queue.Queue),
	}
	for _, class := range classes {
		queue, err := queue.NewQueue(class)
		if err != nil {
			return nil, err
		}
		sc.classNameQueueMap[class] = queue
	}

	return sc, nil
}

func (sc *StateController) Save() {

}

func (sc *StateController) generateList(class string) string {
	usrs := sc.classNameQueueMap[class].All()
	res := ""
	for i, usr := range usrs {
		res += strconv.Itoa(i+1) + ". @" + usr.TgNick + "\n"
	}
	return res
}

func (sc *StateController) Process(s string, tgNick string, chatID int64) (*tgbotapi.ReplyKeyboardMarkup, string) {
	if _, ok := sc.userStateMap[tgNick]; !ok {
		sc.userStateMap[tgNick] = StateChooseQueue
		return queueChooseQueueKeyboard, "Вэлкам!"
	}

	c := Parse(s)
	if c.commandID == CmdUnknown {
		return nil, "418: закипел"
	}
	switch sc.userStateMap[tgNick] {
	case StateChooseQueue:
		sc.userStateMap[tgNick] = StateInQueue
		if c.commandID == CmdIDGoToQueue {
			if class, ok := sc.classNameQueueMap[c.text]; ok {
				sc.userClassMap[tgNick] = c.text
				switch _, userState := class.GetUser(tgNick); userState {
				case queue.UserInQueue:
					return createCommandKeyboard(CmdSkip, CmdCheck, CmdOut), sc.generateList(c.text)
				case queue.UserSkipping:
					return inQueueKeyboard, ""
				case queue.UserNotRegisterd:
					class.Add(queue.User{
						TgNick: tgNick,
						ChatID: chatID,
					})
					return inQueueKeyboard, "Добавил тебя :)\n" + sc.generateList(c.text)
				}
			}
		}

	case StateInQueue:
		switch c.commandID {
		case CmdIDOut:
			sc.userStateMap[tgNick] = StateChooseQueue
			return queueChooseQueueKeyboard, "Выбирай, куда хочешь пойти"

		case CmdIDCheck:
			return inQueueKeyboard, "Очередь на " + sc.userClassMap[tgNick] + "\n" + sc.generateList(sc.userClassMap[tgNick])

		case CmdIDSkip:

		}
	}

	return nil, "418: закипел: badly"
}

func Parse(s string) command {
	for _, cmd := range commands {
		if strings.HasPrefix(s, cmd.text) {
			resp := cmd
			resp.text = strings.TrimPrefix(s, cmd.text)
			fmt.Println("known: ", resp)
			return resp
		}
	}
	fmt.Println("unknown: :(")
	return command{
		commandID: CmdUnknown,
	}
}
