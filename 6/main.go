package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const port = ":8080"

var (
	upgrader = websocket.Upgrader{ // создание объекта Upgrader для преобразования HTTP соединений в WebSocket
		CheckOrigin: func(r *http.Request) bool {
			return true // разрешение соединения из любого источника
		},
	}
	clients   = make(map[*websocket.Conn]bool) // хранение подключенных клиентов
	broadcast = make(chan string)              // канал для отправки сообщений между клиентами
	mu        sync.Mutex                       // мьютекс для синхронизации доступа к clients
)

func main() {
	/*
			Веб-сокеты:
		•	Реализуйте сервер на основе веб-сокетов для чата.
		•	Клиенты должны подключаться к серверу, отправлять и получать сообщения.
		•	Сервер должен поддерживать несколько клиентов и рассылать им сообщения, отправленные любым подключённым клиентом.

	*/

	http.HandleFunc("/ws", handleConnections) // устанавливаем обработчик для WebSocket соединений на маршруте /ws
	go handleMessages()                       // запускаем горутину для обработки входящих сообщений

	fmt.Println("сервер запущен на ", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("ошибка запуска сервера:", err) // обработка ошибки запуска сервера
	}
}

// устанавливает соединение с клиентом
func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // обновление соединения до WebSocket
	if err != nil {                          // проверка на ошибки
		fmt.Println("ошибка при установке соединения:", err)
		return
	}
	defer conn.Close() // закрытие соединения при выходе из функции

	// добавление клиента в массив подключенных клиентов
	mu.Lock()            // защита доступа к shared resource
	clients[conn] = true // добавление клиента в список
	mu.Unlock()          // освобождение блокировки

	// бесконечный цикл для получения сообщений от клиента
	for {
		_, msg, err := conn.ReadMessage() // читение сообщения от клиента
		if err != nil {                   // проверка на ошибки
			fmt.Println("ошибка чтения сообщения:", err)
			break
		}
		broadcast <- string(msg) // отправка сообщения в канал для рассылки
	}
}

// обрабатывает все сообщения и рассылает их всем клиентам
func handleMessages() {
	for {
		msg := <-broadcast            // получение сообщения из канала
		mu.Lock()                     // защита доступа к shared resource
		for client := range clients { // проход по всем подключенным клиентам
			if err := client.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil { // отправление сообщения клиенту
				fmt.Println("ошибка отправки сообщения:", err) // проверка на ошибки
				client.Close()                                 // закрытие соединения с клиентом
				delete(clients, client)                        // удаление клиента из списка
			}
		}
		mu.Unlock() // освобождение блокировки
	}
}
