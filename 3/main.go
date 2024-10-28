package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	timeout = 1 * time.Second
	port    = ":8080"
)

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
			Асинхронная обработка клиентских соединений:
		•	Добавьте в сервер многопоточную обработку нескольких клиентских соединений.
		•	Используйте горутины для обработки каждого нового соединения.
		•	Реализуйте механизм graceful shutdown: сервер должен корректно завершать все активные соединения при остановке.

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
		time.Sleep(timeout)
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
