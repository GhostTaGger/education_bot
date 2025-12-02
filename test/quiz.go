package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Структуры викторины
type Question struct {
	ID      int
	Text    string
	Answers []Answer
}

type Answer struct {
	ID         int
	QuestionID int
	Text       string
	IsCorrect  bool
}

// Данные вопросов
var questions = []Question{
	{
		ID:   1,
		Text: "Столица России?",
		Answers: []Answer{
			{ID: 1, QuestionID: 1, Text: "Москва", IsCorrect: true},
			{ID: 2, QuestionID: 1, Text: "Санкт-Петербург", IsCorrect: false},
			{ID: 3, QuestionID: 1, Text: "Киев", IsCorrect: false},
			{ID: 4, QuestionID: 1, Text: "Минск", IsCorrect: false},
		},
	},
	{
		ID:   2,
		Text: "2 + 2 = ?",
		Answers: []Answer{
			{ID: 5, QuestionID: 2, Text: "4", IsCorrect: true},
			{ID: 6, QuestionID: 2, Text: "5", IsCorrect: false},
			{ID: 7, QuestionID: 2, Text: "3", IsCorrect: false},
			{ID: 8, QuestionID: 2, Text: "22", IsCorrect: false},
		},
	},
	{
		ID:   3,
		Text: "Самая большая планета Солнечной системы?",
		Answers: []Answer{
			{ID: 9, QuestionID: 3, Text: "Юпитер", IsCorrect: true},
			{ID: 10, QuestionID: 3, Text: "Земля", IsCorrect: false},
			{ID: 11, QuestionID: 3, Text: "Марс", IsCorrect: false},
			{ID: 12, QuestionID: 3, Text: "Сатурн", IsCorrect: false},
		},
	},
}

// Функции викторины
func showAllQuestions() string {
	var result strings.Builder
	result.WriteString("📝 Вопросы:\n\n")

	for _, q := range questions {
		result.WriteString(fmt.Sprintf("🆔 %d. %s\n", q.ID, q.Text))
		for _, a := range q.Answers {
			result.WriteString(fmt.Sprintf("   %d) %s\n", a.ID, a.Text))
		}
		result.WriteString("\n")
	}
	result.WriteString("💡 /answer [номер_вопроса] [номер_ответа]")

	return result.String()
}

func getRandomQuestion() string {
	if len(questions) == 0 {
		return "❌ Нет вопросов"
	}

	q := questions[0] // Первый как случайный
	result := fmt.Sprintf("🎲 Вопрос: %s\n\nВарианты:\n", q.Text)

	for _, a := range q.Answers {
		result += fmt.Sprintf("%d. %s\n", a.ID, a.Text)
	}
	result += fmt.Sprintf("\n💡 Ответ: /answer %d [номер]", q.ID)

	return result
}

func handleAnswerCommand(chatID, command string) {
	parts := strings.Fields(command)
	if len(parts) != 3 {
		sendMessage(chatID, "❌ Формат: /answer [номер_вопроса] [номер_ответа]")
		return
	}

	qID, err1 := strconv.Atoi(parts[1])
	aID, err2 := strconv.Atoi(parts[2])

	if err1 != nil || err2 != nil {
		sendMessage(chatID, "❌ Номера должны быть числами")
		return
	}

	// Находим вопрос
	var question *Question
	for _, q := range questions {
		if q.ID == qID {
			question = &q
			break
		}
	}

	if question == nil {
		sendMessage(chatID, "❌ Вопрос не найден")
		return
	}

	// Проверяем ответ
	var isCorrect bool
	var correctAnswer string

	for _, a := range question.Answers {
		if a.ID == aID {
			isCorrect = a.IsCorrect
		}
		if a.IsCorrect {
			correctAnswer = a.Text
		}
	}

	if isCorrect {
		sendMessage(chatID, "✅ Правильно! 🎉")
	} else {
		sendMessage(chatID, fmt.Sprintf("❌ Неправильно! Правильно: %s", correctAnswer))
	}
}
