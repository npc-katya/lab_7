package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

const port = ":8080"

func main() {

	/*
			Реализация TCP-клиента:
		•	Разработайте TCP-клиента, который подключается к вашему серверу.
		•	Клиент должен отправлять сообщение, введённое пользователем, и ожидать ответа.
		•	После получения ответа от сервера клиент завершает соединение.

	*/

	// подключение к серверу
	conn, err := net.Dial("tcp", "localhost"+port)
	if err != nil {
		fmt.Println("ошибка подключения:", err)
		return
	}
	defer conn.Close()

	// ввод сообщения
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("введите сообщение:")
	message, _ := reader.ReadString('\n')

	// отправление сообщения на сервер
	_, err = io.WriteString(conn, message)
	if err != nil {
		fmt.Println("ошибка отправки сообщения:", err)
		return
	}

	// получение ответа сервера
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		fmt.Println("получен ответ:", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("ошибка сканирования ответа:", err)
	}
}
