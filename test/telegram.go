package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Структуры
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int64  `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

type Chat struct {
	ID int64 `json:"id"`
}

// Глобальные переменные
var botToken string

// Инициализация
func init() {
	if token := os.Getenv("TELEGRAM_BOT_TOKEN"); token != "" {
		botToken = token
	} else {
		botToken = "ВСТАВЬ ТОКЕН"
		log.Println("⚠️  Используется демо-токен")
	}
}

// Telegram API функции
const telegramAPI = "https://api.telegram.org/bot"

func sendMessage(chatID string, text string) error {
	url := telegramAPI + botToken + "/sendMessage"

	reqBody := map[string]string{
		"chat_id": chatID,
		"text":    text,
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Telegram API: %s", string(body))
	}
	return nil
}

func getUpdates(offset int) ([]Update, error) {
	url := fmt.Sprintf("%s%s/getUpdates?offset=%d&timeout=30", telegramAPI, botToken, offset)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	json.Unmarshal(body, &result)
	return result.Result, nil
}

// Обработчик сообщений
func processMessage(msg Message) {
	chatID := strconv.FormatInt(msg.Chat.ID, 10)
	text := msg.Text

	if text == "" {
		return
	}

	switch {
	case text == "/start":
		sendMessage(chatID, `🎯 Добро пожаловать в викторину!
Доступные команды:
/questions — показать все вопросы
/quiz — случайный вопрос
/answer [номер_вопроса] [номер_ответа] — ответить на вопрос`)

	case text == "/questions":
		sendMessage(chatID, showAllQuestions())

	case text == "/quiz":
		sendMessage(chatID, getRandomQuestion())

	case len(text) >= 8 && text[:7] == "/answer":
		handleAnswerCommand(chatID, text)

	default:
		sendMessage(chatID, "❌ Неизвестная команда. Напишите /start")
	}
}
