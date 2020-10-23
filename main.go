package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const tgtoken = ""

func main() {
	bot, err := tgbotapi.NewBotAPI(tgtoken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		log.Println("command:", update.Message.Command())

		switch update.Message.Command() {
		case "eddie":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			fmt.Println("command:", update.Message.Command())
			bot.Send(msg)
		case "ipcheck":
			var reslut string
			cmd := strings.Split(update.Message.CommandArguments(), " ")
			if len(cmd) != 2 {
				reslut = "參數錯誤 Ex. /ipcheck {ip} {mask}"
			} else {
				if IPCheck(cmd[0], cmd[1]) {
					reslut = "IP在CIDR中"
				} else {
					reslut = "IP不在CIDR中"
				}
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reslut)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "error cmd")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}

//IPCheck 是否IP在CIDR遮罩中
func IPCheck(ip string, mask string) bool {
	checkedIPNet, _ := parseIP(ip)
	//判斷遮罩
	_, ipnetBlackList, parseErr := net.ParseCIDR(mask)
	if parseErr != nil {
		fmt.Println(parseErr)
		return false
	}

	checkResult := ipnetBlackList.Contains(checkedIPNet)

	return checkResult
}

//是否是IPv6 IPv4
func parseIP(s string) (net.IP, int) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, 0
	}

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return ip, 4
		case ':':
			return ip, 6
		}
	}
	return nil, 0
}
