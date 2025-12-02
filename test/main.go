package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	log.Println("🚀 Бот запущен. Ожидание сообщений...")

	// Проверка токена
	if _, err := getUpdates(0); err != nil {
		log.Printf("❌ Ошибка токена: %v", err)
		return
	}
	log.Println("✅ Токен проверен успешно")

	var offset int
	var mu sync.Mutex

	for {
		mu.Lock()
		updates, err := getUpdates(offset)
		mu.Unlock()

		if err != nil {
			log.Printf("❌ Ошибка получения обновлений: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		for _, update := range updates {
			if update.Message.Text != "" {
				log.Printf("📨 Сообщение: %s", update.Message.Text)
				go processMessage(update.Message)
			}

			mu.Lock()
			if update.UpdateID >= offset {
				offset = update.UpdateID + 1
			}
			mu.Unlock()
		}

		time.Sleep(500 * time.Millisecond)
	}
}
