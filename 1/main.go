package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const port = ":8080"

// считывание и вывод сообщения от клиента
func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("получено сообщение:", message)

		_, err := io.WriteString(conn, "сообщение получено!\n")
		if err != nil {
			fmt.Println("ошибка отправки ответа:", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("ошибка сканирования:", err)
	}
}

func main() {

	/*
			Создание TCP-сервера:
		•	Реализуйте простой TCP-сервер, который слушает указанный порт и принимает входящие соединения.
		•	Сервер должен считывать сообщения от клиента и выводить их на экран.
		•	По завершении работы клиенту отправляется ответ с подтверждением получения сообщения.

	*/

	// создание и запуск сервера
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("ошибка создания слушателя:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("сервер запущен на порту %s...\n", port)

	// Обработка сигналов для graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// завершение работы сервера
	go func() {
		<-c
		fmt.Println("завершение работы сервера...")
		listener.Close()
		os.Exit(0)
	}()

	// приём сообщения от клиента
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("ошибка принятия соединения:", err)
			continue
		}
		go handleConnection(conn)
	}
}
