package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
	"strconv"
	"strings"
)

func postToTelegram(exams TimedExams, telegramChannel string, log zerolog.Logger) {
	if config.BotToken == "" || telegramChannel == "" || len(exams) == 0 {
		return
	}

	channel, err := strconv.ParseInt(telegramChannel, 10, 64)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse bot token")
	}

	log.Debug().Msg("sending message to telegram channel")

	telegramBot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create telegram bot")
	}

	telegramBot.Debug = *debugFlag

	messages := createDiffMessages(exams)
	log.Debug().Int("n", len(messages)).Msg("messages to send")
	for _, msg := range messages {
		msg := tgbotapi.NewMessage(channel, msg)
		msg.ParseMode = tgbotapi.ModeMarkdown
		_, err = telegramBot.Send(msg)
		if err != nil {
			log.Fatal().Err(err).Msg("could not send message to telegram channel")
		}
	}
}

func createDiffMessages(diff TimedExams) []string {
	var b strings.Builder
	messages := make([]string, 0, len(diff)/10+1)

	b.WriteString("📅 **Nuovi appelli pubblicati:**\n")

	for i, exam := range diff {
		b.WriteString("\n")

		b.WriteString(fmt.Sprintf("**%s**", exam.Exam.Course.Title))

		if exam.Exam.Course.Code != "" {
			b.WriteString(fmt.Sprintf(" (%s)", exam.Exam.Course.Code))
		}

		if exam.Exam.Course.Teacher != "" {
			b.WriteString(fmt.Sprintf("\n%s", exam.Exam.Course.Teacher))
		}

		b.WriteString(fmt.Sprintf("\n%s\n\n", exam.Exam.Time.Format("02/01/2006 15:04")))

		// Send message every 10 exams
		if i%10 == 0 && i != 0 {
			messages = append(messages, b.String())
			b.Reset()
		}
	}

	// Add last message
	if b.Len() > 0 {
		messages = append(messages, b.String())
	}

	return messages
}
