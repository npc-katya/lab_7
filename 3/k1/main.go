package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

const port = ":8080"

func main() {

	// подключение к серверу
	conn, err := net.Dial("tcp", "localhost"+port)
	if err != nil {
		fmt.Println("ошибка подключения:", err)
		return
	}
	defer conn.Close()

	for i := 0; i < 20; i++ {
		// ввод сообщения
		message := strconv.Itoa(i) + " первый клиент\n"

		// отправление сообщения на сервер
		_, err = io.WriteString(conn, message)
		if err != nil {
			fmt.Println("ошибка отправки сообщения:", err)
			return
		}

		// получение ответа сервера
		scanner := bufio.NewScanner(conn)
		if scanner.Scan() {
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("ошибка сканирования ответа:", err)
		}
		time.Sleep(1 * time.Second)
	}

}
