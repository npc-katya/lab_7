package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gorilla/websocket" // для работы с веб-сокетами
)

func main() {
	serverAddr := "ws://127.0.0.1:8080/ws"                        // адрес сервера
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil) // установка соединения с сервером
	if err != nil {                                               // проверка на наличие ошибок
		fmt.Println("ошибка подключения к серверу:", err)
		return
	}
	defer conn.Close() // закрытие соединение при выходе из функции

	go listenForMessages(conn)

	// считывание ввода пользователя
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("для выхода введите 'exit'")
	fmt.Println("введите сообщения:")
	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "exit" { // условие для выхода
			break
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil { // отправляем сообщение на сервер
			fmt.Println("ошибка отправки сообщения:", err)
			break
		}
	}
}

// слушает сообщения от сервера и выводит их на экран
func listenForMessages(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage() // чтение сообщения от сервера
		if err != nil {                   // проверка на наличие ошибок
			fmt.Println("ошибка чтения сообщения:", err)
			return
		}
		fmt.Println("получено:", string(msg)) // вывод полученного сообщения на экран
	}
}
